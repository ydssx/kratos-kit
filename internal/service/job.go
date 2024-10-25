package service

import (
	"context"

	jobv1 "github.com/ydssx/kratos-kit/api/job/v1"
	"github.com/ydssx/kratos-kit/internal/job"
	"github.com/ydssx/kratos-kit/pkg/errors"

	"github.com/hibiken/asynq"
)

type JobService struct {
	cli *asynq.Client
	ipt *asynq.Inspector
}

func NewJobService(cli *asynq.Client, ipt *asynq.Inspector) *JobService {
	return &JobService{cli: cli, ipt: ipt}
}

// Enqueue 将任务添加到队列中。
// 它接受 EnqueueRequest 并返回 EnqueueResponse。
// 它会根据请求设置 asynq 选项,如重试次数、延迟执行时间等。
// 然后使用 asynq client 将任务信息提交到队列。
// 如果成功,返回任务 ID;如果失败,返回错误。
func (j *JobService) Enqueue(ctx context.Context, req *jobv1.EnqueueRequest) (*jobv1.EnqueueResponse, error) {
	opts := []asynq.Option{asynq.MaxRetry(0)}
	if req.RetryTime > 0 {
		opts = append(opts, asynq.MaxRetry(int(req.RetryTime)))
	}
	if req.ProcessAt.IsValid() {
		opts = append(opts, asynq.ProcessAt(req.ProcessAt.AsTime()))
	}
	if req.ProcessIn.IsValid() {
		opts = append(opts, asynq.ProcessIn(req.ProcessIn.AsDuration()))
	}
	if req.Retention.IsValid() {
		opts = append(opts, asynq.Retention(req.Retention.AsDuration()))
	}
	if err := job.ValidateTask(req.JobType); err != nil {
		return nil, errors.Wrapf(err, "failed to validate task: %v", err)
	}
	taskInfo, err := j.cli.EnqueueContext(ctx, asynq.NewTask(req.JobType.String(), req.Payload), opts...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to send task: %v", err)
	}
	return &jobv1.EnqueueResponse{TaskId: taskInfo.ID}, nil
}

// QueryTasks 根据任务 ID 查询任务详情。
// 它接受 QueryTasksRequest 并返回 QueryTasksResponse。
// 对于每个任务 ID,它会从 asynq Inspector 获取任务详情。
// 然后将任务 ID、状态和结果组装到响应中返回。
// 如果查询失败,返回错误。
func (j *JobService) QueryTasks(ctx context.Context, req *jobv1.QueryTasksRequest) (resp *jobv1.QueryTasksResponse, err error) {
	resp = new(jobv1.QueryTasksResponse)
	for _, taskId := range req.TaskIds {
		taskInfo, err := j.ipt.GetTaskInfo("default", taskId)
		if err != nil {
			return nil, err
		}
		resp.Tasks = append(resp.Tasks, &jobv1.QueryTasksResponse_TaskInfo{
			TaskId: taskId,
			Result: taskInfo.Result,
			Status: taskInfo.State.String(),
		})
	}
	return
}
