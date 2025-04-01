package tasks

import (
	"context"
	"log/slog"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

// NewPool creates a new pool with the given name, cap and config.
func NewPool(capacity int32) *pool {
	p := &pool{}
	p.Capacity.Store(capacity)
	return p
}

type pool struct {
	sync.Mutex

	// linked list of tasks (tail push head pop) ; 默认处理的是并发处理任务不多的情况
	head *task
	tail *task

	// Capacity of the pool, the maximum number of goroutines that are actually working
	Capacity atomic.Int32
	// Record the number of running workers
	worker atomic.Int32

	length atomic.Int32 // 当前任务数
}

// task list
type task struct {
	// context is trade 执行链的上下文, 主要用于 panic handler 追查链
	ctx  context.Context
	fn   func(ctx context.Context)
	next *task
}

func (p *pool) Push(tail *task) {
	p.Lock()
	// tail push
	if p.tail != nil {
		// normal case, tail push
		p.tail.next = tail
	} else {
		// first push, head = tail
		p.head = tail
	}
	p.tail = tail
	p.Unlock()

	p.length.Add(1)
}

func (p *pool) Pop() (head *task) {
	p.Lock()
	head = p.head
	if head == nil {
		p.Unlock()
		return
	}

	// normal case head != nil
	p.head = head.next
	if p.head == nil {
		p.tail = nil // 队列为空，tail 随同 head 指回 nil
	}

	p.Unlock()
	p.length.Add(-1)
	return
}

// Go pool add task before run
// c 可能需要 chain.CopyTrace(c, keys) 基于业务独立考虑 // context 脱敏 & 延长生命周期
func (p *pool) Go(c context.Context, fn func(context.Context)) {
	tail := &task{
		ctx: c, // 需要自行进行 context 脱敏 & 延长生命周期
		fn:  fn,
	}

	// tail push
	p.Push(tail)

	// The current number of workers is less than the upper limit p.cap.
	worker := p.worker.Load()
	if worker < p.Capacity.Load() && p.worker.CompareAndSwap(worker, worker+1) {
		go p.running()
	}
}

func (p *pool) running() {
	for {
		// pop head after run task
		head := p.Pop()
		if head == nil {
			break
		}

		func() {
			defer func() {
				if cover := recover(); cover != nil {
					slog.ErrorContext(head.ctx,
						"tasks pool worker run panic error",
						slog.Any("error", cover),
						slog.String("stack", string(debug.Stack())),
					)
				}
			}()

			head.fn(head.ctx)
		}()
	}

	p.worker.Add(-1)
}

// p.Capacity p.Worker() p.Len() 属于运行时内部监控

func (p *pool) Worker() int32 {
	return p.worker.Load()
}

func (p *pool) Len() int32 {
	return p.length.Load()
}
