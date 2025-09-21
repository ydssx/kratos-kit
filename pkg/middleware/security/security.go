package security

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"golang.org/x/time/rate"
)

// SecurityHeaders adds security-related headers to the response
func SecurityHeaders() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				header := tr.ReplyHeader()
				// Security headers
				header.Set("X-Content-Type-Options", "nosniff")
				header.Set("X-Frame-Options", "DENY")
				header.Set("X-XSS-Protection", "1; mode=block")
				header.Set("Referrer-Policy", "strict-origin-when-cross-origin")
				header.Set("Content-Security-Policy", "default-src 'self'")
				header.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			}
			return handler(ctx, req)
		}
	}
}

// RateLimiter implements a token bucket rate limiter
func RateLimiter(r float64, b int) middleware.Middleware {
	limiter := rate.NewLimiter(rate.Limit(r), b)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if !limiter.Allow() {
				return nil, errors.New(429, "TOO_MANY_REQUESTS", "rate limit exceeded")
			}
			return handler(ctx, req)
		}
	}
}

// CSRF middleware for CSRF protection
func CSRF() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				// Only check POST/PUT/DELETE/PATCH methods
				if tr.Kind() == transport.KindHTTP {
					method := tr.RequestHeader().Get("Method")
					if method != "GET" && method != "HEAD" && method != "OPTIONS" {
						token := tr.RequestHeader().Get("X-CSRF-Token")
						cookie := tr.RequestHeader().Get("Cookie")
						if !validateCSRFToken(token, cookie) {
							return nil, errors.New(403, "INVALID_CSRF_TOKEN", "invalid or missing CSRF token")
						}
					}
				}
			}
			return handler(ctx, req)
		}
	}
}

func validateCSRFToken(token, cookie string) bool {
	// 实际项目中应该从cookie中提取token并与请求中的token比较
	return token != "" && cookie != ""
}

// GenerateCSRFToken generates a new CSRF token
func GenerateCSRFToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

// XSS middleware for XSS protection
func XSS() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				if tr.Kind() == transport.KindHTTP {
					// Sanitize input parameters
					sanitizeHeaders(tr.RequestHeader())
				}
			}
			return handler(ctx, req)
		}
	}
}

func sanitizeHeaders(header http.Header) {
	for key, values := range header {
		for i, value := range values {
			values[i] = sanitizeString(value)
		}
		header[key] = values
	}
}

func sanitizeString(s string) string {
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}
