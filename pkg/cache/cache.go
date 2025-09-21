package cache

import (
	"context"
	"time"

	"github.com/ydssx/kratos-kit/pkg/logger"
	"golang.org/x/sync/singleflight"
)

type Cache interface {
	// Get 从缓存中获取指定key的值,并反序列化到result中
	Get(ctx context.Context, key string, result interface{}) error
	// Set 将指定的key/value对设置到缓存中,并设置过期时间
	Set(ctx context.Context, key string, value interface{}, expire time.Duration) error
	// Delete 从缓存中删除指定key的值
	Delete(ctx context.Context, key string) error
	// Clear 清空缓存中的所有键值对
	Clear(ctx context.Context) error
}

var g singleflight.Group

// WithCache 通用缓存装饰器
func WithCache[T any](c Cache, ctx context.Context, key string, duration time.Duration, fn func() (T, error)) (T, error) {
	var data T
	err := c.Get(ctx, key, &data)
	if err == nil {
		return data, nil
	}

	// 使用singleflight防止缓存击穿
	v, err, _ := g.Do(key, func() (interface{}, error) {
		var data T
		if err := c.Get(ctx, key, &data); err == nil {
			return data, nil
		}

		d, err := fn()
		if err != nil {
			return d, err
		}

		if err := c.Set(ctx, key, d, duration); err != nil {
			logger.Errorf(ctx, "cache set error: %v", err)
		}
		return d, nil
	})

	if err != nil {
		return data, err
	}

	return v.(T), nil
}
