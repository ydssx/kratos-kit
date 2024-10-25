package redis

import (
	"context"

	"github.com/ydssx/kratos-kit/pkg/errors"
	"github.com/ydssx/kratos-kit/pkg/logger"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
)

// NewRedis 连接Redis并返回Client对象
func NewRedis(opt *redis.Options) (*redis.Client, error) {
	cli := redis.NewClient(opt)
	_, err := cli.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.Wrap(err, "redis connect failed")
	}
	logger.Info(context.Background(), "init redis success")
	return cli, nil
}

// NewRedisCluster 连接Redis集群并返回ClusterClient对象
func NewRedisCluster(opt *redis.ClusterOptions) *redis.ClusterClient {
	cli := redis.NewClusterClient(opt)
	_, err := cli.Ping(context.Background()).Result()
	if err != nil {
		log.Error(err)
	}
	return cli
}
