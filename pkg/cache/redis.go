package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/ydssx/kratos-kit/pkg/errors"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

var (
	_           Cache = (*RedisCache)(nil)
	cachePrefix       = "cache:"
)

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

// Get 从redis中获取指定key的值,并反序列化到result中
// 如果key不存在,将返回key not found错误
// 如果发生其他错误,将直接返回错误
func (c *RedisCache) Get(ctx context.Context, key string, result interface{}) error {
	val, err := c.client.Get(ctx, cachePrefix+key).Result()
	if err == redis.Nil {
		return fmt.Errorf("key %s not found", key)
	}
	if err != nil {
		return errors.Wrap(err, "get redis key error")
	}
	err = json.Unmarshal([]byte(val), &result)
	return errors.Wrap(err, "unmarshal redis value error")
}

// Set 将指定的key/value对设置到redis中,并设置过期时间
// 如果value不是json可序列化的,将返回错误
// 过期时间会在指定的过期时间上再加上0-1秒的随机值,是为了防止大量key同时过期
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return errors.Wrap(err, "marshal value error")
	}

	randomExpire := expire + time.Duration(rand.Int63n(int64(time.Second)))

	err = c.client.Set(ctx, cachePrefix+key, string(data), randomExpire).Err()
	if err != nil {
		return errors.Wrap(err, "set redis key error")
	}
	return nil
}

// Delete deletes the key from redis.
// It returns an error if there was a problem deleting the key.
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, cachePrefix+key).Err()
	return err
}

// Clear 清空缓存中的所有键值对
func (c *RedisCache) Clear(ctx context.Context) error {
	// 用scan扫描key
	iter := c.client.Scan(ctx, 0, cachePrefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		err := c.client.Del(ctx, iter.Val()).Err()
		if err != nil {
			return errors.Wrap(err, "clear redis key error")
		}
	}
	return nil
}
