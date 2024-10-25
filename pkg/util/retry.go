package util

import (
	"fmt"
	"math/rand"
	"time"
)

type RetryOption func(*retryOptions)

type retryOptions struct {
	attempts int
	sleep    time.Duration
	jitter   time.Duration
}

// WithAttempts 设置重试次数
func WithAttempts(attempts int) RetryOption {
	return func(o *retryOptions) {
		o.attempts = attempts
	}
}

// WithSleep 设置重试间隔
func WithSleep(sleep time.Duration) RetryOption {
	return func(o *retryOptions) {
		o.sleep = sleep
	}
}

// WithJitter 设置重试间隔的抖动
func WithJitter(jitter time.Duration) RetryOption {
	return func(o *retryOptions) {
		o.jitter = jitter
	}
}

// Retry 重试函数, 默认重试3次，每次间隔1秒
func Retry(fn func() error, opts ...RetryOption) error {
	options := &retryOptions{
		attempts: 3,           // 默认重试次数
		sleep:    time.Second, // 默认重试间隔
		jitter:   time.Second, // 默认重试间隔的抖动
	}

	for _, opt := range opts {
		opt(options)
	}

	var err error
	for i := 0; i < options.attempts; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		sleepTime := options.sleep + time.Duration(rand.Intn(int(options.jitter)))
		time.Sleep(sleepTime)
	}
	return fmt.Errorf("after %d attempts, last error: %s", options.attempts, err)
}
