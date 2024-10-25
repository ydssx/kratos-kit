package service

import (
	"context"

	userv1 "github.com/ydssx/kratos-kit/api/user/v1"
	"github.com/ydssx/kratos-kit/internal/biz"
	"github.com/ydssx/kratos-kit/pkg/util"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GreeterService is a greeter service.
type UserService struct {
	userv1.UserServiceHTTPServer
	uc *biz.UserUseCase
}

func NewUserService(uc *biz.UserUseCase) *UserService {
	return &UserService{uc: uc}
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	return s.uc.Login(ctx, req)
}

// GetUser 获取用户信息
func (s *UserService) GetUser(ctx context.Context, g *emptypb.Empty) (*userv1.GetUserResponse, error) {
	return s.uc.GetUser(ctx, g)
}

// Create 创建用户
func (s *UserService) Create(ctx context.Context, req *userv1.CreateRequest) (*userv1.LoginResponse, error) {
	return s.uc.Create(ctx, req)
}

// IsAccountExist 检测账号是否已存在
func (s *UserService) IsAccountExist(ctx context.Context, req *userv1.IsAccountExistRequest) (*userv1.IsAccountExistResponse, error) {
	return s.uc.IsAccountExist(ctx, req)
}

// GoogleLogin 获取Google登录URL
func (s *UserService) GoogleLogin(ctx context.Context, req *emptypb.Empty) (*userv1.GoogleLoginResponse, error) {
	return s.uc.GoogleLogin(ctx, req)
}

// GoogleCallback 处理Google回调
func (s *UserService) GoogleCallback(ctx *gin.Context) {
	res, err := s.uc.GoogleCallback(ctx)
	if err != nil {
		util.FailWithError(ctx, err)
		return
	}
	util.OKWithData(ctx, res)
}

// Logout 用户登出
func (s *UserService) Logout(ctx context.Context, req *emptypb.Empty) (res *userv1.LoginResponse, err error) {
	return s.uc.Logout(ctx, req)
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, req *userv1.RegisterRequest) (res *userv1.LoginResponse, err error) {
	return s.uc.Register(ctx, req)
}

// SendVerificationCode 发送验证码
func (s *UserService) SendVerificationCode(ctx context.Context, req *userv1.SendVerificationCodeRequest) (res *emptypb.Empty, err error) {
	return s.uc.SendVerificationCode(ctx, req)
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (res *emptypb.Empty, err error) {
	return s.uc.UpdateUser(ctx, req)
}
