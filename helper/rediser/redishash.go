package rediser

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

// HSet 设置哈希表字段值
func (r *Client) HSet(ctx context.Context, key, field string, value any) error {
	err := r.UniversalClient.HSet(ctx, key, field, value).Err()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HSet error", slog.String("key", key), slog.String("field", field), slog.String("error", err.Error()))
	}
	return err
}

// HGet 获取哈希表字段值
func (r *Client) HGet(ctx context.Context, key, field string) (val string, ok bool, err error) {
	val, err = r.UniversalClient.HGet(ctx, key, field).Result()
	if err != nil {
		if err == redis.Nil {
			slog.InfoContext(ctx, "Redis HGet field not found", slog.String("key", key), slog.String("field", field), slog.String("error", err.Error()))
			return "", false, nil
		}
		slog.ErrorContext(ctx, "Redis HGet error", slog.String("key", key), slog.String("field", field), slog.String("error", err.Error()))
		return
	}
	return val, true, nil
}
