package server

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

func NewJobServer(c *conf.Bootstrap, serviceSet *biz.UsecaseSet) *JobServer {
	cfg := &queue.Config{
		RedisAddr:     c.Data.Redis.Addr,
		RedisPassword: c.Data.Redis.Password,
		RedisDB:       int(c.Data.Redis.Db),
		Concurrency:   int(c.Asynq.Concurrency),
		ReadTimeout:   c.Data.Redis.ReadTimeout.AsDuration(),
		WriteTimeout:  c.Data.Redis.WriteTimeout.AsDuration(),
		BaseContext: func() context.Context {
			return biz.NewContextWithUsecaseSet(context.Background(), serviceSet)
		},
	}

	client, err := queue.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	// 注册任务处理器
	registerJobHandler(client)
	// 注册定时任务
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
		logger.Errorf(ctx, "failed to stop job server: %v", err)
		return err
	}
	logger.Info(ctx, "job server stopped")
	return nil
}

// registerJobHandler registers all job handlers defined in jobHandlerMap to the client
func registerJobHandler(client *queue.Client) {
	for k, v := range job.JobHandlerMap {
		client.RegisterHandler(k.String(), queue.HandleFunc(v))
	}
}

// registerCronJob registers all cron jobs defined in cronJobMap
func registerCronJob(client *queue.Client) {
	for spec, jobType := range job.CronJobMap {
		err := job.ValidateTask(jobType)
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
