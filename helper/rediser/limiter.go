// Package rediser provides Redis-based helpers including a rate limiter.
package rediser

import (
	"context"
	"log/slog"
	"time"

	"github.com/wangzhione/sbp/structs"
)

// Limiter 限制固定时间内最多请求 N 次
type Limiter struct {
	R   *Client
	Key string        // 限流 Key
	TTL time.Duration // 限制时长, 如 10*time.Minute
	N   int64         // 最大允许的请求次数
}

// NewLimiter 创建限流器
func NewLimiter(r *Client, key string, ttl time.Duration, limit ...int64) (rate *Limiter) {
	rate = &Limiter{
		R:   r,
		Key: key,
		TTL: ttl,
		N:   structs.Max(1, structs.Max(limit...)), // 默认 ttl 时间内, 只能有一次请求
	}
	return
}

// Allow 在固定时间内只能请求 N 次
func (rate *Limiter) Allow(ctx context.Context) bool {
	// 递增请求次数
	count, err := rate.R.Incr(ctx, rate.Key)
	if err != nil {
		slog.ErrorContext(ctx, "Redis Incr error", slog.String("key", rate.Key), slog.String("error", err.Error()))
		return false
	}

	// 第一次请求，设置 TTL
	if count == 1 {
		err = rate.R.Expire(ctx, rate.Key, rate.TTL)
		if err != nil {
			// 理论上不会出现, 如果出现需要人工干预
			slog.ErrorContext(
				ctx,
				"Redis Expire panic error",
				slog.String("rate.Key", rate.Key),
				slog.String("error", err.Error()),
				slog.Float64("rate.TTL", rate.TTL.Seconds()),
			)
		}
	}

	// 如果请求次数超过限制，则拒绝请求
	if count > rate.N {
		slog.InfoContext(ctx, "Rate limit exceeded", slog.String("key", rate.Key), slog.Int64("count", count))
		return false
	}

	return true
}
