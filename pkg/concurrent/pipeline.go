package concurrent

import (
	"context"
	"sync"
)

// Pipeline 表示一个数据处理管道
type Pipeline[T any] struct {
	stages []func(in <-chan T) <-chan T
}

// NewPipeline 创建一个新的Pipeline
func NewPipeline[T any]() *Pipeline[T] {
	return &Pipeline[T]{
		stages: make([]func(in <-chan T) <-chan T, 0),
	}
}

// AddStage 向Pipeline添加一个处理阶段
func (p *Pipeline[T]) AddStage(stage func(in <-chan T) <-chan T) *Pipeline[T] {
	p.stages = append(p.stages, stage)
	return p
}

// Run 执行Pipeline
func (p *Pipeline[T]) Run(ctx context.Context, source <-chan T) <-chan T {
	out := source
	for _, stage := range p.stages {
		out = stage(out)
	}
	return out
}

// RunWithWorkers 使用指定数量的worker执行Pipeline
func (p *Pipeline[T]) RunWithWorkers(ctx context.Context, source <-chan T, numWorkers int) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		for item := range source {
			result := item
			for _, stage := range p.stages {
				select {
				case <-ctx.Done():
					return
				default:
					ch := make(chan T, 1)
					ch <- result
					close(ch)
					result = <-stage(ch)
				}
			}
			select {
			case <-ctx.Done():
				return
			case out <- result:
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
