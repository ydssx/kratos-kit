package concurrent

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

// Future 代表一个异步操作的结果
type Future[T any] struct {
	result T
	err    error
	done   chan struct{}
}

// NewFuture 创建一个新的Future
func NewFuture[T any](f func() (T, error)) *Future[T] {
	future := &Future[T]{
		done: make(chan struct{}),
	}
	go func() {
		defer close(future.done)
		future.result, future.err = f()
	}()
	return future
}

// Await 等待异步操作完成并返回结果
func (f *Future[T]) Await() (T, error) {
	<-f.done
	return f.result, f.err
}

// IsDone 判断异步操作是否完成
func (f *Future[T]) IsDone() bool {
	select {
	case <-f.done:
		return true
	default:
		return false
	}
}

// WaitFutures 等待多个异步操作完成
func WaitFutures[T any](concurrency int, funcs ...func() (T, error)) []*Future[T] {
	sem := semaphore.NewWeighted(int64(concurrency))
	ctx := context.Background()
	futures := make([]*Future[T], len(funcs))

	for i, f := range funcs {
		_ = sem.Acquire(ctx, 1)
		futures[i] = NewFuture(func() (T, error) {
			defer sem.Release(1)
			return f()
		})
	}

	// 等待所有操作完成
	for _, future := range futures {
		if !future.IsDone() {
			<-future.done
		}
	}

	return futures
}

// RunFutures 并发执行多个函数，限制并发数
func RunFutures[T any](concurrency int, fns ...func() (T, error)) []*Future[T] {
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
			futures[i] = NewFuture(fn)
		}(i, fn)
	}

	wg.Wait()
	return futures
}
