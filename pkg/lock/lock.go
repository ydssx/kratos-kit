package lock

import (
	"context"
	"sync"
	"time"

	"github.com/ydssx/kratos-kit/pkg/errors"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

type Locker interface {
	// Lock 获取锁
	Lock(ctx context.Context, key string, opt ...LockerOption) error
	// Unlock 释放锁
	Unlock(ctx context.Context, key string) error
	// TryLock 尝试获取锁
	TryLock(ctx context.Context, key string, opt ...LockerOption) error
}

type lockOption struct {
	// 锁的有效期
	ttl time.Duration
	// 重试次数
	tries int
	// 重试间隔
	delay time.Duration
	// 超时时间
	timeout time.Duration
}

type LockerOption func(*lockOption)

// WithTTL 设置锁的有效期
func WithTTL(ttl time.Duration) LockerOption {
	return func(opt *lockOption) {
		opt.ttl = ttl
	}
}

// WithTries 设置重试次数
func WithTries(tries int) LockerOption {
	return func(opt *lockOption) {
		opt.tries = tries
	}
}

// WithDelay 设置重试间隔
func WithDelay(delay time.Duration) LockerOption {
	return func(opt *lockOption) {
		opt.delay = delay
	}
}

func WithTimeout(timeout time.Duration) LockerOption {
	return func(opt *lockOption) {
		opt.timeout = timeout
	}
}

type RedisLocker struct {
	*redislock.Client
	lockMap sync.Map
}

var _ Locker = (*RedisLocker)(nil)

func NewLocker(cli *redis.Client) *RedisLocker {
	return &RedisLocker{Client: redislock.New(cli)}
}

// Lock 获取锁
//
// 实现了分布式锁的功能，通过Redis实现。
//
// ctx: 上下文
// key: 锁的key
// opt: 可选参数
//
// 可选参数:
//
//	WithTTL(ttl time.Duration): 设置锁的有效期，默认10s。
//	WithTries(tries int): 设置重试次数，默认3次。
//	WithDelay(delay time.Duration): 设置重试间隔，默认500ms。
func (r *RedisLocker) Lock(ctx context.Context, key string, opt ...LockerOption) error {
	lo := &lockOption{
		ttl:   10 * time.Second,
		tries: 3,
		delay: 100 * time.Millisecond,
	}
	for _, o := range opt {
		o(lo)
	}
	lock, err := r.Client.Obtain(ctx, key, lo.ttl, &redislock.Options{
		RetryStrategy: redislock.LinearBackoff(lo.delay),
	})
	if err != nil {
		return errors.Errorf("redislock: lock %s failed: %v", key, err)
	}

	// Add the lock to the map
	r.lockMap.Store(key, lock)
	return nil
}

// Unlock releases the lock obtained by Lock or TryLock.
// Returns an error if the lock is not held or the underlying Redis
// command fails.
//
// It returns an error with the message "redislock: not obtain lock"
// if the lock is not held.
func (r *RedisLocker) Unlock(ctx context.Context, key string) error {
	lockInterface, exists := r.lockMap.LoadAndDelete(key)
	if !exists {
		return redislock.ErrLockNotHeld
	}

	lock := lockInterface.(*redislock.Lock)
	// Release the lock
	if err := lock.Release(ctx); err != nil && !errors.Is(err, redislock.ErrLockNotHeld) {
		return errors.Errorf("redislock: unlock %s failed: %v", lock.Key(), err)
	}

	return nil
}

// TryLock tries to obtain a lock for the given key. It returns
// nil if the lock is obtained, otherwise it returns an error.
//
// The context will be used to establish a deadline for the lock
// acquisition.
//
// The key will be locked for the duration specified in the options
// parameter (WithTTL). If no duration is specified, the lock will
// be obtained with a default TTL of 10s.
//
// If the lock is not obtained within the specified deadline, an
// error is returned.
func (r *RedisLocker) TryLock(ctx context.Context, key string, opt ...LockerOption) error {
	// Set a deadline for the lock acquisition
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
	defer cancel()

	// Try to obtain the lock
	if err := r.Lock(ctx, key, opt...); err != nil {
		return errors.Wrapf(err, "redislock: try lock %s failed", key)
	}

	return nil
}
