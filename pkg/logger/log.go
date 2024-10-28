package logger

import (
	"context"
	"fmt"

	"github.com/ydssx/kratos-kit/constants"
	"github.com/ydssx/kratos-kit/pkg/middleware/kratos"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

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

// 优化日志方法,减少代码重复
func logWithLevel(ctx context.Context, level log.Level, format string, args ...interface{}) {
	var msg string
	if format == "" {
		msg = fmt.Sprint(args...)
	} else {
		msg = fmt.Sprintf(format, args...)
	}

	kv := []interface{}{log.DefaultMessageKey, msg}
	kv = append(kv, LogInject(ctx)...)
	DefaultLogger.Log(level, kv...)
}

// 简化各个日志级别的方法
func Info(ctx context.Context, args ...interface{}) {
	logWithLevel(ctx, log.LevelInfo, "", args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	logWithLevel(ctx, log.LevelInfo, format, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	logWithLevel(ctx, log.LevelError, "", args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	logWithLevel(ctx, log.LevelError, format, args...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	logWithLevel(ctx, log.LevelFatal, format, args...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	logWithLevel(ctx, log.LevelFatal, "", args...)
}

func Log(ctx context.Context, level log.Level, msg string, keyvals ...interface{}) {
	kv := []interface{}{log.DefaultMessageKey, msg}
	kv = append(kv, LogInject(ctx)...)
	DefaultLogger.Log(level, append(kv, keyvals...)...)
}

func Debug(ctx context.Context, args ...interface{}) {
	logWithLevel(ctx, log.LevelDebug, "", args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	logWithLevel(ctx, log.LevelDebug, format, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	logWithLevel(ctx, log.LevelWarn, "", args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	logWithLevel(ctx, log.LevelWarn, format, args...)
}

func WithFields(fields ...zap.Field) *Logger {
	return DefaultLogger.WithFields(fields...)
}

// 添加一些实用的辅助方法
func WithContext(ctx context.Context) *ContextLogger {
	return &ContextLogger{
		logger: DefaultLogger,
		ctx:    ctx,
	}
}

type ContextLogger struct {
	logger *Logger
	ctx    context.Context
}

// 扩展 ContextLogger 的方法
func (l *ContextLogger) Info(msg string, fields ...zap.Field) {
	l.logger.Log(log.LevelInfo, l.getLogFields(msg, fields...)...)
}

func (l *ContextLogger) Error(msg string, fields ...zap.Field) {
	l.logger.Log(log.LevelError, l.getLogFields(msg, fields...)...)
}

func (l *ContextLogger) Debug(msg string, fields ...zap.Field) {
	l.logger.Log(log.LevelDebug, l.getLogFields(msg, fields...)...)
}

func (l *ContextLogger) Warn(msg string, fields ...zap.Field) {
	l.logger.Log(log.LevelWarn, l.getLogFields(msg, fields...)...)
}

// 添加更多实用方法
func (l *ContextLogger) WithError(err error) *ContextLogger {
	if err == nil {
		return l
	}
	return l.WithFields(zap.Error(err))
}

func (l *ContextLogger) WithValue(key string, value interface{}) *ContextLogger {
	return l.WithFields(zap.Any(key, value))
}

// 优化全局方法
func WithError(err error) *Logger {
	if err == nil {
		return DefaultLogger
	}
	return DefaultLogger.WithFields(zap.Error(err))
}

func WithValue(key string, value interface{}) *Logger {
	return DefaultLogger.WithFields(zap.Any(key, value))
}

// 添加通用的日志字段处理方法
func (l *ContextLogger) getLogFields(msg string, fields ...zap.Field) []interface{} {
	data := l.logger.opts.logInject(l.ctx)
	kv := []interface{}{log.DefaultMessageKey, msg}
	kv = append(kv, data...)

	// 处理额外的字段
	for _, field := range fields {
		kv = append(kv, field.Key, field.Interface)
	}

	return kv
}

// 添加带格式化的方法
func (l *ContextLogger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

func (l *ContextLogger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

func (l *ContextLogger) Debugf(format string, args ...interface{}) {
	l.Debug(fmt.Sprintf(format, args...))
}

func (l *ContextLogger) Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}

// 添加 WithFields 方法
func (l *ContextLogger) WithFields(fields ...zap.Field) *ContextLogger {
	return &ContextLogger{
		logger: l.logger.WithFields(fields...),
		ctx:    l.ctx,
	}
}

// 添加缺失的方法
func (l *ContextLogger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Log(log.LevelFatal, l.getLogFields(msg, fields...)...)
}

func (l *ContextLogger) Fatalf(format string, args ...interface{}) {
	l.Fatal(fmt.Sprintf(format, args...))
}

// 添加日志级别检查
func (l *ContextLogger) IsDebugEnabled() bool {
	return l.logger.Zlog.Core().Enabled(zapcore.DebugLevel)
}

func (l *ContextLogger) IsInfoEnabled() bool {
	return l.logger.Zlog.Core().Enabled(zapcore.InfoLevel)
}

func (l *ContextLogger) IsWarnEnabled() bool {
	return l.logger.Zlog.Core().Enabled(zapcore.WarnLevel)
}

func (l *ContextLogger) IsErrorEnabled() bool {
	return l.logger.Zlog.Core().Enabled(zapcore.ErrorLevel)
}

// 其他方法类似...
