package chango

import (
	"fmt"
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
