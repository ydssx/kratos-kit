package biz

import (
	"context"
	"net/url"
	"strings"
	"sync"
	"time"

	userv1 "github.com/ydssx/kratos-kit/api/user/v1"
	"github.com/ydssx/kratos-kit/constants"
	"github.com/ydssx/kratos-kit/internal/middleware"
	"github.com/ydssx/kratos-kit/models"
	"github.com/ydssx/kratos-kit/pkg/cache"
	"github.com/ydssx/kratos-kit/pkg/email"
	"github.com/ydssx/kratos-kit/pkg/errors"
	"github.com/ydssx/kratos-kit/pkg/jwt"
	"github.com/ydssx/kratos-kit/pkg/lock"
	"github.com/ydssx/kratos-kit/pkg/logger"
	"github.com/ydssx/kratos-kit/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/oauth2"
	goauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

// UserRepo 是用户仓库接口
type UserRepo interface {
	// CreateUser 创建用户
	CreateUser(ctx context.Context, user *models.User) (userId int, err error)
	// UpdateUser 更新用户
	UpdateUser(ctx context.Context, userId int, user interface{}) error
	// UpdateUserByID 根据ID更新用户
	UpdateUserByID(ctx context.Context, uid int, user *models.User) error
	// DeleteUser 删除用户
	DeleteUser(ctx context.Context, id uint) error
	// ListUser 获取用户列表
	ListUser(ctx context.Context, cond *ListUserCond) []models.User
	// GetUserByID 根据用户ID获取用户
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
	// GetUserByUUID 根据用户UUID获取用户
	GetUserByUUID(ctx context.Context, uuid string) (*models.User, error)
	// GetUserByName 根据用户名获取用户
	GetUserByName(ctx context.Context, username string) (*models.User, error)
	// GetUserByEmail 根据邮箱名获取用户
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	// GetUserByPhone 根据手机号码获取用户
	GetUserByPhone(ctx context.Context, phoneNumber string) (*models.User, error)
	// GetUserByIDWithLock 根据用户ID获取用户（加锁）
	GetUserByIDWithLock(ctx context.Context, id uint) (*models.User, error)
	// GetUserByBrowserFingerprint 根据指纹获取用户
	GetUserByBrowserFingerprint(ctx context.Context, fingerprint string) (*models.User, error)
	// GetUserByGoogleID 根据Google ID获取用户
	GetUserByGoogleID(ctx context.Context, googleID string) (*models.User, error)
	// GetTotalUsers 获取总用户数
	GetTotalUsers(ctx context.Context) (int, error)
	// GetUsersByIDs 根据用户ID获取用户
	GetUsersByIDs(ctx context.Context, ids ...int) ([]models.User, error)
}

type (
	ListUserCond struct {
		Type *models.UserType
		// 积分数大于
		PointsGt int64
	}
	ListRoleCond struct {
		Page  int64  // 页码
		Limit int64  // 条数
		Name  string // 角色名称
	}
	AdminUserCollectCond struct {
		VisitorsNum int32
		PayNum      int32
		NonPayNum   int32
		FirstNum    int32
		RenewNum    int32
	}
)

type UserUseCase struct {
	repo              UserRepo
	log               *log.Helper
	tm                Transaction // 事务管理器
	commonUc          *CommonUseCase
	mu                sync.Mutex
	locker            lock.Locker
	googleOauthConfig *oauth2.Config
	cache             cache.Cache
	email             *email.Email
}

func NewUserUseCase(
	userRepo UserRepo,
	logger log.Logger,
	transaction Transaction,
	commonUc *CommonUseCase,
	locker lock.Locker,
	googleOauthConfig *oauth2.Config,
	cache cache.Cache,
	email *email.Email,
) *UserUseCase {
	return &UserUseCase{
		repo:              userRepo,
		log:               log.NewHelper(logger),
		tm:                transaction,
		commonUc:          commonUc,
		locker:            locker,
		googleOauthConfig: googleOauthConfig,
		cache:             cache,
		email:             email,
	}
}

// Login 用户登陆api
func (uc *UserUseCase) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	// 验证请求参数
	if req.Email == "" {
		return nil, errors.NewUserError("email is required")
	}

	// 根据邮箱获取用户
	user, err := uc.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.NewUserError("user not found")
	}

	// 验证密码或验证码
	if req.Password != "" {
		// 密码登录
		if util.MD5(req.Password) != user.PasswordHash {
			return nil, errors.NewUserError("password is incorrect")
		}
	} else if req.Code != "" {
		// 验证码登录
		code, err := uc.GetVerificationCode(ctx, req.Email)
		if err != nil || code != req.Code {
			return nil, errors.NewUserError("verification code is incorrect or expired")
		}
	} else {
		return nil, errors.NewUserError("please provide a password or verification code")
	}

	// 返回登录响应
	return &userv1.LoginResponse{
		Uuid: user.UUID,
	}, nil
}

// GetUser 根据userID获取用户信息
func (uc *UserUseCase) GetUser(ctx context.Context, g *emptypb.Empty) (*userv1.GetUserResponse, error) {
	// 请求头中获取uuid
	claim := middleware.GetClaims(ctx)
	userId := uint(claim.Uid)
	// 获取用户信息
	user, err := uc.repo.GetUserByID(ctx, userId)
	if err != nil {
		return nil, errors.New("failed to get user")
	}
	// 查询红点
	// 返回
	return &userv1.GetUserResponse{
		Id:         int32(user.ID),
		Username:   user.Username,
		Email:      user.Email,
		AvatarPath: user.AvatarPath,
	}, nil
}

// Create 创建用户
func (uc *UserUseCase) Create(ctx context.Context, req *userv1.CreateRequest) (res *userv1.LoginResponse, err error) {
	// TODO 根据x-u-key和token校验客户端动态加密结果
	header := middleware.GetHeaderInfo(ctx)
	res = &userv1.LoginResponse{}

	if req.XUKey != "" {
		err = uc.locker.Lock(ctx, req.XUKey, lock.WithTTL(time.Second*2))
		if err != nil {
			return nil, errors.New("failed to lock")
		}
		defer func() {
			err = uc.locker.Unlock(ctx, req.XUKey)
			if err != nil {
				logger.Errorf(ctx, "failed to unlock: %s", err.Error())
			}
		}()
	}

	user, err := uc.repo.GetUserByBrowserFingerprint(ctx, req.XUKey)
	if err == nil && user.IPAddress == header.ClientIP && user.Email == "" {
		res.Uuid = user.UUID
		return
	}

	// 创建用户
	uuid := util.GetUUID()
	res.Uuid = uuid
	_, err = uc.repo.CreateUser(ctx, &models.User{
		UUID:               uuid,
		Platform:           header.Platform,
		BrowserFingerprint: req.XUKey,
		IPAddress:          header.ClientIP,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create user")
	}

	return res, nil
}

// IsAccountExist 检测邮箱账号是否存在
func (uc *UserUseCase) IsAccountExist(ctx context.Context,
	req *userv1.IsAccountExistRequest,
) (*userv1.IsAccountExistResponse, error) {
	_, err := uc.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err.Error() == "user not found" {
			return &userv1.IsAccountExistResponse{
				IsExist: false,
			}, nil
		}
		return nil, err
	}
	return &userv1.IsAccountExistResponse{
		IsExist: true,
	}, nil
}

// 处理Google回调
func (uc *UserUseCase) GoogleCallback(ctx *gin.Context) (*userv1.LoginResponse, error) {
	code := ctx.Query("code")
	// state := ctx.Query("state")
	jwtToken := ctx.Query("token")
	code, _ = url.QueryUnescape(code)

	var userInfo *goauth2.Userinfo
	var err error
	var token *oauth2.Token
	if jwtToken != "" {
		googleClaims, err := uc.parseJWTToken(jwtToken)
		if err != nil {
			return nil, errors.Wrap(err, "failed to verify jwt token")
		}
		userInfo = googleClaims.Userinfo
		userInfo.Id = googleClaims.Sub
	} else {
		token, err = uc.googleOauthConfig.Exchange(ctx, code)
		if err != nil {
			return nil, errors.Wrap(err, "failed to exchange token")
		}

		service, err := goauth2.NewService(ctx, option.WithTokenSource(uc.googleOauthConfig.TokenSource(ctx, token)))
		if err != nil {
			return nil, errors.Wrap(err, "failed to create service")
		}

		userInfo, err = service.Userinfo.Get().Do()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get user info")
		}
	}

	user, err := uc.repo.GetUserByGoogleID(ctx, userInfo.Id)
	if err != nil {
		user, err = uc.repo.GetUserByEmail(ctx, userInfo.Email)
		if err != nil {
			uuid := util.GetUUID()
			_, err := uc.repo.CreateUser(ctx, &models.User{
				UUID:       uuid,
				Email:      userInfo.Email,
				Username:   userInfo.Name,
				FirstName:  userInfo.GivenName,
				LastName:   userInfo.FamilyName,
				AvatarPath: userInfo.Picture,
				GoogleId:   userInfo.Id,
			})
			if err != nil {
				return nil, errors.Wrap(err, "failed to create user")
			}
			return &userv1.LoginResponse{
				Uuid: uuid,
			}, nil
		}
	}

	if token != nil {
	}
	// 更新用户的Google ID和头像
	user.GoogleId = userInfo.Id
	user.AvatarPath = userInfo.Picture
	user.Email = userInfo.Email
	user.Username = userInfo.Name
	user.FirstName = userInfo.GivenName
	user.LastName = userInfo.FamilyName
	err = uc.repo.UpdateUserByID(ctx, int(user.ID), user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update user")
	}

	return &userv1.LoginResponse{
		Uuid: user.UUID,
	}, nil
}

// GoogleLogin Google登录
func (uc *UserUseCase) GoogleLogin(ctx context.Context, req *emptypb.Empty) (res *userv1.GoogleLoginResponse, err error) {
	res = new(userv1.GoogleLoginResponse)

	claim := middleware.GetClaims(ctx)

	res.Url = uc.googleOauthConfig.AuthCodeURL(claim.Uuid, oauth2.AccessTypeOffline)

	return
}

// JWT令牌解析
func (uc *UserUseCase) parseJWTToken(token string) (user *GoogleClaims, err error) {
	claims, err := jwt.DecodeJWT(token)
	if err != nil {
		return nil, err
	}
	err = util.MapDecode(claims.Payload, &user)
	if err != nil {
		return nil, err
	}
	return
}

// Logout 用户登出
func (uc *UserUseCase) Logout(ctx context.Context, req *emptypb.Empty) (res *userv1.LoginResponse, err error) {
	res = new(userv1.LoginResponse)
	header := middleware.GetHeaderInfo(ctx)

	uuid := util.GetUUID()
	user := &models.User{
		UUID:      uuid,
		Type:      int(models.UserTypeLogout),
		IPAddress: header.ClientIP,
		Platform:  header.Platform,
	}

	_, err = uc.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create user")
	}

	res.Uuid = uuid

	return
}

// Register 用户注册
func (uc *UserUseCase) Register(ctx context.Context, req *userv1.RegisterRequest) (res *userv1.LoginResponse, err error) {
	res = new(userv1.LoginResponse)

	// 验证验证码
	storedCode, err := uc.GetVerificationCode(ctx, req.Email)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get verification code")
	}
	if storedCode != req.Code {
		return nil, errors.NewUserError("verification code is incorrect")
	}

	// 检查邮箱是否已被注册
	existingUser, err := uc.repo.GetUserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(err, "failed to check if email exists")
	}
	if existingUser != nil {
		return nil, errors.NewUserError("email already registered")
	}

	newUser := &models.User{
		UUID:         util.GetUUID(),
		Email:        req.Email,
		PasswordHash: util.MD5(req.Password),
		IPAddress:    middleware.GetHeaderInfo(ctx).ClientIP,
		Username:     strings.Split(req.Email, "@")[0],
	}

	_, err = uc.repo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create user")
	}

	res.Uuid = newUser.UUID

	// 删除验证码
	// err = uc.cache.Delete(ctx, "verification:"+req.Email)
	// if err != nil {
	// 	uc.log.Errorf("删除验证码失败: %v", err)
	// }

	return res, nil
}

func (uc *UserUseCase) SetVerificationCode(ctx context.Context, email, code string) error {
	return uc.cache.Set(ctx, "verification:"+email, code, 5*time.Minute)
}

func (uc *UserUseCase) GetVerificationCode(ctx context.Context, email string) (code string, err error) {
	err = uc.cache.Get(ctx, "verification:"+email, &code)
	return
}

// SendVerificationCode 发送验证码
func (uc *UserUseCase) SendVerificationCode(ctx context.Context, req *userv1.SendVerificationCodeRequest) (res *emptypb.Empty, err error) {
	_, err = uc.GetVerificationCode(ctx, req.Email)
	if err == nil {
		return nil, errors.ErrVerificationCodeSent
	}

	// 生成随机验证码
	code := util.GenerateCode(6)

	// 发送验证码邮件
	err = uc.email.SendVerificationCode(req.Email, code, constants.EmailLoginTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send verification code")
	}

	// 将验证码存储到 Redis
	err = uc.SetVerificationCode(ctx, req.Email, code)
	if err != nil {
		return nil, errors.Wrap(err, "failed to store verification code")
	}

	return &emptypb.Empty{}, nil
}

// UpdateUser 更新用户信息
func (uc *UserUseCase) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (res *emptypb.Empty, err error) {
	res = new(emptypb.Empty)
	userID := middleware.GetClaims(ctx).Uid

	// 检查用户是否存在
	user, err := uc.repo.GetUserByID(ctx, uint(userID))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user")
	}
	if user == nil {
		return nil, errors.NewUserError("user not found")
	}

	// 更新用户信息
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.AvatarPath != "" {
		user.AvatarPath = req.AvatarPath
	}

	// 保存更新后的用户信息
	err = uc.repo.UpdateUserByID(ctx, int(user.ID), user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update user")
	}

	return res, nil
}
