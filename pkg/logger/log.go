package logger

import (
	"context"
	"fmt"

	"github.com/ydssx/kratos-kit/constants"
	"github.com/ydssx/kratos-kit/pkg/middleware/kratos"

	"github.com/go-kratos/kratos/v2/log"
)

var LogInject = func(ctx context.Context) (data []interface{}) {
	if traceID := kratos.TraceIDFromContext(ctx); traceID != "" {
		data = append(data, constants.LogKeyTraceID, traceID)
	}
	if operation := kratos.OperationFromContext(ctx); operation != "" {
		data = append(data, constants.LogKeyOperation, operation)
	}
	if claims := kratos.GetClaims(ctx); claims != nil {
		if claims.Uid != 0 {
			data = append(data, constants.LogKeyUserID, claims.Uid)
		}
	}
	return data
}

func Info(ctx context.Context, msg ...interface{}) {
	kv := []interface{}{log.DefaultMessageKey, fmt.Sprint(msg...)}
	DefaultLogger.Log(log.LevelInfo, append(kv, LogInject(ctx)...))
}

func Infof(ctx context.Context, format string, msg ...interface{}) {
	kv := []interface{}{log.DefaultMessageKey, fmt.Sprintf(format, msg...)}
	DefaultLogger.Log(log.LevelInfo, append(kv, LogInject(ctx)...))
}

func Error(ctx context.Context, msg ...interface{}) {
	kv := []interface{}{log.DefaultMessageKey, fmt.Sprint(msg...)}
	DefaultLogger.Log(log.LevelError, append(kv, LogInject(ctx)...))
}

func Errorf(ctx context.Context, format string, msg ...interface{}) {
	kv := []interface{}{log.DefaultMessageKey, fmt.Sprintf(format, msg...)}
	DefaultLogger.Log(log.LevelError, append(kv, LogInject(ctx)...))
}

func Fatalf(ctx context.Context, format string, msg ...interface{}) {
	kv := []interface{}{log.DefaultMessageKey, fmt.Sprintf(format, msg...)}
	DefaultLogger.Log(log.LevelFatal, append(kv, LogInject(ctx)...))
}

func Fatal(ctx context.Context, msg ...interface{}) {
	kv := []interface{}{log.DefaultMessageKey, fmt.Sprint(msg...)}
	DefaultLogger.Log(log.LevelFatal, append(kv, LogInject(ctx)...))
}

func Log(ctx context.Context, level log.Level, msg string, keyvals ...interface{}) {
	kv := []interface{}{log.DefaultMessageKey, msg}
	kv = append(kv, LogInject(ctx)...)
	DefaultLogger.Log(level, append(kv, keyvals...))
}

func Debug(ctx context.Context, msg ...interface{}) {
	kv := []interface{}{log.DefaultMessageKey, fmt.Sprint(msg...)}
	DefaultLogger.Log(log.LevelDebug, append(kv, LogInject(ctx)...))
}

func Debugf(ctx context.Context, format string, msg ...interface{}) {
	kv := []interface{}{log.DefaultMessageKey, fmt.Sprintf(format, msg...)}
	DefaultLogger.Log(log.LevelDebug, append(kv, LogInject(ctx)...))
}

func Warn(ctx context.Context, msg ...interface{}) {
	kv := []interface{}{log.DefaultMessageKey, fmt.Sprint(msg...)}
	DefaultLogger.Log(log.LevelWarn, append(kv, LogInject(ctx)...))
}

func Warnf(ctx context.Context, format string, msg ...interface{}) {
	kv := []interface{}{log.DefaultMessageKey, fmt.Sprintf(format, msg...)}
	DefaultLogger.Log(log.LevelWarn, append(kv, LogInject(ctx)...))
}
