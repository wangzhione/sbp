package cache

import (
	"sync/atomic"
	"time"
)

// Timer holds representation of current time.
type Timer interface {
	// Give current time (in seconds)
	Now() int64
}

// Timer that must be stopped.
type StoppableTimer interface {
	Timer

	// Release resources of the timer, functionality may or may not be affected
	// It is not called automatically, so user must call it just once
	Stop()
}

// Default timer reads Unix time always when requested
type defaultTimer struct{}

func (timer defaultTimer) Now() int64 {
	return time.Now().Unix()
}

// Cached timer stores Unix time every second and returns the cached value
type cachedTimer struct {
	now    int64
	ticker *time.Ticker
	done   chan struct{}
}

// Create cached timer and start runtime timer that updates time every second
func NewCachedTimer() StoppableTimer {
	timer := &cachedTimer{
		now:    time.Now().Unix(),
		ticker: time.NewTicker(time.Second),
		done:   make(chan struct{}),
	}

	go timer.update()

	return timer
}

func (timer *cachedTimer) Now() int64 {
	return atomic.LoadInt64(&timer.now)
}

// Stop runtime timer and finish routine that updates time
func (timer *cachedTimer) Stop() {
	timer.ticker.Stop()
	close(timer.done)
}

// Periodically check and update  of time
func (timer *cachedTimer) update() {
	for {
		select {
		case <-timer.done:
			return
		case <-timer.ticker.C:
			atomic.StoreInt64(&timer.now, time.Now().Unix())
		}
	}
}
