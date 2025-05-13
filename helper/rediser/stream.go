package rediser

import (
	"context"
	"log/slog"
	"strings"

	"github.com/redis/go-redis/v9"
)

// 删除 Stream（即整个消息队列）
// err := r.RDB().Del(ctx, stream).Err()
// 删除某个 Group
// err := r.RDB().Do(ctx, "XGROUP", "DESTROY", stream, group).Err()
// 一般现实业务, 不知道什么时候需要程序主动去清理清理这些信息. 往往依赖资深开放手工操作

func IsStreamGroupExists(err error) bool {
	// XGROUP CREATE task_stream worker_group [0 or $]
	// (error) BUSYGROUP Consumer Group name already exists
	return strings.HasPrefix(err.Error(), "BUSYGROUP")
}

func (r *Client) XDel(ctx context.Context, stream string, ids ...string) (err error) {
	result, err := r.RDB().XDel(ctx, stream, ids...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "r.RDB().XDel error",
			"Stream", stream, "ids", ids, "err", err, "result", result)
		return err
	}

	return
}

func (r *Client) XAck(ctx context.Context, stream, group string, ids ...string) (err error) {
	result, err := r.RDB().XAck(ctx, stream, group, ids...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "q.R.XAck error",
			"Stream", stream, "Group", group, "ids", ids, "err", err, "result", result)
		return err
	}

	return
}

func (r *Client) XReadGroup(ctx context.Context, xreadgroupargs *redis.XReadGroupArgs) (msg redis.XMessage, err error) {
	// 开放 XReadGroup XAck XDel 自行去定义操作

	res, err := r.RDB().XReadGroup(ctx, xreadgroupargs).Result()
	if err != nil {
		slog.ErrorContext(ctx, "r.RDB().XReadGroup error",
			"Streams", xreadgroupargs.Streams, "Group", xreadgroupargs.Group, "Consumer", xreadgroupargs.Consumer, "err", err)
		return
	}
	if len(res) == 0 || len(res[0].Messages) == 0 {
		slog.InfoContext(ctx, "r.RDB().XReadGroup returned no message",
			"Streams", xreadgroupargs.Streams, "Group", xreadgroupargs.Group, "Consumer", xreadgroupargs.Consumer, "err", err)
		return
	}

	msg = res[0].Messages[0]
	return
}
