package chango

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// 定义一个任务结构，实现 Tasker 接口
type MyTask struct {
	id int
}

func (t *MyTask) Do() {
	fmt.Printf("[Task %d] Start at %v\n", t.id, time.Now().Format("15:04:05.000"))
	time.Sleep(500 * time.Millisecond) // 模拟耗时任务
	fmt.Printf("[Task %d] Done at %v\n", t.id, time.Now().Format("15:04:05.000"))
}

func TestNewPool(t *testing.T) {
	// 创建一个池子，最多 4 个并发 worker，缓冲区 10
	pool := NewPool[*MyTask](4, 10)

	// 投入 8 个任务
	for i := 1; i <= 8; i++ {
		pool.Push(&MyTask{id: i})
	}

	// 等待所有任务完成（这里只是简单 sleep 模拟，正式应使用 sync.WaitGroup）
	time.Sleep(5 * time.Second)
}

// mockTask 是用于性能测试的假任务，实现 Tasker 接口
type mockTask struct {
	wg *sync.WaitGroup
}

func (m mockTask) Do() {
	time.Sleep(1 * time.Millisecond) // 模拟执行耗时
	m.wg.Done()
}

const (
	maxGoWorker = 100
	bufferSize  = 1000
)

func BenchmarkPool(b *testing.B) {
	pool := NewPool[mockTask](maxGoWorker, bufferSize)

	// 重置定时器，不统计初始化开销

	for range 1000 {
		var wg sync.WaitGroup
		wg.Add(maxGoWorker) // 每次压测发送 maxGoWorker 个任务

		for range maxGoWorker {
			pool.Push(mockTask{wg: &wg})
		}

		wg.Wait()
	}
}

/*
goos: windows
goarch: amd64
pkg: github.com/wangzhione/sbp/helper/safego/chango
cpu: AMD Ryzen 9 7945HX3D with Radeon Graphics
BenchmarkPool
BenchmarkPool-32
       1	1524289700 ns/op	  273680 B/op	    2814 allocs/op

BenchmarkRawGoroutine
BenchmarkRawGoroutine-32
       1	1536055500 ns/op	11422928 B/op	  201430 allocs/op
*/

func BenchmarkRawGoroutine(b *testing.B) {
	for range 1000 {
		var wg sync.WaitGroup
		wg.Add(maxGoWorker)

		for range maxGoWorker {
			go func() {
				time.Sleep(1 * time.Millisecond)
				wg.Done()
			}()
		}

		wg.Wait()
	}
}
