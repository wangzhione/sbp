package rediser

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	redis.UniversalClient
}

// Close 关闭数据库连接, 必须主动去执行, 否则无法被回收
func (r *Client) Close(ctx context.Context) (err error) {
	if r != nil && r.UniversalClient != nil {
		err = r.UniversalClient.Close()
		// 创建和关闭都是很重的操作需要格外小心
		slog.InfoContext(ctx, "r.UniversalClient.Close() info", "reason", err)
	}
	return
}

// Do 执行原生 Redis 命令
func (r *Client) Do(ctx context.Context, args ...any) (result any, err error) {
	result, err = r.UniversalClient.Do(ctx, args...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis Do error", slog.Any("args", args), slog.String("error", err.Error()))
		return
	}
	return
}

// LPush 向列表左侧插入值
func (r *Client) LPush(ctx context.Context, key string, values ...any) error {
	err := r.UniversalClient.LPush(ctx, key, values...).Err()
	if err != nil {
		slog.ErrorContext(ctx, "Redis LPush error", slog.String("key", key), slog.String("error", err.Error()))
	}
	return err
}

// RPop 从列表右侧弹出值
func (r *Client) RPop(ctx context.Context, key string) (value string, ok bool, err error) {
	value, err = r.UniversalClient.RPop(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			slog.InfoContext(ctx, "Redis RPop key empty", slog.String("key", key), slog.String("error", err.Error()))
			return "", false, nil
		}
		slog.ErrorContext(ctx, "Redis RPop error", slog.String("key", key), slog.String("error", err.Error()))
		return
	}
	return value, true, nil
}

// LRange 获取列表 key 中指定区间 [start, stop] 的元素
func (r *Client) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	values, err := r.UniversalClient.LRange(ctx, key, start, stop).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis LRange error",
			slog.String("key", key),
			slog.Int64("start", start),
			slog.Int64("stop", stop),
			slog.String("error", err.Error()))
		return nil, err
	}
	return values, nil
}

// Eval 执行 Lua 脚本
func (r *Client) Eval(ctx context.Context, script string, keys []string, args ...any) (result any, err error) {
	result, err = r.UniversalClient.Eval(ctx, script, keys, args...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis Eval error", slog.String("script", script), slog.String("error", err.Error()))
		return
	}
	return
}

// BLPop 从列表左侧阻塞式弹出
func (r *Client) BLPop(ctx context.Context, timeout time.Duration, keys ...string) (values []string, err error) {
	values, err = r.UniversalClient.BLPop(ctx, timeout, keys...).Result()
	if err != nil {
		if err == redis.Nil {
			slog.InfoContext(ctx, "Redis BLPop timeout or empty", slog.Any("keys", keys))
			return nil, nil
		}
		slog.ErrorContext(ctx, "Redis BLPop failed", slog.Any("keys", keys), slog.String("error", err.Error()))
	}
	return
}

// BRPop 从列表右侧阻塞式弹出
func (r *Client) BRPop(ctx context.Context, timeout time.Duration, keys ...string) (values []string, err error) {
	values, err = r.UniversalClient.BRPop(ctx, timeout, keys...).Result()
	if err != nil {
		if err == redis.Nil {
			slog.InfoContext(ctx, "Redis BRPop timeout or empty", slog.Any("keys", keys))
			return nil, nil
		}
		slog.ErrorContext(ctx, "Redis BRPop failed", slog.Any("keys", keys), slog.String("error", err.Error()))
	}
	return
}
