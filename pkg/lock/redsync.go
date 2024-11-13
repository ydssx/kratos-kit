package lock

import (
	"context"
	"sync"
	"time"

	"github.com/ydssx/kratos-kit/pkg/errors"

	"github.com/go-redsync/redsync/v4"
	syncredis "github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

// Error definitions
var (
	ErrLockFailed    = errors.New("redisSync: lock failed")
	ErrUnlockFailed  = errors.New("redisSync: unlock failed")
	ErrLockNotExists = errors.New("redisSync: lock not exists")
)

// 在文件顶部添加以下代码以确保接口实现
var _ Locker = (*RedisSync)(nil)

type RedisSync struct {
	redsync  *redsync.Redsync
	mutexMap sync.Map // 使用sync.Map替代map+mutex
	
	// 默认配置
	defaultTTL   time.Duration
	defaultTries int
	defaultDelay time.Duration
}

// RedisOption 定义Redis配置选项
type RedisOption func(*RedisSync)

// WithDefaultTTL 设置默认过期时间
func WithDefaultTTL(ttl time.Duration) RedisOption {
	return func(rs *RedisSync) {
		rs.defaultTTL = ttl
	}
}

// WithDefaultTries 设置默认重试次数
func WithDefaultTries(tries int) RedisOption {
	return func(rs *RedisSync) {
		rs.defaultTries = tries
	}
}

// WithDefaultDelay 设置默认重试延迟
func WithDefaultDelay(delay time.Duration) RedisOption {
	return func(rs *RedisSync) {
		rs.defaultDelay = delay
	}
}

func NewRedisSync(cli *redis.Client, opts ...RedisOption) *RedisSync {
	pool := syncredis.NewPool(cli)
	rs := &RedisSync{
		redsync:      redsync.New(pool),
		defaultTTL:   8 * time.Second,  // 默认值
		defaultTries: 32,
		defaultDelay: 500 * time.Millisecond,
	}
	
	for _, opt := range opts {
		opt(rs)
	}
	
	return rs
}

func (r *RedisSync) Lock(ctx context.Context, key string, opt ...LockerOption) error {
	m := r.newMutex(opt, key)
	if err := m.LockContext(ctx); err != nil {
		return errors.Wrap(err, ErrLockFailed.Error())
	}
	r.mutexMap.Store(key, m)
	return nil
}

func (r *RedisSync) Unlock(ctx context.Context, key string) error {
	value, exists := r.mutexMap.Load(key)
	if !exists {
		return ErrLockNotExists
	}

	m, ok := value.(*redsync.Mutex)
	if !ok {
		return ErrLockNotExists
	}

	ok, err := m.UnlockContext(ctx)
	if err != nil {
		if _, is := err.(*redsync.ErrTaken); !is {
			return errors.Errorf("redisSync: unlock %s failed: %v", m.Name(), err)
		}
	}
	if !ok {
		return ErrUnlockFailed
	}
	
	r.mutexMap.Delete(key)
	return nil
}

func (r *RedisSync) TryLock(ctx context.Context, key string, opt ...LockerOption) error {
	m := r.newMutex(opt, key)
	if err := m.TryLockContext(ctx); err != nil {
		return errors.Wrap(err, ErrLockFailed.Error())
	}
	r.mutexMap.Store(key, m)
	return nil
}

func (r *RedisSync) newMutex(opt []LockerOption, key string) *redsync.Mutex {
	var o lockOption
	for _, f := range opt {
		f(&o)
	}
	
	opts := []redsync.Option{
		redsync.WithExpiry(r.defaultTTL),
		redsync.WithTries(r.defaultTries),
		redsync.WithRetryDelay(r.defaultDelay),
	}

	if o.ttl != 0 {
		opts = append(opts, redsync.WithExpiry(o.ttl))
	}
	if o.tries != 0 {
		opts = append(opts, redsync.WithTries(o.tries))
	}
	if o.delay != 0 {
		opts = append(opts, redsync.WithRetryDelay(o.delay))
	}
	
	return r.redsync.NewMutex(key, opts...)
}

// Cleanup 清理过期的锁
func (r *RedisSync) Cleanup() {
	r.mutexMap.Range(func(key, value interface{}) bool {
		if m, ok := value.(*redsync.Mutex); ok {
			// 尝试解锁,忽略错误
			_, _ = m.Unlock()
		}
		r.mutexMap.Delete(key)
		return true
	})
}
