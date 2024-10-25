package middleware

import (
	"context"
	"fmt"
	stdhttp "net/http"
	"time"

	"github.com/ydssx/kratos-kit/pkg/errors"
	"github.com/ydssx/kratos-kit/pkg/middleware/kratos"

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
			ctx = kratos.NewTraceIDContext(ctx, traceID)
			if info, ok := http.RequestFromServerContext(ctx); ok {
				operation = info.URL.Path
			}
			ctx = kratos.NewOperationContext(ctx, operation)
			ts, _ := transport.FromServerContext(ctx)
			ts.ReplyHeader().Set("traceID", traceID)

			klog := log.Context(ctx)
			klog.Log(log.LevelInfo,
				"path", operation,
				"kind", "server",
				"component", kind,
				"args", extractArgs(req, ts.Operation()),
			)

			reply, err = handler(ctx, req)
			if _, ok := err.(*errors.UserError); ok {
				return
			}
			if ke, ok := err.(*kerrors.Error); ok && ke.Code != kerrors.UnknownCode { // http状态码类错误，比如401,404等
				return
			}
			if errors.Is(err, context.Canceled) {
				return
			}

			if se := kerrors.FromError(err); se != nil {
				code = se.Code
				reason = se.Reason
				if code == stdhttp.StatusUnauthorized {
					return
				}
			}
			level, msg, stack := extractError(err)
			klog.Log(level,
				"msg", msg,
				"path", operation,
				"latency", time.Since(startTime).String(),
				"code", code,
				"reason", reason,
				"stack", stack,
			)
			return
		}
	}
}

// extractArgs returns the string of the req
func extractArgs(req interface{}, _ string) string {
	if stringer, ok := req.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%+v", req)
}

// extractError returns the string of the error
func extractError(err error) (lvl log.Level, msg string, stack string) {
	if err != nil {
		return log.LevelError, err.Error(), fmt.Sprintf("%+v", err)
	}
	return log.LevelInfo, "", ""
}
