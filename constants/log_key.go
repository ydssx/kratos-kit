package constants

type LogKey string

const (
	LogKeyTraceID   LogKey = "traceID" // 链路跟踪ID
	LogKeyOperation LogKey = "path"    // 操作路径
	LogKeyStack     LogKey = "stack"   // 堆栈信息
	LogKeyUserID    LogKey = "userID"  // 用户ID
)
