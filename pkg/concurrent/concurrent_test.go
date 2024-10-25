package concurrent

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)

	os.Exit(m.Run())
}

func TestGroup(t *testing.T) {
	t.Run("Basic Execution", func(t *testing.T) {
		g := NewGroup(context.Background())
		result := 0

		err := g.Run(
			func() error {
				result += 1
				return nil
			},
			func() error {
				result += 2
				return nil
			},
		)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result != 3 {
			t.Errorf("Expected result 3, got %d", result)
		}
	})

	t.Run("Error Handling", func(t *testing.T) {
		g := NewGroup(context.Background())
		expectedErr := errors.New("test error")

		err := g.Run(
			func() error {
				return expectedErr
			},
			func() error {
				return nil
			},
		)

		if err != expectedErr {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("Fast Fail", func(t *testing.T) {
		g := NewGroup(context.Background(), WithFastFail(true))
		expectedErr := errors.New("fast fail error")

		err := g.Run(
			func() error {
				time.Sleep(100 * time.Millisecond)
				return nil
			},
			func() error {
				return expectedErr
			},
		)

		if err != expectedErr {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("Timeout", func(t *testing.T) {
		g := NewGroup(context.Background(), WithTimeout(50*time.Millisecond))

		err := g.Run(
			func() error {
				time.Sleep(100 * time.Millisecond)
				return nil
			},
		)

		if !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("Expected deadline exceeded error, got %v", err)
		}
	})

	t.Run("Panic Recovery", func(t *testing.T) {
		g := NewGroup(context.Background(), WithRecover(true))

		err := g.Run(
			func() error {
				panic("test panic")
			},
		)

		if err == nil {
			t.Error("Expected error from panic, got nil")
		}
		if !errors.Is(err, context.Canceled) {
			t.Errorf("Expected context canceled error, got %v", err)
		}
	})

	t.Run("Semaphore Limit", func(t *testing.T) {
		g := NewGroup(context.Background(), WithSemaphore(2))
		running := make(chan struct{}, 3)
		done := make(chan struct{})

		err := g.Run(
			func() error {
				running <- struct{}{}
				<-done
				return nil
			},
			func() error {
				running <- struct{}{}
				<-done
				return nil
			},
			func() error {
				running <- struct{}{}
				<-done
				return nil
			},
		)

		// 等待一小段时间让goroutines启动
		time.Sleep(50 * time.Millisecond)

		// 检查同时运行的goroutine数量
		if len(running) > 2 {
			t.Error("More than 2 goroutines running simultaneously")
		}

		// 清理
		close(done)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}
