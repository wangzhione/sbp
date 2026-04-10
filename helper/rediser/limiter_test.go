package rediser

import (
	"testing"
	"time"
)

func TestLimiter_Allow(t *testing.T) {
	r := requireRedis(t)

	key := "mykey"
	ttl := 3 * time.Second
	rate := NewLimiter(r, key, ttl)

	t.Log(rate.Allow(ctx)) // true
	t.Log(rate.Allow(ctx)) // false
	time.Sleep(time.Second)
	t.Log(rate.Allow(ctx)) // false
	time.Sleep(time.Second)
	t.Log(rate.Allow(ctx)) // false
	time.Sleep(time.Second)
	t.Log(rate.Allow(ctx)) // true
}
