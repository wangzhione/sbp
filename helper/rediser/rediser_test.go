package rediser

import (
	"testing"

	"github.com/wangzhione/sbp/chain"
)

var ctx = chain.Context()

func TestClient_Eval(t *testing.T) {
	r, err := NewDefaultRedis(ctx, command)
	if err != nil {
		t.Fatal("fatal new redis", err, command)
	}

	script := `
		local value = redis.call("GET", KEYS[1])
		if value then
			return value
		else
			redis.call("SET", KEYS[1], ARGV[1], "EX", 10)
			return ARGV[1]
		end
	`

	key := "mykey"
	defaultValue := "default_value"

	result, err := r.Eval(ctx, script, []string{key}, defaultValue)
	if err != nil {
		t.Fatal("Failed to execute Lua script:", err)
	}

	t.Log("Success Lua script result:", result)
}
