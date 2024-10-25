package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/ydssx/kratos-kit/constants"
	"github.com/ydssx/kratos-kit/pkg/logger"

	"github.com/axiaoxin-com/goutils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gormlogger "gorm.io/gorm/logger"
)

var (
	// GormLoggerName gorm logger 名称
	GormLoggerName = "gorm"
	// GormLoggerCallerSkip caller skip
	GormLoggerCallerSkip = 2
)

// GormLogger 使用 zap 来打印 gorm 的日志
// 初始化时在内部的 logger 中添加 trace id 可以追踪 sql 执行记录
type GormLogger struct {
	// 日志级别
	logLevel zapcore.Level
	// 指定慢查询时间
	slowThreshold time.Duration
	// Trace 方法打印日志是使用的日志 level
	traceWithLevel zapcore.Level
}

var gormLogLevelMap = map[gormlogger.LogLevel]zapcore.Level{
	gormlogger.Info:  zap.InfoLevel,
	gormlogger.Warn:  zap.WarnLevel,
	gormlogger.Error: zap.ErrorLevel,
}

// LogMode 实现 gorm logger 接口方法
func (g GormLogger) LogMode(gormLogLevel gormlogger.LogLevel) gormlogger.Interface {
	zaplevel, exists := gormLogLevelMap[gormLogLevel]
	if !exists {
		zaplevel = zap.DebugLevel
	}
	newlogger := g
	newlogger.logLevel = zaplevel
	return &newlogger
}

// CtxLogger 创建打印日志的 ctxlogger
func (g GormLogger) CtxLogger(ctx context.Context) *zap.Logger {
	log := logger.DefaultLogger.Zlog.Named(GormLoggerName).WithOptions(zap.AddCallerSkip(GormLoggerCallerSkip), zap.AddCaller())
	kvs := logger.LogInject(ctx)
	for i := 0; i < len(kvs); i += 2 {
		log = log.With(zap.Any(string(kvs[i].(constants.LogKey)), kvs[i+1]))
	}
	return log
}

// Info 实现 gorm logger 接口方法
func (g GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if g.logLevel <= zap.InfoLevel {
		g.CtxLogger(ctx).Sugar().Infof(msg, data...)
	}
}

// Warn 实现 gorm logger 接口方法
func (g GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if g.logLevel <= zap.WarnLevel {
		g.CtxLogger(ctx).Sugar().Warnf(msg, data...)
	}
}

// Error 实现 gorm logger 接口方法
func (g GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if g.logLevel <= zap.ErrorLevel {
		g.CtxLogger(ctx).Sugar().Errorf(msg, data...)
	}
}

// Trace implements the gorm logger interface method
func (g GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	latency := time.Since(begin).Milliseconds()
	sql, rows := fc()
	sql = goutils.RemoveDuplicateWhitespace(sql, true)
	logger := g.CtxLogger(ctx)

	logFields := []zap.Field{
		zap.String("latency", fmt.Sprintf("%.3fms", float64(latency))),
		zap.Int64("rows", rows),
		zap.String("sql", sql),
	}

	switch {
	case err != nil:
		logFields = append(logFields, zap.Error(err))
		logger.Error("", logFields...)
	case g.slowThreshold != 0 && float64(latency) > float64(g.slowThreshold.Milliseconds()):
		logFields = append(logFields, zap.Float64("threshold", float64(g.slowThreshold.Milliseconds())))
		logger.Warn("slow sql", logFields...)
	default:
		logger.Debug("", logFields...)
	}
}

// NewGormLogger 返回带 zap logger 的 GormLogger
func NewGormLogger(logLevel zapcore.Level, traceWithLevel zapcore.Level, slowThreshold time.Duration) GormLogger {
	return GormLogger{
		logLevel:       logLevel,
		slowThreshold:  slowThreshold,
		traceWithLevel: traceWithLevel,
	}
}
