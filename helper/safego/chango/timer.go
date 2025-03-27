package chango

import "time"

func StopTimer(timer *time.Timer) {
	// For a Timer created with NewTimer,
	// Reset should be invoked only on stopped or expired timers with drained channels.
	if !timer.Stop() {
		select {
		case <-timer.C: // try to drain the channel
		default:
		}
	}
}

func ResetTimer(timer *time.Timer, life time.Duration) {
	StopTimer(timer)
	timer.Reset(life)
}

type Timer struct {
	*time.Timer
	Life time.Duration // time.Timer.C ticker
}

func NewTimer(d time.Duration) *Timer {
	return &Timer{
		Timer: time.NewTimer(d),
		Life:  d,
	}
}

func (t *Timer) Stop() {
	StopTimer(t.Timer)
}

func (t *Timer) Reset() {
	ResetTimer(t.Timer, t.Life)
}
