package data

import (
	"context"

	"github.com/ydssx/kratos-kit/common"
	"github.com/ydssx/kratos-kit/internal/biz"
	"github.com/ydssx/kratos-kit/pkg/cache"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewUserRepoCacheDecorator,
	common.NewRedisCLient,
	NewRedisCache,
	common.NewMysqlDB,
	NewTransaction,
	NewUserRepo,
)

// Data .
type Data struct {
	rdb *goredis.Client
	db  *gorm.DB
}

type contextTxKey struct{}

// NewData returns a new instance of Data along with a cleanup function and error (if any).
//
// It takes the following parameters:
//   - logger: an instance of log.Logger used for logging
//   - rdb: a pointer to a goredis.Client used for Redis operations
//   - db: a pointer to a gorm.DB used for database operations
//   - collection: a pointer to a mongo.Collection used for MongoDB operations
//
// It returns the following:
//   - data: an instance of Data containing the initialized Redis, database, and MongoDB clients
//   - error: an error, if any, encountered during the initialization process
func NewData(ctx context.Context, logger log.Logger, rdb *goredis.Client, db *gorm.DB) (*Data, error) {
	cleanup := func() {
		log.Info("closing the data resources")

		if err := rdb.Close(); err != nil {
			log.Error("close redis failed:", err)
		}
	}
	context.AfterFunc(ctx, cleanup)

	return &Data{rdb: rdb, db: db}, nil
}

// InTx 在一个数据库事务中执行函数 fn。
// 使用给定的上下文 ctx 创建一个事务,在事务中执行 fn,如果成功提交事务,否则回滚。
// fn 函数在事务上下文中执行,可以通过 ctx 访问事务数据库连接。
func (d *Data) InTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

// WithTx 从上下文中获取事务数据库连接并在事务上下文中执行函数 fn。
// 如果上下文中不存在事务数据库连接,则创建一个新的事务。
func (d *Data) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if d.IsInTx(ctx) {
		return fn(ctx)
	}
	return d.InTx(ctx, fn)
}

// DB 返回当前上下文绑定的数据库连接。
// 如果上下文中存在事务(tx),则返回事务连接,否则返回默认数据库连接。
func (d *Data) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return d.db.WithContext(ctx)
}

// GetDB 返回新的数据库连接。
func (d *Data) GetDB(ctx context.Context) *gorm.DB {
	return d.db.WithContext(ctx)
}

// IsInTx 检查当前上下文是否在事务中。
func (d *Data) IsInTx(ctx context.Context) bool {
	_, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	return ok
}

func NewTransaction(d *Data) biz.Transaction {
	return d
}

func NewRedisCache(client *goredis.Client) cache.Cache {
	return cache.NewRedisCache(client)
}
