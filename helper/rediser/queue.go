package rediser

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/wangzhione/sbp/chain"
	"github.com/wangzhione/sbp/structs"
)

// stream 模拟 分布式 queue

// Queue represents a Redis Stream task queue with single group.
type Queue struct {
	R        *Client // *redis.Client
	Stream   string
	Group    string
	Consumer string
	MaxLen   int64 // 默认 0, 无限
}

func IsStreamGroupExists(err error) bool {
	// XGROUP CREATE task_stream worker_group [0 or $]
	// (error) BUSYGROUP Consumer Group name already exists
	return strings.HasPrefix(err.Error(), "BUSYGROUP")
}

// NewQueue initializes the stream queue, ensuring stream & group exist.
// maxLen 默认填写 0
func (r *Client) NewQueue(ctx context.Context, stream, group string, maxLen ...int64) (q *Queue, err error) {
	result, err := r.XGroupCreateMkStream(ctx, stream, group, "0").Result()
	if err != nil {
		if IsStreamGroupExists(err) {
			// 如果提示已经创建了 Group 默认吃掉这个 error
			err = nil
		} else {
			slog.ErrorContext(ctx, "XGroupCreateMkStream stream group error",
				"Stream", stream, "Group", group, "MaxLen", maxLen, "result", result)
			return
		}
	}

	consumer := chain.Hostname + "." + chain.UUID()

	// 没有错误, 或者 group 已经存在
	q = &Queue{
		R:        r,
		Stream:   stream,
		Group:    group,
		Consumer: consumer,
		MaxLen:   structs.Max(maxLen...),
	}

	return
}

// Produce pushes a new task into the stream. return insert stream id
func (q *Queue) Produce(ctx context.Context, values map[string]any) (msgID string, err error) {
	xaddargs := &redis.XAddArgs{
		Stream: q.Stream,
		MaxLen: q.MaxLen, // MaxLen = 0, Redis 会一直保留所有历史消息, Stream 会无限增长, 不会触发裁剪策略
		Approx: true,     // 默认 MaxLen + Approx 策略, 近似修剪（~）
		Values: values,
	}

	// XAddArgs.Values 支持以下格式：
	// - map[string]any{"k1": "v1", "k2": "v2"} ✅ 推荐
	// - []any{"k1", "v1", "k2", "v2"}          ✅ 自定义顺序
	// - []string{"k1", "v1", "k2", "v2"}       ✅ 简洁写法

	/*
	   XAddArgs.Values 类型	是否支持	示例值
	   string	✅	"hello"
	   []byte	✅	[]byte(\"binary\")
	   int, int64, float64	✅	123, 45.6
	   bool	✅	true, false
	   时间类型（如 time.Time）	✅	自动转换为字符串
	   任意可被 fmt.Sprint 转换为字符串的值	✅	自动调用内部序列化
	*/

	msgID, err = q.R.XAdd(ctx, xaddargs).Result()
	if err != nil {
		slog.ErrorContext(ctx, "q.R.XAdd error",
			"Stream", q.Stream, "Group", q.Group, "Consumer", q.Consumer, "err", err)
		return
	}
	return
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

func (q *Queue) XAck(ctx context.Context, msgID string) (err error) {
	err = q.R.XAck(ctx, q.Stream, q.Group, msgID)
	if err != nil {
		slog.ErrorContext(ctx, "q.R.XAck error", "err", err, "Consumer", q.Consumer)
	}

	// 在 Queue 中 ack 应答是 集合 delete 业务一起的
	err = q.R.XDel(ctx, q.Stream, msgID)
	if err != nil {
		slog.ErrorContext(ctx, "q.R.XDel error",
			"Group", q.Group, "Consumer", q.Consumer, "err", err)
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

// Consume reads one task and calls handler, then ACK + DEL.
// block time.Duration  默认 -1 无限等待数据到来; 0 zero, 有无结果都立即返回 XReadGroup
func (q *Queue) Consume(ctx context.Context, block time.Duration, handler func(values map[string]any) error) (err error) {
	xreadgroupargs := &redis.XReadGroupArgs{
		Group:    q.Group,
		Consumer: q.Consumer,
		Streams:  []string{q.Stream, ">"}, // 从 q.Stream 这个 Stream 中，读取 q.Group 尚未读取的新消息
		Count:    1,
		Block:    block,
	}

	// 🚨 注意：BLOCK 0（协议） ⬌ Block: -1（go-redis）
	// XREADGROUP GROUP mygroup consumer-name STREAMS mystream > BLOCK 0
	// BLOCK 0 就是 无限阻塞
	// BLOCK 5000 表示最多阻塞 5 秒（超时返回 nil）

	msg, err := q.R.XReadGroup(ctx, xreadgroupargs)
	if err != nil {
		return
	}

	// msg.ID = Queue.Produce msgID
	slog.InfoContext(ctx, "Consume handler begin", "msgID", msg.ID, "values", msg.Values)
	defer func() {
		slog.InfoContext(ctx, "Consume handler end", "msgID", msg.ID, "reason", err)
	}()
	if err := handler(msg.Values); err != nil {
		slog.ErrorContext(ctx, "Consume handler end error",
			"Stream", q.Stream, "Group", q.Group, "Consumer", q.Consumer, "msgID", msg.ID, "values", msg.Values, "err", err)
		return err
	}

	// XReadGroup -> XAck 随后 清理 stream 中 msg.ID
	return q.XAck(ctx, msg.ID)
}

// 删除 Stream（即整个消息队列）
// err := r.RDB().Del(ctx, stream).Err()
// 删除某个 Group
// err := r.RDB().Do(ctx, "XGROUP", "DESTROY", stream, group).Err()
// 一般现实业务, 不知道什么时候需要程序主动去清理清理这些信息. 往往依赖资深开放手工操作
