package rediser

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/wangzhione/sbp/structs"
	"github.com/wangzhione/sbp/system"
)

// stream 模拟 分布式 queue

// Queue represents a Redis Stream task queue with single group.
// Queue 内部设计, 默认给服务做简单解耦操作, 不是消息发布和订阅, 而是类似 任务队列概念, 发布任务, 执行任务, 任务执行完毕
type Queue struct {
	R        *Client // *rediser.Client
	Stream   string
	Group    string
	Consumer string
	MaxLen   int64 // 默认 0, 无限
}

func (q *Queue) Init(ctx context.Context) (err error) {
	if q.Consumer == "" {
		// 内部定义启动这个 队列 随后 Queue.Consume 发给 redis 的消费者名称
		q.Consumer = system.Hostname + "." + system.UUID()[:6]
	}

	if q.Group == "" {
		q.Group = q.Stream
	}

	result, err := q.R.XGroupCreateMkStream(ctx, q.Stream, q.Group, "0").Result()
	if err != nil {
		if IsStreamGroupExists(err) {
			// 如果提示已经创建了 Group 默认吃掉这个 error
			err = nil
		} else {
			slog.ErrorContext(ctx, "XGroupCreateMkStream stream group error",
				"Stream", q.Stream, "Group", q.Group, "MaxLen", q.MaxLen, "result", result)
			return
		}
	}

	return
}

// NewQueue initializes the {name} stream queue, ensuring stream & group exist.
// maxLen 默认不填写 , 默认设置为 0 , 这个 queue 理论上不受长度限制
// 有 maxLen 当超长时候, 会丢弃早期消息
func (r *Client) NewQueue(ctx context.Context, name string, maxLen ...int64) (*Queue, error) {
	// 没有错误, 或者 group 已经存在
	q := &Queue{
		R:      r,
		Stream: name,
		MaxLen: structs.Max(maxLen...),
	}

	err := q.Init(ctx)
	if err != nil {
		return nil, err
	}
	return q, err
}

// Produce pushes a new task into the stream. return insert stream id
func (q *Queue) Produce(ctx context.Context, values map[string]any) (msgID string, err error) {
	xaddargs := &redis.XAddArgs{
		Stream: q.Stream,
		MaxLen: q.MaxLen, // MaxLen = 0, Redis 会一直保留所有历史消息, Stream 会无限增长, 不会触发裁剪策略
		Approx: true,     // 默认 MaxLen + Approx 策略, 近似修剪（~）🧹 删除规则：从最早的消息开始（左边裁剪）
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

// Ack Acknowledge character 执行完任务消息, 之后需要应答, 通知远端消费完成
func (q *Queue) Ack(ctx context.Context, msgID string) (err error) {
	err = q.R.XAck(ctx, q.Stream, q.Group, msgID)
	if err != nil {
		slog.ErrorContext(ctx, "q.R.XAck panic error", "err", err, "Consumer", q.Consumer)
	}

	// 在 Queue 中 ack 应答是 集合 delete 业务一起的
	err = q.R.XDel(ctx, q.Stream, msgID)
	if err != nil {
		slog.ErrorContext(ctx, "q.R.XDel panic error",
			"Group", q.Group, "Consumer", q.Consumer, "err", err)
		return err
	}

	return
}

// Consume reads one task and calls handler, then ACK + DEL.
// block time.Duration  默认 -1 无限等待数据到来; 0 zero, 有无结果都立即返回 XReadGroup
func (q *Queue) Consume(ctx context.Context, block time.Duration, consume func(values map[string]any) error) (err error) {
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

	// 默认 return err != nil, 消费失败, 不 XAck + XDel
	if err := consume(msg.Values); err != nil {
		slog.ErrorContext(ctx, "Consume consume handler end error",
			"Stream", q.Stream, "Group", q.Group, "Consumer", q.Consumer, "msgID", msg.ID, "values", msg.Values, "err", err)
		return err
	}

	// XReadGroup -> XAck 随后 清理 stream 中 msg.ID
	return q.Ack(ctx, msg.ID)
}
