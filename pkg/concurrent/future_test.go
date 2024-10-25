package concurrent

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestFuture(t *testing.T) {
	t.Run("Basic Future", func(t *testing.T) {
		ctx := context.Background()
		expected := 42
		f := NewFuture(ctx, func(ctx context.Context) (int, error) {
			return expected, nil
		})

		result, err := f.Await()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result != expected {
			t.Errorf("Expected %d, got %d", expected, result)
		}
	})

	t.Run("Future with Error", func(t *testing.T) {
		ctx := context.Background()
		expectedErr := errors.New("test error")
		f := NewFuture(ctx, func(ctx context.Context) (int, error) {
			return 0, expectedErr
		})

		_, err := f.Await()
		if err != expectedErr {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("Future with Timeout", func(t *testing.T) {
		ctx := context.Background()
		f := NewFuture(ctx, func(ctx context.Context) (int, error) {
			time.Sleep(100 * time.Millisecond)
			return 42, nil
		})

		_, err := f.AwaitWithTimeout(50 * time.Millisecond)
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("Expected deadline exceeded error, got %v", err)
		}
	})

	t.Run("Future with Cancel", func(t *testing.T) {
		ctx := context.Background()
		f := NewFuture(ctx, func(ctx context.Context) (int, error) {
			select {
			case <-ctx.Done():
				return 0, ctx.Err()
			case <-time.After(100 * time.Millisecond):
				return 42, nil
			}
		})

		go func() {
			time.Sleep(50 * time.Millisecond)
			f.Cancel()
		}()

		_, err := f.Await()
		if !errors.Is(err, context.Canceled) {
			t.Errorf("Expected context canceled error, got %v", err)
		}
	})

	t.Run("WaitFutures", func(t *testing.T) {
		ctx := context.Background()
		funcs := []func(context.Context) (int, error){
			func(ctx context.Context) (int, error) { return 1, nil },
			func(ctx context.Context) (int, error) { return 2, nil },
			func(ctx context.Context) (int, error) { return 3, nil },
		}

		futures := WaitFutures(ctx, 2, funcs...)

		sum := 0
		for _, f := range futures {
			result, err := f.Await()
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			sum += result
		}

		if sum != 6 {
			t.Errorf("Expected sum 6, got %d", sum)
		}
	})

	t.Run("RunFutures", func(t *testing.T) {
		ctx := context.Background()
		funcs := []func() (int, error){
			func() (int, error) {
				time.Sleep(50 * time.Millisecond)
				return 1, nil
			},
			func() (int, error) {
				time.Sleep(50 * time.Millisecond)
				return 2, nil
			},
			func() (int, error) {
				time.Sleep(50 * time.Millisecond)
				return 3, nil
			},
		}

		start := time.Now()
		futures := RunFutures(ctx, 2, funcs...)
		
		sum := 0
		for _, f := range futures {
			result, err := f.Await()
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			sum += result
		}

		duration := time.Since(start)
		if duration < 100*time.Millisecond {
			t.Error("Concurrent execution did not work as expected")
		}

		if sum != 6 {
			t.Errorf("Expected sum 6, got %d", sum)
		}
	})

	t.Run("Future Concurrency Safety", func(t *testing.T) {
		ctx := context.Background()
		f := NewFuture(ctx, func(ctx context.Context) (int, error) {
			time.Sleep(50 * time.Millisecond)
			return 42, nil
		})

		// 并发调用 Await
		done := make(chan struct{})
		go func() {
			for i := 0; i < 10; i++ {
				go func() {
					_, _ = f.Await()
				}()
			}
			close(done)
		}()

		select {
		case <-done:
		case <-time.After(200 * time.Millisecond):
			t.Error("Concurrent Await calls did not complete in time")
		}
	})
}
