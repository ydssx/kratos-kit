package concurrent

import (
	"context"
	"sync"
	"time"
)

// Future 代表一个异步操作的结果
type Future[T any] struct {
	result T
	err    error
	done   chan struct{}
	once   sync.Once
	ctx    context.Context
	cancel context.CancelFunc
}

// NewFuture 创建一个新的Future
func NewFuture[T any](ctx context.Context, f func(context.Context) (T, error)) *Future[T] {
	ctx, cancel := context.WithCancel(ctx)
	future := &Future[T]{
		done:   make(chan struct{}),
		ctx:    ctx,
		cancel: cancel,
	}

	go func() {
		defer future.once.Do(func() {
			close(future.done)
			cancel()
		})
		future.result, future.err = f(ctx)
	}()
	return future
}

// Await 等待异步操作完成并返回结果
func (f *Future[T]) Await() (T, error) {
	select {
	case <-f.ctx.Done():
		if f.err == nil {
			f.err = f.ctx.Err()
		}
	case <-f.done:
	}
	return f.result, f.err
}

// AwaitWithTimeout 带超时的等待
func (f *Future[T]) AwaitWithTimeout(timeout time.Duration) (T, error) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-timer.C:
		var zero T
		return zero, context.DeadlineExceeded
	case <-f.done:
		return f.result, f.err
	}
}

// Cancel 取消异步操作
func (f *Future[T]) Cancel() {
	f.cancel()
}

// WaitFutures 等待多个异步操作完成
func WaitFutures[T any](ctx context.Context, concurrency int, funcs ...func(context.Context) (T, error)) []*Future[T] {
	if concurrency <= 0 {
		concurrency = len(funcs)
	}

	futures := make([]*Future[T], len(funcs))
	sem := make(chan struct{}, concurrency)

	var wg sync.WaitGroup
	wg.Add(len(funcs))

	for i, f := range funcs {
		go func(i int, f func(context.Context) (T, error)) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			futures[i] = NewFuture(ctx, f)
		}(i, f)
	}

	wg.Wait()
	return futures
}

// RunFutures 并发执行多个函数，限制并发数
func RunFutures[T any](ctx context.Context, concurrency int, fns ...func() (T, error)) []*Future[T] {
	futures := make([]*Future[T], len(fns))
	var wg sync.WaitGroup
	var semaphore chan struct{}

	if concurrency > 0 {
		semaphore = make(chan struct{}, concurrency)
	}

	for i, fn := range fns {
		wg.Add(1)
		go func(i int, fn func() (T, error)) {
			defer wg.Done()
			if concurrency > 0 {
				semaphore <- struct{}{}
				defer func() { <-semaphore }()
			}
			futures[i] = NewFuture(ctx, func(ctx context.Context) (T, error) {
				return fn()
			})
		}(i, fn)
	}

	wg.Wait()
	return futures
}
