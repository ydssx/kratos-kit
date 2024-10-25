package kratos

import (
	"context"
	"fmt"
	"time"

	"github.com/ydssx/kratos-kit/pkg/errors"

	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/uuid"
)

func TraceServer() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				code      int32
				reason    string
				kind      string
				operation string
			)

			startTime := time.Now()
			traceID := uuid.NewString()
			ctx = NewTraceIDContext(ctx, traceID)
			if info, ok := http.RequestFromServerContext(ctx); ok {
				operation = info.URL.Path
				// info.Response.Header.Set("traceID", traceID)
			}
			ctx = NewOperationContext(ctx, operation)
			ts, _ := transport.FromServerContext(ctx)
			ts.ReplyHeader().Set("traceID", traceID)

			klog := log.Context(ctx)
			klog.Log(log.LevelInfo,
				"path", operation,
				"kind", "server",
				"component", kind,
				"args", extractArgs(req),
			)
			reply, err = handler(ctx, req)
			if err == nil {
				return
			}
			if _, ok := err.(*errors.UserError); ok { // 给用户的提示性错误
				return
			}
			if ke, ok := err.(*kerrors.Error); ok && ke.Code != kerrors.UnknownCode { // http状态码类错误，比如401,404等
				return
			}

			klog.Log(log.LevelError,
				"msg", err.Error(),
				"path", operation,
				"latency", time.Since(startTime).String(),
				"code", code,
				"reason", reason,
			)
			return
		}
	}
}

type traceIDKey struct{}

func NewTraceIDContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, id)
}

func TraceIDFromContext(ctx context.Context) string {
	id, ok := ctx.Value(traceIDKey{}).(string)
	if !ok {
		id = ""
	}
	return id
}

func TraceID() log.Valuer {
	return func(ctx context.Context) interface{} {
		return TraceIDFromContext(ctx)
	}
}

type operationKey struct{}

func NewOperationContext(ctx context.Context, operation string) context.Context {
	return context.WithValue(ctx, operationKey{}, operation)
}

func OperationFromContext(ctx context.Context) string {
	op, ok := ctx.Value(operationKey{}).(string)
	if !ok {
		op = ""
	}
	return op
}

// extractArgs returns the string of the req
func extractArgs(req interface{}) string {
	if stringer, ok := req.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%+v", req)
}

// extractError returns the string of the error
func extractError(err error) (log.Level, string) {
	if err != nil {
		return log.LevelError, err.Error()
	}
	return log.LevelInfo, ""
}
