package admin

import (
	"context"

	"github.com/ydssx/kratos-kit/common/conf"
	"github.com/ydssx/kratos-kit/internal/biz"
	"github.com/ydssx/kratos-kit/internal/job"
	"github.com/ydssx/kratos-kit/pkg/logger"
	"github.com/ydssx/kratos-kit/pkg/queue"
)

type JobServer struct {
	client *queue.Client
}

func NewJobServer(c *conf.Bootstrap, serviceSet *biz.AdminUseCase) *JobServer {
	cfg := &queue.Config{
		RedisAddr:     c.Data.JobRedis.Addr,
		RedisPassword: c.Data.JobRedis.Password,
		RedisDB:       int(c.Data.JobRedis.Db),
		Concurrency:   int(c.Asynq.Concurrency),
		ReadTimeout:   c.Data.JobRedis.ReadTimeout.AsDuration(),
		WriteTimeout:  c.Data.JobRedis.WriteTimeout.AsDuration(),
		BaseContext: func() context.Context {
			return biz.WithAdminUseCase(context.Background(), serviceSet)
		},
	}

	client, err := queue.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	// 注册管理员任务处理器
	registerJobHandler(client)
	// 注册管理员定时任务
	registerCronJob(client)

	return &JobServer{client: client}
}

// Start starts the JobServer
func (j *JobServer) Start(ctx context.Context) error {
	return j.client.Start()
}

// Stop stops the JobServer gracefully
func (j *JobServer) Stop(ctx context.Context) error {
	err := j.client.Close()
	if err != nil {
		logger.Errorf(ctx, "failed to stop admin job server: %v", err)
		return err
	}
	logger.Info(ctx, "admin job server stopped")
	return nil
}

// registerJobHandler registers all admin job handlers defined in AdminJobHandlerMap to the client
func registerJobHandler(client *queue.Client) {
	for k, v := range job.AdminJobHandlerMap {
		client.RegisterHandler(k.String(), queue.HandleFunc(v))
	}
}

// registerCronJob registers all admin cron jobs defined in AdminCronJobMap
func registerCronJob(client *queue.Client) {
	for spec, jobType := range job.AdminCronJobMap {
		err := job.ValidateAdminTask(jobType)
		if err != nil {
			panic(err)
		}

		task := &queue.Task{
			TypeName: jobType.String(),
			Payload:  nil,
		}

		err = client.EnqueuePeriodicTask(context.Background(), task, spec)
		if err != nil {
			panic(err)
		}
	}
}
