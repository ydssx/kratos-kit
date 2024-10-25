package admin

import (
	"context"
	"time"

	"github.com/ydssx/kratos-kit/internal/biz"
	"github.com/ydssx/kratos-kit/internal/job"

	common2 "github.com/ydssx/kratos-kit/common"
	"github.com/ydssx/kratos-kit/common/conf"
	"github.com/ydssx/kratos-kit/pkg/logger"

	"github.com/hibiken/asynq"
	"golang.org/x/sync/errgroup"
)

type JobServer struct {
	sr  *asynq.Server
	sd  *asynq.Scheduler
	mux *asynq.ServeMux
}

func NewJobServer(c *conf.Bootstrap, serviceSet *biz.AdminUseCase) *JobServer {
	opt := common2.InitJobRedisOpt(c)

	server := asynq.NewServer(opt, asynq.Config{
		Concurrency:  int(c.Asynq.Concurrency),
		ErrorHandler: asynq.ErrorHandlerFunc(reportError),
		BaseContext: func() context.Context {
			return biz.WithAdminUseCase(context.Background(), serviceSet)
		},
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
		StrictPriority: c.Asynq.StrictPriority,
	})

	mux := asynq.NewServeMux()
	registerJobHandler(mux)

	scheduler := asynq.NewScheduler(opt, &asynq.SchedulerOpts{Location: time.Local})
	registerCronJob(scheduler)

	return &JobServer{sr: server, mux: mux, sd: scheduler}
}

// Start starts the JobServer, including the scheduler and the server.
// It starts the scheduler and the server concurrently using goroutines.
// It returns an error if any of the concurrent operations fail.
func (j *JobServer) Start(ctx context.Context) error {
	// Start the scheduler and the server concurrently.
	eg := errgroup.Group{}

	// Start the scheduler.
	eg.Go(j.sd.Start)

	// Start the server.
	eg.Go(func() error {
		// Start the server with the registered job handlers.
		return j.sr.Start(j.mux)
	})

	// Wait for all the concurrent operations to complete.
	// If any of the operations fail, the function returns the error.
	return eg.Wait()
}

// Stop 停止 JobServer,包括停止调度器和服务器。
// 依次调用服务器和调度器的 Shutdown 方法进行优雅停止。
func (j *JobServer) Stop(ctx context.Context) error {
	j.sr.Stop()
	j.sr.Shutdown()
	j.sd.Shutdown()
	logger.Info(ctx, "job server stopped")
	return nil
}

func reportError(ctx context.Context, task *asynq.Task, err error) {
	logger.Errorf(ctx, "执行任务失败,task_type:%s ,err: %v", task.Type(), err)
}

// registerJobHandler 注册 jobHandlerMap 中定义的所有 job 的处理函数到 ServeMux。
// 它会遍历 jobHandlerMap,并为每个 job 注册对应的处理函数到 mux。
// mux 会根据请求中的 job name 来路由到相应的处理函数。
func registerJobHandler(mux *asynq.ServeMux) {
	for k, v := range job.AdminJobHandlerMap {
		mux.HandleFunc(k.String(), v)
	}
}

// registerCronJob 注册定时任务的处理函数。
// 它会遍历 cronJobMap 中定义的所有定时任务,并在调度器 sd 中注册对应的处理函数。
// 如果某个定时任务在 jobHandlerMap 中没有找到对应的处理函数,会 panic。
// 注册成功后,定时任务会按照 cronJobMap 中定义的时间表定期执行。
func registerCronJob(sd *asynq.Scheduler) {
	for k, jobType := range job.AdminCronJobMap {
		err := job.ValidateAdminTask(jobType)
		if err != nil {
			panic(err)
		}
		_, err = sd.Register(k, asynq.NewTask(jobType.String(), nil))
		if err != nil {
			panic(err)
		}
	}
}
