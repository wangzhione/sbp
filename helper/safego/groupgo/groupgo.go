// Package groupgo provides a goroutine group with concurrency limits and error handling.
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
// 具有 limit 限制 goroutine group call 方式
type Group struct {
	c   context.Context
	wg  sync.WaitGroup
	sem chan struct{}

	errOnce sync.Once
	err     error // 默认只纪录第一个错误
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
	// 兼容, nil Group wait 模式
	if g == nil {
		return nil
	}

	g.wg.Wait()
	return g.err
}

// Go calls the given function in a new goroutine.
// It blocks until the new goroutine can be added without the number of
// active goroutines in the group exceeding the configured limit.
// 会纪录首次出现 error 信息
func (g *Group) Go(f func(ctx context.Context) error) {
	g.sem <- struct{}{}

	g.wg.Add(1)
	go func() {
		defer func() {
			if cover := recover(); cover != nil {
				g.errOnce.Do(func() {
					g.err = fmt.Errorf("panic: groupgo.Group.Go %#v", cover)
				})

				// 遇到启动不起来, 异常退出, 打印堆栈方便排除问题
				slog.ErrorContext(g.c, "Group Go panic error",
					slog.Any("error", cover),
					slog.String("stack", string(debug.Stack())), // 记录详细的堆栈信息
				)
			}

			g.done()
		}()

		if err := f(g.c); err != nil {
			slog.ErrorContext(g.c, "Group f call error", "error", err)

			g.errOnce.Do(func() {
				g.err = err
			})
		}
	}()
}
