package concurrent

import (
	"context"
	"testing"
	"time"
)

func TestPipeline(t *testing.T) {
	t.Run("NewPipeline", func(t *testing.T) {
		p := NewPipeline[int]()
		if p == nil {
			t.Error("NewPipeline should return a non-nil Pipeline")
		}
		if len(p.stages) != 0 {
			t.Error("NewPipeline should create an empty stages slice")
		}
	})

	t.Run("WithBuffer", func(t *testing.T) {
		bufferSize := 10
		p := NewPipeline[int](WithBuffer[int](bufferSize))
		if p.buffer != bufferSize {
			t.Errorf("Expected buffer size %d, got %d", bufferSize, p.buffer)
		}
	})

	t.Run("Pipeline Processing", func(t *testing.T) {
		ctx := context.Background()
		p := NewPipeline[int]()

		// 添加两个处理阶段
		p.stages = append(p.stages, func(in <-chan int) <-chan int {
			out := make(chan int)
			go func() {
				defer close(out)
				for n := range in {
					out <- n * 2
				}
			}()
			return out
		})

		p.stages = append(p.stages, func(in <-chan int) <-chan int {
			out := make(chan int)
			go func() {
				defer close(out)
				for n := range in {
					out <- n + 1
				}
			}()
			return out
		})

		// 创建输入通道
		input := make(chan int)
		go func() {
			defer close(input)
			for i := 1; i <= 3; i++ {
				input <- i
			}
		}()

		// 运行管道
		output := p.Run(ctx, input)

		// 验证结果
		expected := []int{3, 5, 7} // (1*2)+1, (2*2)+1, (3*2)+1
		i := 0
		for result := range output {
			if result != expected[i] {
				t.Errorf("Expected %d, got %d", expected[i], result)
			}
			i++
		}
	})

	t.Run("Pipeline With Workers", func(t *testing.T) {
		ctx := context.Background()
		p := NewPipeline[int](WithBuffer[int](5))

		// 添加处理阶段
		p.stages = append(p.stages, func(in <-chan int) <-chan int {
			out := make(chan int)
			go func() {
				defer close(out)
				for n := range in {
					time.Sleep(10 * time.Millisecond) // 模拟耗时操作
					out <- n * 2
				}
			}()
			return out
		})

		// 创建输入通道
		input := make(chan int)
		go func() {
			defer close(input)
			for i := 1; i <= 10; i++ {
				input <- i
			}
		}()

		// 使用3个worker运行管道
		output := p.RunWithWorkers(ctx, input, 3)

		// 收集结果
		results := make([]int, 0, 10)
		for result := range output {
			results = append(results, result)
		}

		// 验证结果数量
		if len(results) != 10 {
			t.Errorf("Expected 10 results, got %d", len(results))
		}

		// 验证所有结果都是输入的2倍
		for i := range results {
			expected := (i + 1) * 2
			found := false
			for _, r := range results {
				if r == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected to find %d in results", expected)
			}
		}
	})

	t.Run("Context Cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		p := NewPipeline[int]()

		// 添加一个会阻塞的处理阶段
		p.stages = append(p.stages, func(in <-chan int) <-chan int {
			out := make(chan int)
			go func() {
				defer close(out)
				for n := range in {
					time.Sleep(100 * time.Millisecond)
					out <- n
				}
			}()
			return out
		})

		input := make(chan int)
		go func() {
			defer close(input)
			for i := 1; i <= 5; i++ {
				input <- i
			}
		}()

		output := p.Run(ctx, input)

		// 立即取消上下文
		cancel()

		// 验证是否及时退出
		timeout := time.After(200 * time.Millisecond)
		select {
		case <-timeout:
			t.Error("Pipeline did not cancel in time")
		case _, ok := <-output:
			if ok {
				t.Error("Expected pipeline to be cancelled")
			}
		}
	})
}
