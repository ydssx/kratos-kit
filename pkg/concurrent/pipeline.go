package concurrent

import (
	"context"
	"sync"
)

// Pipeline 表示一个数据处理管道
type Pipeline[T any] struct {
	stages []func(in <-chan T) <-chan T
	buffer int // 添加缓冲区大小配置
}

// NewPipeline 创建一个新的Pipeline
func NewPipeline[T any](opts ...PipelineOption[T]) *Pipeline[T] {
	p := &Pipeline[T]{
		stages: make([]func(in <-chan T) <-chan T, 0),
		buffer: 1, // 默认缓冲区大小
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// PipelineOption 定义管道选项类型
type PipelineOption[T any] func(*Pipeline[T])

// WithBuffer 设置管道缓冲区大小
func WithBuffer[T any](size int) PipelineOption[T] {
	return func(p *Pipeline[T]) {
		p.buffer = size
	}
}

// Run 执行Pipeline
func (p *Pipeline[T]) Run(ctx context.Context, source <-chan T) <-chan T {
	if len(p.stages) == 0 {
		return source
	}

	out := make(chan T, p.buffer)
	go func() {
		defer close(out)
		current := source
		for _, stage := range p.stages {
			select {
			case <-ctx.Done():
				return
			default:
				current = stage(current)
			}
		}
		for item := range current {
			select {
			case <-ctx.Done():
				return
			case out <- item:
			}
		}
	}()
	return out
}

// RunWithWorkers 使用指定数量的worker执行Pipeline
func (p *Pipeline[T]) RunWithWorkers(ctx context.Context, source <-chan T, numWorkers int) <-chan T {
	if numWorkers <= 0 {
		numWorkers = 1
	}

	out := make(chan T, p.buffer)
	var wg sync.WaitGroup
	
	// 创建工作池
	pool := make(chan struct{}, numWorkers)

	worker := func() {
		defer wg.Done()
		for item := range source {
			select {
			case <-ctx.Done():
				return
			case pool <- struct{}{}: // 获取工作槽
				var processedItem T
				processedItem = item
				var ok bool
				for _, stage := range p.stages {
					if processedItem, ok = p.processStage(ctx, stage, processedItem); !ok {
						<-pool // 释放工作槽
						return
					}
				}
				select {
				case <-ctx.Done():
					<-pool // 释放工作槽
					return
				case out <- processedItem:
				}
				<-pool // 释放工作槽
			}
		}
	}

	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// processStage 处理单个管道阶段
func (p *Pipeline[T]) processStage(ctx context.Context, stage func(in <-chan T) <-chan T, input T) (T, bool) {
	ch := make(chan T, 1)
	ch <- input
	close(ch)
	
	result := stage(ch)
	select {
	case <-ctx.Done():
		var zero T
		return zero, false
	case output, ok := <-result:
		if !ok {
			var zero T
			return zero, false
		}
		return output, true
	}
}
