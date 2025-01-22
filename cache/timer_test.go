package cache

import (
	"testing"
	"time"
)

func Test_cachedTimer_Stop(t *testing.T) {
	stopTimer := NewCachedTimer()

	time.Sleep(time.Second)

	stopTimer.Stop()
}
