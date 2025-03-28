package timer

import "time"

// 要求 go version >= 1.23
// time.Timer 1.23 之后语义发生变化, 修复 1.23 之前啰嗦 stop select case 写法

// Don’t pool what you don’t need to. GC is cheap. Your code’s clarity is expensive.
// 官方 net/http 标准库也没使用 timer pool（他们是高性能场景，虽然不极端）

type Timer struct {
	*time.Timer
	Life time.Duration // time.Timer.C ticker
}

func NewTimer(d time.Duration) *Timer {
	return &Timer{Timer: time.NewTimer(d), Life: d}
}

func (t *Timer) Reset() {
	t.Timer.Reset(t.Life)
}
