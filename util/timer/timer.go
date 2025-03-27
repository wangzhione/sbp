package timer

import "time"

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
	// For a Timer created with NewTimer,
	// Reset should be invoked only on stopped or expired timers with drained channels.
	if !t.Timer.Stop() {
		select {
		case <-t.C: // try to drain the channel
		default:
		}
	}
}

func (t *Timer) Reset() {
	t.Stop()
	t.Timer.Stop()
}
