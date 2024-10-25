package data

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/ydssx/kratos-kit/pkg/logger"

	"github.com/ydssx/kratos-kit/internal/biz"
	"github.com/ydssx/kratos-kit/models"

	"github.com/go-kratos/kratos/v2/log"
)

var _ biz.UserRepo = (*userRepo)(nil)

type userRepo struct {
	data *Data
	log  *log.Helper
}

// GetUserByGoogleID implements biz.UserRepo.
func (r *userRepo) GetUserByGoogleID(ctx context.Context, googleID string) (*models.User, error) {
	return models.NewUserModel(r.data.DB(ctx)).SetGoogleID(googleID).FirstOne()
}

// GetUserByBrowserFingerprint implements biz.UserRepo.
func (r *userRepo) GetUserByBrowserFingerprint(ctx context.Context, fingerprint string) (*models.User, error) {
	return models.NewUserModel(r.data.DB(ctx)).SetBrowserFingerprint(fingerprint).FirstOne()
}

func (r *userRepo) GetUserVisitCount(ctx context.Context, domain string, startTime time.Time, endTime time.Time) (int64, error) {
	return models.NewUserLoginLogModel(r.data.DB(ctx)).LeftJoinUser().SetDomain(domain).GteLoginDate(startTime).LteLoginDate(endTime).Count()
}

// GetUserByIDWithLock implements biz.UserRepo.
func (r *userRepo) GetUserByIDWithLock(ctx context.Context, id uint) (*models.User, error) {
	if !r.data.IsInTx(ctx) {
		return nil, errors.New("context is not in tx")
	}
	return models.NewUserModel(r.data.DB(ctx)).SetIds(int(id)).XLock().FirstOne()
}

func NewUserRepo(data *Data, logger log.Logger) *userRepo {
	return &userRepo{data: data, log: log.NewHelper(logger)}
}

// GetUserByPhone implements biz.UserRepo.
func (r *userRepo) GetUserByPhone(ctx context.Context, phoneNumber string) (*models.User, error) {
	user, err := models.NewUserModel(r.data.DB(ctx)).FirstOne()
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// GetUserByName implements biz.UserRepo.
func (r *userRepo) GetUserByName(ctx context.Context, username string) (*models.User, error) {
	user, err := models.NewUserModel(r.data.DB(ctx)).SetUsername(username).FirstOne()
	logger.Warnf(ctx, "[GetUserByName]:err:%v::user:%v", err, user)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// GetUserByEmail implements biz.UserRepo.
func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := models.NewUserModel(r.data.DB(ctx)).SetEmail(email).FirstOne()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepo) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	user, err := models.NewUserModel(r.data.DB(ctx)).SetIds(int(id)).FirstOne()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err

	}
	return user, nil
}

func (r *userRepo) UpdateUser(ctx context.Context, userId int, user interface{}) error {
	return models.NewUserModel(r.data.DB(ctx)).SetIds(int(userId)).Updates(user)
}

// GetUserByUUID 根据用户UUID获取用户
func (r *userRepo) GetUserByUUID(ctx context.Context, uuid string) (*models.User, error) {
	user, err := models.NewUserModel(r.data.DB(ctx)).SetUUIds(uuid).FirstOne()
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *userRepo) UpdateUserByID(ctx context.Context, uid int, user *models.User) error {
	return models.NewUserModel(r.data.DB(ctx)).SetIds(uid).Updates(user)
}

func (r *userRepo) DeleteUser(ctx context.Context, id uint) error {
	return models.NewUserModel(r.data.DB(ctx)).SetIds(int(id)).Delete()
}

// CreateUser implements biz.UserRepo.
func (r *userRepo) CreateUser(ctx context.Context, user *models.User) (userId int, err error) {
	userInfo, err := models.NewUserModel(r.data.DB(ctx)).Create(*user)
	if err != nil {
		return 0, err
	}
	return int(userInfo.ID), nil
}

// ListUser 根据条件查询用户列表
// ctx 上下文
// cond 查询条件
// 返回用户列表
func (r *userRepo) ListUser(ctx context.Context, cond *biz.ListUserCond) []models.User {
	if cond == nil {
		cond = new(biz.ListUserCond)
	}

	userModel := models.NewUserModel(r.data.DB(ctx)).PointsGt(int(cond.PointsGt))
	if cond.Type != nil {
		userModel.SetUserType(*cond.Type)
	}

	users, _ := userModel.List()

	return users
}

// GetTotalUsers 获取总用户数
func (r *userRepo) GetTotalUsers(ctx context.Context) (int, error) {
	total, err := models.NewUserModel(r.data.DB(ctx)).Count()
	if err != nil {
		return 0, err
	}
	return int(total), nil
}

// GetUsersByIDs 根据用户ID获取用户
func (r *userRepo) GetUsersByIDs(ctx context.Context, ids ...int) ([]models.User, error) {
	users, err := models.NewUserModel(r.data.DB(ctx)).SetIds(ids...).List()
	if err != nil {
		return nil, err
	}
	return users, nil
}
