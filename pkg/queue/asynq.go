package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
	"github.com/ydssx/kratos-kit/pkg/errors"
	"github.com/ydssx/kratos-kit/pkg/logger"
)

// Client wraps asynq client
type Client struct {
	client    *asynq.Client
	srv       *asynq.Server
	mux       *asynq.ServeMux
	scheduler *asynq.Scheduler
}

// Config holds the configuration for asynq client
type Config struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration

	// Queue configurations
	Concurrency int
	RetryLimit  int

	BaseContext func() context.Context
}

// NewClient creates a new asynq client
func NewClient(cfg *Config) (*Client, error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}

	redisOpt := asynq.RedisClientOpt{
		Addr:         cfg.RedisAddr,
		Password:     cfg.RedisPassword,
		DB:           cfg.RedisDB,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	client := asynq.NewClient(redisOpt)
	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: cfg.Concurrency,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			BaseContext:  cfg.BaseContext,
			ErrorHandler: asynq.ErrorHandlerFunc(reportError),
			RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
				return time.Duration(n) * time.Minute
			},
		},
	)
	mux := asynq.NewServeMux()
	scheduler := asynq.NewScheduler(
		redisOpt,
		&asynq.SchedulerOpts{Location: time.Local},
	)

	return &Client{
		client:    client,
		srv:       srv,
		mux:       mux,
		scheduler: scheduler,
	}, nil
}

// Close closes the client
func (c *Client) Close() error {
	if c.client != nil {
		c.client.Close()
	}
	if c.srv != nil {
		c.srv.Stop()
		c.srv.Shutdown()
	}
	if c.scheduler != nil {
		c.scheduler.Shutdown()
	}
	return nil
}

// Task represents an async task
type Task struct {
	TypeName string
	Payload  interface{}
}

// EnqueueTask enqueues a task
func (c *Client) EnqueueTask(ctx context.Context, task *Task, opts ...asynq.Option) error {
	payload, err := json.Marshal(task.Payload)
	if err != nil {
		return errors.Errorf("failed to marshal task payload: %v", err)
	}

	t := asynq.NewTask(task.TypeName, payload)
	_, err = c.client.EnqueueContext(ctx, t, opts...)
	if err != nil {
		return errors.Errorf("failed to enqueue task: %v", err)
	}

	return nil
}

// HandleFunc represents a task handler function
type HandleFunc func(context.Context, *asynq.Task) error

// RegisterHandler registers a task handler
func (c *Client) RegisterHandler(taskType string, handler HandleFunc) {
	c.mux.HandleFunc(taskType, handler)
}

// Start starts the task processor
func (c *Client) Start() error {
	err := c.srv.Start(c.mux)
	if err != nil {
		return errors.Errorf("failed to start task processor: %v", err)
	}
	err = c.scheduler.Start()
	if err != nil {
		return errors.Errorf("failed to start scheduler: %v", err)
	}
	return nil
}

// EnqueueTaskWithDelay enqueues a task with delay
func (c *Client) EnqueueTaskWithDelay(ctx context.Context, task *Task, delay time.Duration) error {
	return c.EnqueueTask(ctx, task, asynq.ProcessIn(delay))
}

// EnqueuePeriodicTask enqueues a periodic task
func (c *Client) EnqueuePeriodicTask(ctx context.Context, task *Task, spec string) error {
	payload, err := json.Marshal(task.Payload)
	if err != nil {
		return errors.Errorf("failed to marshal task payload: %v", err)
	}

	t := asynq.NewTask(task.TypeName, payload)
	_, err = c.scheduler.Register(spec, t)
	if err != nil {
		return errors.Errorf("failed to register periodic task: %v", err)
	}

	return nil
}

func reportError(ctx context.Context, task *asynq.Task, err error) {
	logger.Errorf(ctx, "执行任务失败,task_type:%s ,err: %v", task.Type(), err)
}
