package limit

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

type LimiterType int

const (
	LimiterTypeConnection    LimiterType = iota // 连接数和ip地址限流
	LimiterTypeSlidingWindow                    // 滑动窗口限流
	LimiterTypeTokenBucket                      // 令牌桶限流
)

type option struct {
	ratePerSecond int
	burst         int
}

type Option func(*option)

// WithRatePerSecond returns an Option that sets the maximum number of requests
// allowed per second. This can be used to rate limit requests.
func WithRatePerSecond(ratePerSecond int) Option {
	return func(o *option) { o.ratePerSecond = ratePerSecond }
}

// WithBurst returns an Option that sets the maximum number of requests allowed
// to burst before rate limiting applies. This can be used in conjunction with
// WithRatePerSecond to allow short bursts above the sustained rate limit.
func WithBurst(burst int) Option {
	return func(o *option) { o.burst = burst }
}

type Limiter interface {
	Allow(key string, opts ...Option) bool
}

type RedisLimiter struct {
	*redis_rate.Limiter
}

func NewRedisLimiter(rdb *redis.Client) *RedisLimiter {
	return &RedisLimiter{redis_rate.NewLimiter(rdb)}
}

// Limit 检查给定的context中的限流key,如果允许则返回nil,否则返回错误。
// 它会从context中提取key,然后使用Allow方法检查key的限流状态。
// 如果允许,则返回nil,否则返回一个rate limited的错误。
func (l *RedisLimiter) Limit(ctx context.Context) error {
	key := LimitKeyFromCtx(ctx).(string)
	if l.Allow(key) {
		return nil
	}
	return errors.New("rate limited.")
}

// Allow checks if the given key is allowed by the rate limiter.
// It applies the given options to configure the rate and burst limits.
// It uses the redis client to check against the limiter and returns
// whether the request is allowed.
func (l *RedisLimiter) Allow(key string, opts ...Option) bool {
	opt := &option{ratePerSecond: 10, burst: 20}
	for _, v := range opts {
		v(opt)
	}

	r, err := l.Limiter.Allow(context.Background(), key, perSecond(opt.ratePerSecond, opt.burst))
	if err != nil {
		panic(err)
	}

	return r.Allowed > 0
}

func (l *RedisLimiter) Reset(key string) {
	_ = l.Limiter.Reset(context.Background(), key)
}

func perSecond(rate, burst int) redis_rate.Limit {
	burst = int(math.Max(float64(rate), float64(burst)))
	return redis_rate.Limit{
		Rate:   rate,
		Period: time.Second,
		Burst:  burst,
	}
}

type limitKey struct{}

func LimitKeyFromCtx(ctx context.Context) any {
	return ctx.Value(limitKey{})
}

func CtxWithLimitKey(ctx context.Context, value any) context.Context {
	return context.WithValue(ctx, limitKey{}, value)
}
