package logger

import (
	"context"
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

func NewZapLogger(opts ...Option) *zap.Logger {
	options := defaultOptions()
	for _, o := range opts {
		o(options)
	}

	once.Do(func() {
		initWriter(options)
	})

	return initLogger(options.callerSkip, options.level, Writer)
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

// 优化 Logger 结构体,添加配置项
type Logger struct {
	Zlog *zap.Logger
	opts *Options
}

// 添加 Options 结构体用于配置
type Options struct {
	webhookURL    string
	enableWebhook bool
	logInject     func(context.Context) []interface{}
	callerSkip    int
	level         string
	logPath       string // 新增
	maxSize       int    // 新增
	maxBackups    int    // 新增
	maxAge        int    // 新增
	compress      bool   // 新增
	enableConsole bool   // 新增
}

// 添加 Option 函数类型用于配置
type Option func(*Options)

// 创建默认配置
func defaultOptions() *Options {
	return &Options{
		enableWebhook: false,
		logInject:     func(context.Context) []interface{} { return nil },
		callerSkip:    1,
		level:         "info",
		logPath:       "./logs",
		maxSize:       100,
		maxBackups:    10,
		maxAge:        7,
		compress:      true,
		enableConsole: true,
	}
}

// 创建 Logger 时支持可选配置
func NewLogger(zlog *zap.Logger, opts ...Option) *Logger {
	options := defaultOptions()
	for _, o := range opts {
		o(options)
	}
	return &Logger{
		Zlog: zlog,
		opts: options,
	}
}

// 优化 Log 方法,提取公共逻辑
func (l *Logger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.Zlog.Warn("Keyvalues must appear in pairs", zap.Any("keyvalues", keyvals))
		return fmt.Errorf("invalid keyvals length: %d", len(keyvals))
	}

	// 提取并处理日志字段
	fields, stack, kvMap := l.processLogFields(keyvals)

	// 检查日志级别
	ce := l.Zlog.Check(zapcore.Level(level), "")
	if ce == nil {
		return nil
	}

	// 发送webhook消息
	if level >= log.LevelError && l.opts.enableWebhook {
		l.sendWebhookMessage(ce, kvMap, stack)
	}

	ce.Write(fields...)
	return nil
}

// 提取处理日志字段的逻辑
func (l *Logger) processLogFields(keyvals []interface{}) ([]zap.Field, []string, map[string]interface{}) {
	kvMap := make(map[string]interface{}, len(keyvals)/2)
	fields := make([]zap.Field, 0, len(keyvals)/2)
	var stack []string

	for i := 0; i < len(keyvals); i += 2 {
		key := fmt.Sprint(keyvals[i])
		value := keyvals[i+1]
		kvMap[key] = value

		if key == string(constants.LogKeyStack) {
			stack = l.extractStack(value)
			fields = append(fields, zap.String(string(constants.LogKeyStack), strings.Join(stack, "\n")))
		} else {
			fields = append(fields, zap.Any(key, value))
		}
	}

	return fields, stack, kvMap
}

func (l *Logger) Sync() error {
	return l.Zlog.Sync()
}

func (l *Logger) Close() error {
	return l.Sync()
}

func (l *Logger) WithFields(fields ...zap.Field) *Logger {
	return &Logger{
		Zlog: l.Zlog.With(fields...),
		opts: l.opts,
	}
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.Zlog.Info(msg, fields...)
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

// 添加 webhook 相关方法
func (l *Logger) sendWebhookMessage(ce *zapcore.CheckedEntry, kvMap map[string]interface{}, stack []string) {
	message := fmt.Sprintf(msg,
		os.Getenv(string(constants.EnvKeyEnv)),
		ce.Level.String(),
		ce.Time.Format("2006-01-02 15:04:05"),
		kvMap[log.DefaultMessageKey],
		kvMap[string(constants.LogKeyTraceID)],
		kvMap[string(constants.LogKeyOperation)],
		ce.Caller.TrimmedPath(),
		strings.Join(stack, "\n"))

	if l.opts.webhookURL != "" {
		go webhook.NewWebhook(l.opts.webhookURL).SendMessage(message)
	}
}

// 添加提取堆栈的方法
func (l *Logger) extractStack(value interface{}) []string {
	if v, ok := value.([]string); ok {
		return v
	}
	return []string{fmt.Sprintf("%v", value)}
}

// 添加配置选项
func WithWebhook(url string) Option {
	return func(o *Options) {
		o.webhookURL = url
		o.enableWebhook = true
	}
}

func WithLogInject(f func(context.Context) []interface{}) Option {
	return func(o *Options) {
		o.logInject = f
	}
}

// 添加更多实用的配置选项
func WithCallerSkip(skip int) Option {
	return func(o *Options) {
		o.callerSkip = skip
	}
}

func WithLevel(level string) Option {
	return func(o *Options) {
		o.level = level
	}
}

// 添加 WithContext 方法
func (l *Logger) WithContext(ctx context.Context) *ContextLogger {
	return &ContextLogger{
		logger: l,
		ctx:    ctx,
	}
}

// 添加缺失的配置选项方法
func WithLogPath(path string) Option {
	return func(o *Options) {
		o.logPath = path
	}
}

func WithMaxSize(size int) Option {
	return func(o *Options) {
		o.maxSize = size
	}
}

func WithMaxBackups(backups int) Option {
	return func(o *Options) {
		o.maxBackups = backups
	}
}

func WithMaxAge(age int) Option {
	return func(o *Options) {
		o.maxAge = age
	}
}

func WithCompress(compress bool) Option {
	return func(o *Options) {
		o.compress = compress
	}
}

func WithEnableConsole(enable bool) Option {
	return func(o *Options) {
		o.enableConsole = enable
	}
}

// 添加 Writer 初始化方法
func initWriter(opts *Options) {
	lumberjackLogger := &lumberjack.Logger{
		LogPath:    opts.logPath,
		Filename:   lumberjack.GetLogFileName(opts.logPath),
		MaxSize:    opts.maxSize,
		MaxBackups: opts.maxBackups,
		MaxAge:     opts.maxAge,
		Compress:   opts.compress,
	}

	ws := []zapcore.WriteSyncer{zapcore.AddSync(lumberjackLogger)}
	if opts.enableConsole {
		ws = append(ws, zapcore.AddSync(os.Stdout))
	}
	Writer = zapcore.NewMultiWriteSyncer(ws...)
}
