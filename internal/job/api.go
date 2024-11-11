package job

import (
	"context"
	"fmt"
	"os"

	jobv1 "github.com/ydssx/kratos-kit/api/job/v1"
	"github.com/ydssx/kratos-kit/constants"
	"github.com/ydssx/kratos-kit/internal/biz"
	"github.com/ydssx/kratos-kit/pkg/logger"
	"github.com/ydssx/kratos-kit/pkg/util"

	"github.com/hibiken/asynq"
)

type jobHandler func(ctx context.Context, t *asynq.Task) error

var (
	// 定时任务注册
	CronJobMap = map[string]jobv1.JobType{
		// "0 0 * * *":   jobv1.JobType_RESET_DAILY_CHARACTERS,
		// "*/2 * * * *": jobv1.JobType_TASK_TIMEOUT_REFUND, // 每2分钟执行一次
		"10 0 * * *": jobv1.JobType_CLEAN_OLD_LOG_FILES,
	}

	// 任务处理函数注册
	JobHandlerMap = map[jobv1.JobType]jobHandler{
		jobv1.JobType_CLEAN_OLD_LOG_FILES: ClearLogFile,
	}
)

func ValidateTask(jobType jobv1.JobType) error {
	if _, ok := JobHandlerMap[jobType]; !ok {
		return fmt.Errorf("the cron job [%s] does not have any registered handlers", jobType.String())
	}
	return nil
}

// ====================================================================================
//                        以下为定时任务和队列任务处理函数
// ====================================================================================

// 定时任务：清理用户上传文件
func ClearUserUploadFile(ctx context.Context, _ *asynq.Task) error {
	if err := biz.UsecaseSetFromContext(ctx).UploadBiz.CleanUploadFile(ctx); err != nil {
		logger.Errorf(ctx, "清理用户上传文件异常:%s", err.Error())
	}
	return nil
}

// 定时任务：清理日志文件
func ClearLogFile(ctx context.Context, _ *asynq.Task) error {
	err := util.DeleteOldFiles(os.Getenv(string(constants.EnvKeyLogPath)), 30)
	if err != nil {
		logger.Errorf(ctx, "清理日志文件异常:%s", err.Error())
	}
	return nil
}
