package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/ydssx/kratos-kit/constants"
	"github.com/ydssx/kratos-kit/pkg/webhook"
	"github.com/ydssx/kratos-kit/third_party/lumberjack"

	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Writer zapcore.WriteSyncer
	once   sync.Once
)

func NewZapLogger(callerSkip int, logPath string, level string, maxSize int, maxBackups int, maxAge int, compress bool, enableConsole bool) *zap.Logger {
	once.Do(func() {
		Writer = zapcore.AddSync(&lumberjack.Logger{
			LogPath:    logPath,
			Filename:   lumberjack.GetLogFileName(logPath), // 日志文件路径
			MaxSize:    maxSize,                            // 每个日志文件最大尺寸（MB）
			MaxBackups: maxBackups,                         // 保留旧日志文件的个数
			MaxAge:     maxAge,                             // 保留旧日志文件的最大天数
			Compress:   compress,                           // 是否压缩/归档旧文件
		})
		ws := []zapcore.WriteSyncer{Writer}
		if enableConsole {
			ws = append(ws, zapcore.AddSync(os.Stdout))
		}
		Writer = zapcore.NewMultiWriteSyncer(ws...)
	})

	return initLogger(callerSkip, level, Writer)
}

func initLogger(callerSkip int, level string, ws zapcore.WriteSyncer) *zap.Logger {
	// 配置 Encoder
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.MessageKey = ""
	// 创建 Core
	core := zapcore.NewCore(zapcore.NewJSONEncoder(cfg), ws, getLevel(level))

	// 创建 Logger
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(callerSkip), zap.AddStacktrace(zap.DPanicLevel), zap.Hooks(func(e zapcore.Entry) error {
		return nil
	}))

	return logger
}

var _ log.Logger = (*Logger)(nil)

var DefaultLogger = NewLogger(initLogger(2, "info", os.Stdout))

func init() {
	// logger := log.NewFilter(log.GetLogger(), log.FilterLevel(log.LevelInfo))
	log.SetLogger(DefaultLogger)
}

type Logger struct {
	Zlog *zap.Logger
}

func NewLogger(zlog *zap.Logger) *Logger {
	return &Logger{zlog}
}

func (l *Logger) Log(level log.Level, keyvals ...interface{}) error {
	for i, v := range keyvals {
		if r, ok := v.([]interface{}); ok {
			keyvals = append(keyvals[:i], r...)
		}
	}
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.Zlog.Warn("Keyvalues must appear in pairs", zap.Any("keyvalues", keyvals))
		return nil
	}

	kvMap := make(map[string]interface{}, len(keyvals)/2)
	data := make([]zap.Field, 0, len(keyvals)/2)
	stack := make([]string, 0)
	for i := 0; i < len(keyvals); i += 2 {
		key := fmt.Sprint(keyvals[i])
		kvMap[key] = keyvals[i+1]
		if key == string(constants.LogKeyStack) {
			if v, ok := keyvals[i+1].([]string); ok {
				stack = v
				data = append(data, zap.String(string(constants.LogKeyStack), strings.Join(v, "\n")))
			} else {
				stack = []string{fmt.Sprintf("%s", keyvals[i+1])}
				data = append(data, zap.Any(key, keyvals[i+1]))
			}
		} else {
			data = append(data, zap.Any(key, keyvals[i+1]))
		}
	}

	ce := l.Zlog.Check(zapcore.Level(level), "")

	if level >= log.LevelError {
		message := fmt.Sprintf(msg,
			os.Getenv(string(constants.EnvKeyEnv)),
			level.String(),
			ce.Time.Format("2006-01-02 15:04:05"),
			kvMap[log.DefaultMessageKey],
			kvMap[string(constants.LogKeyTraceID)],
			kvMap[string(constants.LogKeyOperation)],
			ce.Caller.TrimmedPath(),
			strings.Join(stack, "\n"))

		go webhook.NewWebhook(os.Getenv(string(constants.EnvKeyDingDingWebhook))).SendMessage(message)
	}

	if ce != nil {
		ce.Write(data...)
	}
	return nil
}

func (l *Logger) Sync() error {
	return l.Zlog.Sync()
}

func (l *Logger) Close() error {
	return l.Sync()
}

func getLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		// 默认返回 InfoLevel
		return zapcore.InfoLevel
	}
}

var msg = `环境：%s
级别：%s
时间：%s
日志：%s
traceID：%s
path：%s
caller：%s
堆栈：%s
`
