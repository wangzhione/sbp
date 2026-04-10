package rediser

import (
	"context"
	"os"
	"testing"

	"github.com/wangzhione/sbp/chain"
)

var ctx = chain.Context()

// testRedisCommand 优先读取环境变量，方便在 CI 或本机显式指定 Redis 连接参数。
// 未配置时回退到 redis-cli 默认参数，对应 localhost:6379。
func testRedisCommand() string {
	if command := os.Getenv("SBP_TEST_REDIS"); command != "" {
		return command
	}
	return "redis-cli"
}

// requireRedis 用于 Redis 集成测试前置检查。
// 如果当前环境没有可用 Redis，则跳过测试，避免默认 go test ./... 因外部依赖失败。
func requireRedis(t *testing.T) *Client {
	t.Helper()

	testctx, cancel := context.WithCancel(ctx)
	defer cancel()

	r, err := NewDefaultRedis(testctx, testRedisCommand())
	if err != nil {
		t.Skipf("skip redis integration test: %v", err)
	}

	t.Cleanup(func() {
		_ = r.Close(testctx)
	})

	return r
}

func TestClient_Eval(t *testing.T) {
	r := requireRedis(t)

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
