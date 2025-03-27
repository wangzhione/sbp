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
		c:         chain.Context(),
		oo:        make(chan T, buffersize),
		sem:       make(chan struct{}, maxgoworker),
		WokerLife: 10 * time.Second, // 1000 * 10ms
	}
}

type Pool[T Tasker] struct {
	// oo T 的任务池
	// 1. T 中如果有 context.Context 请用 chain.CopyTrace
	// 2. p.oo <- T 用于发送任务
	oo  chan T
	sem chan struct{} // make(chan struct{}, max go worker)
	c   context.Context

	WokerLife time.Duration // Pool[T].worker() 存活周期, 默认 10s
}

func (p *Pool[T]) Push(task T) {
	// 无论缓冲区是否满，只要可能，就拉 worker
	select {
	case p.sem <- struct{}{}:
		go p.worker(task)
		return
	default:
		// 达到最大并发，靠现有 worker 消费
	}

	// 提交任务，存在阻塞的可能, 那就等一等
	p.oo <- task
}

func (p *Pool[T]) worker(task T) {
	defer func() {
		if cover := recover(); cover != nil {
			slog.ErrorContext(p.c, "Pool worker panic error",
				slog.Any("error", cover),
				slog.String("stack", string(debug.Stack())),
			)
		}

		<-p.sem
	}()

	// 执行首次任务, 防止首次空转
	task.Do()

	timer := time.NewTimer(p.WokerLife)
	defer timer.Stop()
	for {
		select {
		case task := <-p.oo:
			ResetTimer(timer, p.WokerLife)
			task.Do()
		case <-timer.C:
			return
		}
	}
}

func ResetTimer(timer *time.Timer, life time.Duration) {
	// For a Timer created with NewTimer,
	// Reset should be invoked only on stopped or expired timers with drained channels.
	if !timer.Stop() {
		select {
		case <-timer.C: // try to drain the channel
		default:
		}
	}
	timer.Reset(life)
}
