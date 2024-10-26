package middleware

import (
	"context"
	"fmt"
	stdhttp "net/http"
	"time"

	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/uuid"
	"github.com/ydssx/kratos-kit/pkg/errors"
	"github.com/ydssx/kratos-kit/pkg/middleware/kratos"
)

// TraceServer 返回一个跟踪中间件
func TraceServer() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			// 初始化跟踪上下文
			ctx, operation := initTraceContext(ctx)
			startTime := time.Now()

			// 记录请求日志
			logRequest(ctx, operation, req)

			// 处理请求
			reply, err = handler(ctx, req)

			// 处理特殊错误
			if shouldSkipError(err) {
				return
			}

			// 提取错误信息
			code, reason := extractErrorInfo(err)
			if code == stdhttp.StatusUnauthorized {
				return
			}

			// 记录响应日志
			logResponse(ctx, operation, startTime, code, reason, err)

			return
		}
	}
}

// initTraceContext 初始化跟踪上下文
func initTraceContext(ctx context.Context) (context.Context, string) {
	traceID := uuid.NewString()
	ctx = kratos.NewTraceIDContext(ctx, traceID)

	var operation string
	if info, ok := http.RequestFromServerContext(ctx); ok {
		operation = info.URL.Path
	}
	ctx = kratos.NewOperationContext(ctx, operation)

	if ts, ok := transport.FromServerContext(ctx); ok {
		ts.ReplyHeader().Set("traceID", traceID)
	}

	return ctx, operation
}

// logRequest 记录请求日志
func logRequest(ctx context.Context, operation string, req interface{}) {
	klog := log.Context(ctx)
	klog.Log(log.LevelInfo,
		"event", "request_received",
		"path", operation,
		"kind", "server",
		"args", extractArgs(req),
	)
}

// logResponse 记录响应日志
func logResponse(ctx context.Context, operation string, startTime time.Time, code int32, reason string, err error) {
	level, msg, stack := extractError(err)
	latency := time.Since(startTime)

	klog := log.Context(ctx)
	klog.Log(level,
		"event", "request_completed",
		"path", operation,
		"latency", latency.String(),
		"code", code,
		"reason", reason,
		"error", msg,
		"stack", stack,
	)
}

// shouldSkipError 判断是否需要跳过错误处理
func shouldSkipError(err error) bool {
	if err == nil {
		return true
	}

	switch {
	case errors.Is(err, context.Canceled):
		return true
	case errors.IsUserError(err):
		return true
	case isHttpStatusError(err):
		return true
	}
	return false
}

// isHttpStatusError 判断是否为HTTP状态码错误
func isHttpStatusError(err error) bool {
	if ke, ok := err.(*kerrors.Error); ok && ke.Code != kerrors.UnknownCode {
		return true
	}
	return false
}

// extractErrorInfo 提取错误信息
func extractErrorInfo(err error) (code int32, reason string) {
	if se := kerrors.FromError(err); se != nil {
		return se.Code, se.Reason
	}
	return 0, ""
}

// extractArgs 提取请求参数
func extractArgs(req interface{}) string {
	if stringer, ok := req.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%+v", req)
}

// extractError 提取错误详情
func extractError(err error) (level log.Level, msg string, stack string) {
	if err != nil {
		return log.LevelError, err.Error(), fmt.Sprintf("%+v", err)
	}
	return log.LevelInfo, "", ""
}
