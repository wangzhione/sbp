package rediser

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/redis/go-redis/v9"
)

// 删除整个 Stream：
//
//	err := r.Del(ctx, stream).Err()
//
// 删除某个 Consumer Group：
//
//	err := r.Do(ctx, "XGROUP", "DESTROY", stream, group).Err()
//
// 一般业务程序不主动删除整个 Stream 或 Consumer Group。
// 通常由运维人员根据消息保留策略进行清理。
//
// Stream 中的历史消息建议通过 MAXLEN、XTRIM 等方式统一管理，
// 不建议每个 Consumer 在消费完成后直接删除消息。

// IsStreamGroupExists 判断 XGROUP CREATE 是否因为 Group 已经存在而失败。
//
// 使用示例：
//
//	err := r.XGroupCreateMkStream(ctx, "task_stream", "worker_group", "0").Err()
//
//	if err != nil && !IsStreamGroupExists(err) {
//		return err
//	}
func IsStreamGroupExists(err error) bool {
	// XGROUP CREATE task_stream worker_group 0
	// (error) BUSYGROUP Consumer Group name already exists
	return strings.HasPrefix(err.Error(), "BUSYGROUP")
}

// XDel 从 Stream 中永久删除指定消息。
//
// 注意：
//
//  1. XDel 是全局删除消息，不是仅删除当前 Consumer 或 Consumer Group 的消息。
//
//  2. 如果同一个 Stream 存在多个 Consumer Group，某个 Group 处理完成后执行
//     XDel，可能导致其他 Group 无法再读取该消息的内容。
//
//  3. XDel 不等同于 XAck：
//     XAck 只负责确认当前 Group 已完成消息处理；
//     XDel 负责从 Stream 中删除原始消息。
//
//  4. 一般消费流程只需要 XAck，不需要在每次 XAck 后调用 XDel。
func (r *Client) XDel(ctx context.Context, stream string, ids ...string) error {
	result, err := r.UniversalClient.XDel(ctx, stream, ids...).Result()
	if err != nil {
		slog.ErrorContext(
			ctx,
			"r.UniversalClient.XDel error",
			"Stream", stream,
			"ids", ids,
			"err", err,
			"result", result,
		)
	}

	return err
}

// XAck 确认当前 Consumer Group 已成功处理指定消息。
//
// 注意：
//
//  1. XAck 会把消息从当前 Group 的 PEL（Pending Entries List）中移除。
//
//  2. XAck 不会删除 Stream 中保存的原始消息。
//
//  3. 应当在业务处理成功之后再调用 XAck。
//
//  4. 如果业务处理失败，不应调用 XAck。消息会继续保留在 PEL 中，
//     后续可以由当前 Consumer 重新读取，或者由其他 Consumer 通过
//     XAutoClaim/XClaim 接管。
//
//  5. XReadGroup、业务处理和 XAck 不是原子操作。程序可能在业务处理完成后、
//     XAck 之前崩溃，从而造成消息再次投递。因此业务处理应当保证幂等。
func (r *Client) XAck(ctx context.Context, stream string, group string, ids ...string) error {
	result, err := r.UniversalClient.XAck(ctx, stream, group, ids...).Result()
	if err != nil {
		slog.ErrorContext(
			ctx,
			"r.UniversalClient.XAck error",
			"Stream", stream,
			"Group", group,
			"ids", ids,
			"err", err,
			"result", result,
		)
	}

	return err
}

// XReadGroup 使用 Consumer Group 读取 Stream 消息。
//
// 使用注意：
//
//  1. XReadGroup 一次可能返回多个 Stream，每个 Stream 又可能返回多条消息。
//     因此使用方必须完整遍历返回值：
//
//     for _, stream := range streams {
//     for _, msg := range stream.Messages {
//     // 处理 msg
//     }
//     }
//
//     不能只处理：
//
//     streams[0].Messages[0]
//
//     Redis 返回消息时，就已经将消息投递给当前 Consumer。即使使用方忽略了
//     返回结果中的部分消息，这些消息仍然会进入当前 Group 的 PEL。
//     使用 ">" 读取新消息时，这些被忽略的 Pending 消息不会自动再次返回。
//
//  2. Streams 参数由“Stream 名称”和“读取 ID”两部分组成。
//
//     读取一个 Stream：
//
//     []string{"task_stream", ">"}
//
//     读取两个 Stream：
//
//     []string{
//     "stream1", "stream2",
//     ">",       ">",
//     }
//
//     前半部分是 Stream 名称，后半部分是与每个 Stream 对应的读取 ID。
//
//  3. ID 使用 ">" 表示读取从未投递给当前 Consumer Group 的新消息。
//
//     []string{"task_stream", ">"}
//
//     消息返回后，会记录到当前 Group 的 PEL 中，直到使用方调用 XAck。
//
//  4. ID 使用 "0" 或其他具体 ID 时，读取的是已经分配给当前 Consumer、
//     但是尚未 XAck 的 Pending 消息。
//
//     []string{"task_stream", "0"}
//
//     这种方式只能读取属于当前 Consumer 的 Pending 消息，不能直接读取
//     其他 Consumer 名下的 Pending 消息。
//
//     如果某个 Consumer 已经宕机，需要由其他 Consumer 接管其 Pending 消息，
//     应当使用 XAutoClaim 或 XClaim。
//
//  5. Count 表示每个 Stream 一次最多返回多少条消息。
//
//     Count: 10
//
//     如果只读取一个 Stream，一次最多返回10条。
//     如果同时读取两个 Stream，一次最多可能返回20条。
//     Count 是最大数量，Redis 不保证每次一定返回这么多条。
//
//  6. Block 表示没有新消息时的阻塞等待时间。
//
//     Block: 5 * time.Second
//
//     表示最多等待5秒。等待超时且没有消息时，go-redis 返回 redis.Nil。
//     本方法会把 redis.Nil 当作正常的“暂无消息”，返回空结果和 nil error。
//
//     Block: 0
//
//     表示一直阻塞等待消息。为了方便程序关闭和定期检查 Context，业务代码
//     通常更适合使用一个有限的 Block 时间。
//
//  7. NoAck 默认为 false。
//
//     NoAck: false
//
//     表示消息读取后进入 PEL，业务处理成功后必须调用 XAck。
//
//     NoAck: true
//
//     表示消息读取时直接视为已经确认，不会进入 PEL。如果程序在业务处理过程中
//     崩溃，该消息可能无法重试。只有明确允许消息丢失时才应使用 NoAck=true。
//
//  8. 推荐的消息处理顺序是：
//
//     XReadGroup
//     -> 业务处理
//     -> 业务处理成功
//     -> XAck
//
//     业务处理失败时不要 XAck，让消息继续保留在 PEL 中。
//
//  9. Consumer 名称应当能够唯一标识一个消费者实例。
//
//     多台机器或多个进程不应长期共用同一个 Consumer 名称，否则难以区分
//     Pending 消息实际属于哪个消费者实例。
//
//  10. XReadGroup 和 XAck 提供的是至少一次投递语义，而不是严格的只处理一次。
//     业务代码应使用消息 ID、业务订单号或其他唯一键实现幂等处理。
func (r *Client) XReadGroup(ctx context.Context, xreadgroupargs *redis.XReadGroupArgs) (streams []redis.XStream, err error) {
	// XReadGroup 可能批量返回多个 Stream 和多条消息。
	// 这里返回完整的 []redis.XStream，由使用方负责遍历、处理和 XAck。
	streams, err = r.UniversalClient.XReadGroup(ctx, xreadgroupargs).Result()

	// 使用有限 Block 时，等待超时且没有读取到消息会返回 redis.Nil。
	// 这是正常情况，不需要记录错误日志。
	if errors.Is(err, redis.Nil) {
		err = nil
		return
	}

	if err != nil {
		slog.ErrorContext(
			ctx,
			"r.UniversalClient.XReadGroup error",
			"Streams", xreadgroupargs.Streams,
			"Group", xreadgroupargs.Group,
			"Consumer", xreadgroupargs.Consumer,
			"err", err,
		)
		return
	}

	return
}
