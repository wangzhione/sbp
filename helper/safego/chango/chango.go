// Package chango provides a generic goroutine pool for managing concurrent tasks.
// It supports worker lifecycle management, panic recovery, and context cancellation.
package chango

import (
	"context"
	"log/slog"
	"runtime/debug"
	"time"

	"github.com/wangzhione/sbp/chain"
)

type Tasker interface {
	Do()
}

func NewPool[T Tasker](maxgoworker int, buffersize int) *Pool[T] {
	return &Pool[T]{
		C:         chain.Context(),
		oo:        make(chan T, buffersize),
		sem:       make(chan struct{}, maxgoworker),
		WokerLife: 10 * time.Second, // 默认 10s = 1000 * 10ms
	}
}

type Pool[T Tasker] struct {
	C         context.Context // 可以自行初始化时候设置 context
	WokerLife time.Duration   // Pool[T].worker() 存活周期

	// oo T 的任务池
	// 1. T 中如果有 context.Context 请用 chain.CopyTrace
	// 2. p.oo <- T 用于发送任务
	oo  chan T
	sem chan struct{} // make(chan struct{}, max go worker)
}

func (p *Pool[T]) Push(task T) {
	// 无论缓冲区是否满，只要可能，就拉 worker
	select {
	case <-p.C.Done():
		return
	case p.sem <- struct{}{}:
		go p.worker(task)
		return
	default:
		// 达到最大并发，靠现有 worker 消费
	}

	select {
	case <-p.C.Done():
		return
	case p.oo <- task: // 提交任务，存在阻塞的可能, 那就等一等
		return
	}
}

func (p *Pool[T]) worker(one T) {
	defer func() {
		if cover := recover(); cover != nil {
			slog.ErrorContext(p.C, "Pool worker panic error",
				slog.Any("error", cover),
				slog.String("stack", string(debug.Stack())),
			)
		}

		<-p.sem
	}()

	// 执行首次任务, 防止首次空转
	one.Do()

	// Go 1.23+ safe: Stop 无需 drain (<-timer.C 手工清空)；defer 保证释放
	r := time.NewTimer(p.WokerLife)
	defer r.Stop()

	for {
		select {
		case two := <-p.oo:
			two.Do()
			// 重新来过, 为下次 循环准备; 会 clear old <-r.C 的值
			// https://golang.ac.cn/wiki/Go123Timer
			r.Reset(p.WokerLife)
		case <-r.C:
			// 预防 timer 和 oo 都触发, 导致 worker 消费被遗漏
			select {
			case two := <-p.oo:
				two.Do()
				r.Reset(p.WokerLife)
				// 继续循环，不退出
				continue
			default:
				// 确实没任务，才真正退出
				return
			}

		case <-p.C.Done():
			// Pool 关闭/上层取消
			return
		}
	}
}
