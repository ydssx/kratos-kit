package middleware

import (
	"context"

	"github.com/ydssx/kratos-kit/pkg/limit"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func RateLimit(limiter limit.Limiter) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			r, _ := http.RequestFromServerContext(ctx)
			clientIP := getClientIP(r)
			key := clientIP + "_" + r.URL.Path
			if !limiter.Allow(key) {
				return nil, errors.Forbidden("rate limit", "rate limit")
			}

			return handler(ctx, req)
		}
	}
}
