package job

import (
	"fmt"

	jobv1 "github.com/ydssx/kratos-kit/api/job/v1"
)

var (
	// CronJobMap 定时任务注册
	AdminCronJobMap = map[string]jobv1.AdminJob{}

	// JobHandlerMap 任务处理函数注册
	AdminJobHandlerMap = map[jobv1.AdminJob]jobHandler{}
)

func ValidateAdminTask(jobType jobv1.AdminJob) error {
	if _, ok := AdminJobHandlerMap[jobType]; !ok {
		return fmt.Errorf("the cron job [%s] does not have any registered handlers", jobType.String())
	}
	return nil
}
