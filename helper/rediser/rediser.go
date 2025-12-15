package rediser

import (
	"context"
	"log/slog"

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

// Eval 执行 Lua 脚本
func (r *Client) Eval(ctx context.Context, script string, keys []string, args ...any) (result any, err error) {
	result, err = r.UniversalClient.Eval(ctx, script, keys, args...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis Eval error", slog.String("script", script), slog.String("error", err.Error()))
		return
	}
	return
}
