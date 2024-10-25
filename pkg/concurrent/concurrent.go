package concurrent

import (
	"context"
	"fmt"
	"runtime/debug" // 添加 debug 包导入
	"time"

	"golang.org/x/sync/errgroup"
)

type Group struct {
	limit     int
	eg        *errgroup.Group
	ctx       context.Context
	fastFail  bool
	timeout   time.Duration
	recover   bool
}

type Opt func(*Group)

// WithSemaphore sets the semaphore limit for concurrent execution.
// It takes in an integer value for the semaphore and returns
// an Opt function that can be passed to NewGroup. The Opt function
// sets the limit field on the Group to control the concurrency.
func WithSemaphore(semaphore int) Opt {
	return func(g *Group) { g.limit = semaphore }
}

// WithFastFail returns an Opt function that sets the fastFail field on
// the Group. If set to true, Run will not wait for all goroutines to complete
// before returning the first error.
func WithFastFail(fastFail bool) Opt {
	return func(g *Group) { g.fastFail = fastFail }
}

func WithTimeout(timeout time.Duration) Opt {
	return func(g *Group) { g.timeout = timeout }
}

func WithRecover(recover bool) Opt {
	return func(g *Group) { g.recover = recover }
}

func NewGroup(ctx context.Context, opts ...Opt) *Group {
	if ctx == nil {
		ctx = context.Background()
	}

	g := &Group{ctx: ctx}
	for _, opt := range opts {
		opt(g)
	}

	if g.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, g.timeout)
		go func() {
			<-ctx.Done()
			cancel()
		}()
	}

	eg, ctx := errgroup.WithContext(ctx)
	g.eg = eg
	g.ctx = ctx

	if g.limit > 0 {
		eg.SetLimit(g.limit)
	}

	return g
}

// Run executes a group of functions concurrently.
//
// The Run function takes in multiple functions as arguments and executes them concurrently using goroutines. It adds each function to the underlying errgroup Group using the Go method. If the fastFail flag is set to false, the function waits for all the goroutines to complete by calling the Wait method on the Group. If the fastFail flag is set to true, the function spins off a separate goroutine to wait for the completion of the Group. After all the goroutines have completed, the function returns the error, if any, encountered during the execution of the goroutines.
//
// Parameters:
//   - fs: A variadic parameter that accepts functions of type func() error. These functions are executed concurrently.
//
// Returns:
//   - error: An error that occurred during the execution of the goroutines, if any. If no error occurred, the function returns nil.
func (g *Group) Run(fs ...func() error) error {
	if len(fs) == 0 {
		return nil
	}

	for _, f := range fs {
		f := f // 捕获变量
		g.eg.Go(func() error {
			if g.recover {
				defer func() {
					if r := recover(); r != nil {
						err := fmt.Errorf("panic recovered: %v\nstack:\n%s", r, debug.Stack())
						g.eg.Go(func() error { return err })
					}
				}()
			}
			return f()
		})
	}

	if !g.fastFail {
		return g.eg.Wait()
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- g.eg.Wait()
	}()

	select {
	case <-g.ctx.Done():
		return g.ctx.Err()
	case err := <-errCh:
		return err
	}
}
