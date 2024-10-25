package middleware

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
)

func Timeout(timeout time.Duration) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}

			// 创建一个带超时的上下文
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			// 创建一个 channel 用于捕获 handler 的响应
			done := make(chan struct{})
			var resp interface{}
			var err error

			// 启动一个 goroutine 来执行 handler
			go func() {
				defer close(done)
				newCtx := context.WithoutCancel(ctx)
				resp, err = handler(newCtx, req)
			}()

			// 等待响应或超时
			select {
			case <-done:
				// handler 完成
				return resp, err
			case <-ctx.Done():
				// 超时
				return nil, ctx.Err()
			}
		}
	}
}
