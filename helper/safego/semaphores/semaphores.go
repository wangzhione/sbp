// Package semaphores provides a semaphore implementation for controlling access to resources.
package semaphores

import (
	"context"
	"errors"
)

type Semaphores struct {
	Semaphore chan struct{}
}

func NewSemaphores(size int) (p Semaphores) {
	if size > 0 {
		p.Semaphore = make(chan struct{}, size)
	} else {
		p.Semaphore = make(chan struct{})
	}

	return p
}

var ErrTooMany = errors.New("errors: semaphores too many value")

func (p *Semaphores) TryLock(ctx context.Context) error {
	select {
	case <-ctx.Done():
		// 等待过程中被取消/超时
		return ctx.Err()
	case p.Semaphore <- struct{}{}:
		// 获取成功
		return nil
	default:
		// 已满，直接返回业务错误
		return ErrTooMany
	}
}

func (p *Semaphores) Lock(ctx context.Context) error {
	select {
	case <-ctx.Done():
		// 等待过程中被取消/超时
		return ctx.Err()
	case p.Semaphore <- struct{}{}:
		// 获取成功
		return nil
	}
}

func (p *Semaphores) UnLock() {
	select {
	case <-p.Semaphore:
		// 释放成功
	default:
		// 到这里比较危险, 依赖使用方自己保证最终结果
	}
}
