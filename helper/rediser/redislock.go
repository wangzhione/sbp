package rediser

import (
	"context"
	"log/slog"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/wangzhione/sbp/chain"
)

// 单机简单版本 redis lock
// 适合那些业务, 能忍受极极个别情况下 锁意外 发生两次情况, 例如 redis lock 控制 发 email, 存在发两次可能

type RedisLock struct {
	r     *Client // 支持 *redis.Client 或 *redis.ClusterClient
	key   string  // 推荐 lock:{name} 这样命名
	value string  // 用于标识 owner
}

// TryLock 尝试加锁
// return *RedisLock is not nil, 表示获取到了锁资源
func (r *Client) TryLock(ctx context.Context, key string, ttl time.Duration) (*RedisLock, error) {
	value := chain.UUID()

	// 加入 0 ~ 10ms 的随机抖动，缓解锁同时过期
	ttl += time.Duration(rand.Int63n(int64(10 * time.Millisecond)))

	ok, err := r.SetNX(ctx, key, value, ttl)
	if err != nil {
		slog.ErrorContext(ctx, "TryLock SetNX error", "key", key, "ttl", ttl, "error", err, "value", value)
		return nil, err
	}
	if !ok {
		slog.InfoContext(ctx, "TryLock already held", "key", key, "ttl", ttl, "value", value)
		return nil, nil
	}

	return &RedisLock{r: r, key: key, value: value}, nil
}

// Unlock 安全解锁，使用 Lua 保证原子性（只能删除自己加的锁）
func (l *RedisLock) Unlock(ctx context.Context) (success bool, err error) {
	luaScript := redis.NewScript(`
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
	`)
	res, err := luaScript.Run(ctx, l.r.UniversalClient, []string{l.key}, l.value).Result()
	if err != nil {
		slog.ErrorContext(ctx, "luaScript.Run error", "key", l.key, "value", l.value, "error", err)
		return
	}

	status, ok := res.(int64)
	if !ok {
		slog.ErrorContext(ctx, "Unlock result not int64 panic error", "key", l.key, "value", l.value, "result", res)
		return
	}
	return status == 1, nil
}

// TimeoutLock timeout 不是个准确数字, 是个大致靠近的数字
func (r *Client) TimeoutLock(ctx context.Context, key string, ttl, timeout time.Duration) (*RedisLock, error) {
	const interval = 130 * time.Millisecond // 重试间隔，默认 130 ms

	var timer *time.Timer

	deadline := time.Now().Add(timeout)
	for {
		lock, err := r.TryLock(ctx, key, ttl)
		if err != nil {
			return nil, err // redis 或网络错误，立即返回
		}
		if lock != nil {
			return lock, nil // 获取锁成功
		}

		// 检查超时
		if time.Now().After(deadline) {
			return nil, nil // 超过最大等待时间
		}

		if timer == nil { // lazy 惰性创建 Timer
			timer = time.NewTimer(interval)
			defer timer.Stop()
		} else {
			timer.Reset(interval)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timer.C:
			// retry
		}
	}
}
