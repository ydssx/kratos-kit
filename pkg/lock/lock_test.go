package lock

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisSync_LockerInterface(t *testing.T) {
	// 创建 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()

	// 创建 RedisSync 实例
	rs := NewRedisSync(rdb)

	ctx := context.Background()
	key := "test_lock"

	// 测试基本的锁操作
	t.Run("Basic Lock Operations", func(t *testing.T) {
		// 获取锁
		err := rs.Lock(ctx, key)
		assert.NoError(t, err)

		// 尝试再次获取同一个锁应该失败
		err = rs.TryLock(ctx, key)
		assert.Error(t, err)

		// 释放锁
		err = rs.Unlock(ctx, key)
		assert.NoError(t, err)
	})

	// 测试带选项的锁操作
	t.Run("Lock With Options", func(t *testing.T) {
		err := rs.Lock(ctx, key, 
			WithTTL(time.Second),
			WithTries(3),
			WithDelay(100*time.Millisecond),
		)
		assert.NoError(t, err)

		err = rs.Unlock(ctx, key)
		assert.NoError(t, err)
	})

	// 测试清理功能
	t.Run("Cleanup", func(t *testing.T) {
		err := rs.Lock(ctx, "lock1")
		assert.NoError(t, err)
		
		err = rs.Lock(ctx, "lock2")
		assert.NoError(t, err)

		rs.Cleanup()

		// 验证锁已被清理
		err = rs.Lock(ctx, "lock1")
		assert.NoError(t, err)
		err = rs.Lock(ctx, "lock2")
		assert.NoError(t, err)
	})
}
