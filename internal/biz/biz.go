package biz

import (
	"context"

	"github.com/ydssx/kratos-kit/common"
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

func NewContextWithUsecaseSet(ctx context.Context, serviceSet *UsecaseSet) context.Context {
	return context.WithValue(ctx, usecasetKey{}, serviceSet)
}

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

// 幂等性校验
type IdempotencyChecker interface {
	IsIdempotent(ctx context.Context, uid int, req interface{}) (bool, error)
	MarkIdempotent(ctx context.Context, uid int, req interface{}) error
}
