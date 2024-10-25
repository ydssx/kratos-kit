package constants

type EnvKey string

const (
	EnvKeyDingDingWebhook EnvKey = "DING_DING_WEBHOOK"
	EnvKeyLogPath         EnvKey = "LOG_PATH" // 日志路径
	EnvKeyEnv             EnvKey = "ENV"      // 当前环境
)
