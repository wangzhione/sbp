package rediser

import (
	"testing"
	"time"
)

func TestLimiter_Allow(t *testing.T) {
	command := "redis-cli"

	r, err := NewDefaultRedis(ctx, command)
	if err != nil {
		t.Fatal("fatal new redis", err, command)
	}

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
