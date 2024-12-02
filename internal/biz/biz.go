package biz

import (
	"context"

	"github.com/ydssx/kratos-kit/common"
	"github.com/ydssx/kratos-kit/models"
	"github.com/ydssx/kratos-kit/pkg/lock"
	"github.com/ydssx/kratos-kit/pkg/storage"

	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	common.NewGoogleCloudStorage,
	common.InitGoogleOAuth,
	common.NewEmail,
	common.NewWsService,
	wire.Bind(new(storage.Storage), new(*storage.GoogleCloudStorage)),
	common.NewRedisLocker,
	wire.Bind(new(lock.Locker), new(*lock.RedisLocker)),
	common.NewGeoipDB,
	NewUsecaseSet,
	NewUserUseCase,
	NewUploadUseCase,
	NewCommonUseCase,
	NewAdminUseCase,
	NewAiUseCase,
)

type UsecaseSet struct {
	UserBiz   *UserUseCase
	UploadBiz *UploadUseCase
}

func NewUsecaseSet(
	userBiz *UserUseCase,
	uploadBiz *UploadUseCase,
) *UsecaseSet {
	return &UsecaseSet{
		UserBiz:   userBiz,
		UploadBiz: uploadBiz,
	}
}

type usecasetKey struct{}

// NewContextWithUsecaseSet 创建一个包含UsecaseSet的上下文
func NewContextWithUsecaseSet(ctx context.Context, serviceSet *UsecaseSet) context.Context {
	return context.WithValue(ctx, usecasetKey{}, serviceSet)
}

// UsecaseSetFromContext 从上下文中获取UsecaseSet
func UsecaseSetFromContext(ctx context.Context) *UsecaseSet {
	return ctx.Value(usecasetKey{}).(*UsecaseSet)
}

// Transaction is a interface that can be used to wrap a business logic function with a database transaction.
type Transaction interface {
	// InTx runs the given business logic function within a database transaction. If the function returns an error, the transaction will be rolled back. Otherwise, the transaction will be committed.
	InTx(context.Context, func(ctx context.Context) error) error

	// WithTx 从上下文中获取事务数据库连接并在事务上下文中执行函数 fn。
	// 如果上下文中不存在事务数据库连接,则创建一个新的事务。
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

// IdempotencyChecker 幂等性校验
type IdempotencyChecker interface {
	IsIdempotent(ctx context.Context, uid int, req interface{}) (bool, error)
	MarkIdempotent(ctx context.Context, uid int, req interface{}) error
}

// UserRepo 用户仓库接口
type UserRepo interface {
	// CreateUser 创建用户
	CreateUser(ctx context.Context, user *models.User) (userId int, err error)
	// UpdateUser 更新用户
	UpdateUser(ctx context.Context, userId int, user interface{}) error
	// UpdateUserByID 根据ID更新用户
	UpdateUserByID(ctx context.Context, uid int, user *models.User) error
	// ListUser 获取用户列表
	ListUser(ctx context.Context, cond *ListUserCond) []models.User
	// GetUserByID 根据用户ID获取用户
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
	// GetUserByUUID 根据用户UUID获取用户
	GetUserByUUID(ctx context.Context, uuid string) (*models.User, error)
	// GetUserByEmail 根据邮箱名获取用户
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	// GetUserByBrowserFingerprint 根据指纹获取用户
	GetUserByBrowserFingerprint(ctx context.Context, fingerprint string) (*models.User, error)
	// GetUserByGoogleID 根据Google ID获取用户
	GetUserByGoogleID(ctx context.Context, googleID string) (*models.User, error)
}

type (
	ListUserCond struct {
		Type *models.UserType
		// 积分数大于
		PointsGt int64
	}
)
