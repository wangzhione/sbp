package groupgo

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"sync"
)

// A Group is a collection of goroutines working on subtasks that are part of
// the same overall task.
//
// A zero Group is valid, has no limit on the number of active goroutines,
// and does not cancel on error.
type Group struct {
	c   context.Context
	wg  sync.WaitGroup
	sem chan struct{}
	err error
}

func (g *Group) done() {
	<-g.sem
	g.wg.Done()
}

// NewGroup groupgo 推荐 n > 0 的业务, 最多启动 n 个 goroutine 业务去处理
func NewGroup(ctx context.Context, n int) *Group {
	g := &Group{
		c:   ctx,
		sem: make(chan struct{}, n),
	}
	return g
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *Group) Wait() error {
	g.wg.Wait()
	return g.err
}

// Go calls the given function in a new goroutine.
// It blocks until the new goroutine can be added without the number of
// active goroutines in the group exceeding the configured limit.
//
// The first call to return a non-nil error cancels the group's context, if the
// group was created by calling WithContext. The error will be returned by Wait.
func (g *Group) Go(f func(ctx context.Context) error) {
	g.sem <- struct{}{}

	g.wg.Add(1)
	go func() {
		defer func() {
			if cover := recover(); cover != nil {
				g.err = fmt.Errorf("panic recovered: %#v", cover)

				// 遇到启动不起来, 异常退出, 打印堆栈方便排除问题
				slog.ErrorContext(g.c, "Group Go panic error",
					slog.Any("error", cover),
					slog.String("stack", string(debug.Stack())), // 记录详细的堆栈信息
				)
			}

			g.done()
		}()

		if err := f(g.c); err != nil {
			g.err = err
			slog.ErrorContext(g.c, "Group f call error", "error", err)
		}
	}()
}
