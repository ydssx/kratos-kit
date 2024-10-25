package concurrent

import (
	"context"
	"testing"
	"time"
)

func TestPipeline(t *testing.T) {
	// 测试 NewPipeline
	t.Run("NewPipeline", func(t *testing.T) {
		p := NewPipeline[int]()
		if p == nil {
			t.Error("NewPipeline should return a non-nil Pipeline")
		}
		if len(p.stages) != 0 {
			t.Error("NewPipeline should create an empty stages slice")
		}
	})

	// 测试 AddStage
	t.Run("AddStage", func(t *testing.T) {
		p := NewPipeline[int]()
		stage := func(in <-chan int) <-chan int {
			out := make(chan int)
			go func() {
				defer close(out)
				for n := range in {
					out <- n * 2
				}
			}()
			return out
		}
		p.AddStage(stage)
		if len(p.stages) != 1 {
			t.Error("AddStage should append the stage to stages slice")
		}
	})

	// 测试 Run
	t.Run("Run", func(t *testing.T) {
		p := NewPipeline[int]()
		p.AddStage(func(in <-chan int) <-chan int {
			out := make(chan int)
			go func() {
				defer close(out)
				for n := range in {
					out <- n * 2
				}
			}()
			return out
		}).AddStage(func(in <-chan int) <-chan int {
			out := make(chan int)
			go func() {
				defer close(out)
				for n := range in {
					out <- n + 1
				}
			}()
			return out
		})

		source := make(chan int)
		go func() {
			defer close(source)
			for i := 1; i <= 3; i++ {
				source <- i
			}
		}()

		ctx := context.Background()
		result := p.Run(ctx, source)

		expected := []int{3, 5, 7}
		for i, exp := range expected {
			if res := <-result; res != exp {
				t.Errorf("Expected %d, got %d for input %d", exp, res, i+1)
			}
		}
	})

	// 测试 RunWithWorkers
	t.Run("RunWithWorkers", func(t *testing.T) {
		p := NewPipeline[int]()
		p.AddStage(func(in <-chan int) <-chan int {
			out := make(chan int)
			go func() {
				defer close(out)
				for n := range in {
					out <- n * 2
				}
			}()
			return out
		}).AddStage(func(in <-chan int) <-chan int {
			out := make(chan int)
			go func() {
				defer close(out)
				for n := range in {
					out <- n + 1
				}
			}()
			return out
		})

		source := make(chan int)
		go func() {
			defer close(source)
			for i := 1; i <= 10; i++ {
				source <- i
			}
		}()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result := p.RunWithWorkers(ctx, source, 3)

		expected := make(map[int]bool)
		for i := 1; i <= 10; i++ {
			expected[i*2+1] = true
		}

		for res := range result {
			if !expected[res] {
				t.Errorf("Unexpected result: %d", res)
			}
			delete(expected, res)
		}

		if len(expected) != 0 {
			t.Errorf("Missing expected results: %v", expected)
		}
	})
}