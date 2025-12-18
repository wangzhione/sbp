package rediser

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

// LPush 向列表左侧插入值
// 返回插入后列表的长度
func (r *Client) LPush(ctx context.Context, key string, values ...any) (length int64, err error) {
	length, err = r.UniversalClient.LPush(ctx, key, values...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis LPush error", slog.String("key", key), slog.String("error", err.Error()))
		return
	}
	return length, nil
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

// BLMPop 从多个列表左侧阻塞式弹出指定数量的元素
// direction 为 "LEFT" 或 "RIGHT"，count 为要弹出的元素数量
// 返回键名和对应的值列表
func (r *Client) BLMPop(ctx context.Context, timeout time.Duration, direction string, count int64, keys ...string) (key string, values []string, err error) {
	key, values, err = r.UniversalClient.BLMPop(ctx, timeout, direction, count, keys...).Result()
	if err != nil {
		if err == redis.Nil {
			slog.InfoContext(ctx, "Redis BLMPop timeout or empty", slog.String("direction", direction), slog.Int64("count", count), slog.Any("keys", keys))
			return "", nil, nil
		}
		slog.ErrorContext(ctx, "Redis BLMPop failed", slog.String("direction", direction), slog.Int64("count", count), slog.Any("keys", keys), slog.String("error", err.Error()))
		return "", nil, err
	}
	return key, values, nil
}

// BRPopLPush 从源列表右侧阻塞式弹出一个元素，并将其推入目标列表的左侧
// 如果源列表为空，会阻塞直到有元素可用或超时
// 返回被移动的元素值
func (r *Client) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) (value string, ok bool, err error) {
	value, err = r.UniversalClient.BRPopLPush(ctx, source, destination, timeout).Result()
	if err != nil {
		if err == redis.Nil {
			slog.InfoContext(ctx, "Redis BRPopLPush timeout or empty", slog.String("source", source), slog.String("destination", destination))
			return "", false, nil
		}
		slog.ErrorContext(ctx, "Redis BRPopLPush failed", slog.String("source", source), slog.String("destination", destination), slog.String("error", err.Error()))
		return "", false, err
	}
	return value, true, nil
}

// LIndex 通过索引获取列表中的元素
// index 从 0 开始，负数表示从列表末尾开始计数
// 返回指定索引位置的元素值
func (r *Client) LIndex(ctx context.Context, key string, index int64) (value string, ok bool, err error) {
	value, err = r.UniversalClient.LIndex(ctx, key, index).Result()
	if err != nil {
		if err == redis.Nil {
			slog.InfoContext(ctx, "Redis LIndex key empty or index out of range", slog.String("key", key), slog.Int64("index", index))
			return "", false, nil
		}
		slog.ErrorContext(ctx, "Redis LIndex error", slog.String("key", key), slog.Int64("index", index), slog.String("error", err.Error()))
		return "", false, err
	}
	return value, true, nil
}

// LInsert 在列表的指定元素前或后插入新元素
// op 为 "BEFORE" 或 "AFTER"，pivot 为参考元素，value 为要插入的值
// 返回插入后列表的长度，如果 pivot 不存在返回 -1
func (r *Client) LInsert(ctx context.Context, key, op string, pivot, value any) (length int64, err error) {
	length, err = r.UniversalClient.LInsert(ctx, key, op, pivot, value).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis LInsert error",
			slog.String("key", key),
			slog.String("op", op),
			slog.Any("pivot", pivot),
			slog.Any("value", value),
			slog.String("error", err.Error()))
		return
	}
	return length, nil
}

// LInsertBefore 在列表的指定元素前插入新元素
// pivot 为参考元素，value 为要插入的值
// 返回插入后列表的长度，如果 pivot 不存在返回 -1
func (r *Client) LInsertBefore(ctx context.Context, key string, pivot, value any) (length int64, err error) {
	length, err = r.UniversalClient.LInsertBefore(ctx, key, pivot, value).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis LInsertBefore error",
			slog.String("key", key),
			slog.Any("pivot", pivot),
			slog.Any("value", value),
			slog.String("error", err.Error()))
		return
	}
	return length, nil
}

// LInsertAfter 在列表的指定元素后插入新元素
// pivot 为参考元素，value 为要插入的值
// 返回插入后列表的长度，如果 pivot 不存在返回 -1
func (r *Client) LInsertAfter(ctx context.Context, key string, pivot, value any) (length int64, err error) {
	length, err = r.UniversalClient.LInsertAfter(ctx, key, pivot, value).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis LInsertAfter error",
			slog.String("key", key),
			slog.Any("pivot", pivot),
			slog.Any("value", value),
			slog.String("error", err.Error()))
		return
	}
	return length, nil
}

// LLen 获取列表的长度
// 如果 key 不存在，返回 0
func (r *Client) LLen(ctx context.Context, key string) (length int64, err error) {
	length, err = r.UniversalClient.LLen(ctx, key).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis LLen error", slog.String("key", key), slog.String("error", err.Error()))
		return
	}
	return length, nil
}

// LMPop 从多个列表左侧或右侧弹出指定数量的元素
// direction 为 "LEFT" 或 "RIGHT"，count 为要弹出的元素数量
// 返回键名和对应的值列表
func (r *Client) LMPop(ctx context.Context, direction string, count int64, keys ...string) (key string, values []string, err error) {
	key, values, err = r.UniversalClient.LMPop(ctx, direction, count, keys...).Result()
	if err != nil {
		if err == redis.Nil {
			slog.InfoContext(ctx, "Redis LMPop empty", slog.String("direction", direction), slog.Int64("count", count), slog.Any("keys", keys))
			return "", nil, nil
		}
		slog.ErrorContext(ctx, "Redis LMPop failed", slog.String("direction", direction), slog.Int64("count", count), slog.Any("keys", keys), slog.String("error", err.Error()))
		return "", nil, err
	}
	return key, values, nil
}

// LPop 从列表左侧弹出值
// 返回弹出的元素值
func (r *Client) LPop(ctx context.Context, key string) (value string, ok bool, err error) {
	value, err = r.UniversalClient.LPop(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			slog.InfoContext(ctx, "Redis LPop key empty", slog.String("key", key))
			return "", false, nil
		}
		slog.ErrorContext(ctx, "Redis LPop error", slog.String("key", key), slog.String("error", err.Error()))
		return "", false, err
	}
	return value, true, nil
}

// LPopCount 从列表左侧弹出指定数量的元素
// count 为要弹出的元素数量
// 返回弹出的元素列表
func (r *Client) LPopCount(ctx context.Context, key string, count int) (values []string, err error) {
	values, err = r.UniversalClient.LPopCount(ctx, key, count).Result()
	if err != nil {
		if err == redis.Nil {
			slog.InfoContext(ctx, "Redis LPopCount key empty", slog.String("key", key), slog.Int("count", count))
			return nil, nil
		}
		slog.ErrorContext(ctx, "Redis LPopCount error", slog.String("key", key), slog.Int("count", count), slog.String("error", err.Error()))
		return nil, err
	}
	return values, nil
}

// LPos 查找列表中指定元素的位置
// value 为要查找的元素值，args 为查找选项（如 RANK、MAXLEN 等）
// 返回元素在列表中的索引位置
func (r *Client) LPos(ctx context.Context, key string, value string, args redis.LPosArgs) (index int64, ok bool, err error) {
	index, err = r.UniversalClient.LPos(ctx, key, value, args).Result()
	if err != nil {
		if err == redis.Nil {
			slog.InfoContext(ctx, "Redis LPos value not found", slog.String("key", key), slog.String("value", value))
			return 0, false, nil
		}
		slog.ErrorContext(ctx, "Redis LPos error", slog.String("key", key), slog.String("value", value), slog.String("error", err.Error()))
		return 0, false, err
	}
	return index, true, nil
}

// LPosCount 查找列表中指定元素的所有位置
// value 为要查找的元素值，count 为要返回的最大位置数量，args 为查找选项
// 返回元素在列表中的所有索引位置
func (r *Client) LPosCount(ctx context.Context, key string, value string, count int64, args redis.LPosArgs) (indices []int64, err error) {
	indices, err = r.UniversalClient.LPosCount(ctx, key, value, count, args).Result()
	if err != nil {
		if err == redis.Nil {
			slog.InfoContext(ctx, "Redis LPosCount value not found", slog.String("key", key), slog.String("value", value), slog.Int64("count", count))
			return nil, nil
		}
		slog.ErrorContext(ctx, "Redis LPosCount error", slog.String("key", key), slog.String("value", value), slog.Int64("count", count), slog.String("error", err.Error()))
		return nil, err
	}
	return indices, nil
}

// LPushX 仅当列表存在时，向列表左侧插入值
// 如果列表不存在，操作不会执行
// 返回插入后列表的长度
func (r *Client) LPushX(ctx context.Context, key string, values ...any) (length int64, err error) {
	length, err = r.UniversalClient.LPushX(ctx, key, values...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis LPushX error", slog.String("key", key), slog.String("error", err.Error()))
		return
	}
	return length, nil
}

// LRem 从列表中移除指定数量的元素
// count > 0: 从表头开始搜索，移除 count 个与 value 相等的元素
// count < 0: 从表尾开始搜索，移除 |count| 个与 value 相等的元素
// count = 0: 移除所有与 value 相等的元素
// 返回被移除元素的数量
func (r *Client) LRem(ctx context.Context, key string, count int64, value any) (removed int64, err error) {
	removed, err = r.UniversalClient.LRem(ctx, key, count, value).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis LRem error",
			slog.String("key", key),
			slog.Int64("count", count),
			slog.Any("value", value),
			slog.String("error", err.Error()))
		return
	}
	return removed, nil
}

// LSet 设置列表中指定索引位置的元素值
// index 为索引位置，value 为要设置的值
// 如果索引超出范围，会返回错误
func (r *Client) LSet(ctx context.Context, key string, index int64, value any) error {
	err := r.UniversalClient.LSet(ctx, key, index, value).Err()
	if err != nil {
		slog.ErrorContext(ctx, "Redis LSet error",
			slog.String("key", key),
			slog.Int64("index", index),
			slog.Any("value", value),
			slog.String("error", err.Error()))
		return err
	}
	return nil
}

// LTrim 修剪列表，只保留指定区间内的元素
// start 和 stop 为区间范围，包含两端
// 区间外的元素会被删除
func (r *Client) LTrim(ctx context.Context, key string, start, stop int64) error {
	err := r.UniversalClient.LTrim(ctx, key, start, stop).Err()
	if err != nil {
		slog.ErrorContext(ctx, "Redis LTrim error",
			slog.String("key", key),
			slog.Int64("start", start),
			slog.Int64("stop", stop),
			slog.String("error", err.Error()))
		return err
	}
	return nil
}

// RPopCount 从列表右侧弹出指定数量的元素
// count 为要弹出的元素数量
// 返回弹出的元素列表
func (r *Client) RPopCount(ctx context.Context, key string, count int) (values []string, err error) {
	values, err = r.UniversalClient.RPopCount(ctx, key, count).Result()
	if err != nil {
		if err == redis.Nil {
			slog.InfoContext(ctx, "Redis RPopCount key empty", slog.String("key", key), slog.Int("count", count))
			return nil, nil
		}
		slog.ErrorContext(ctx, "Redis RPopCount error", slog.String("key", key), slog.Int("count", count), slog.String("error", err.Error()))
		return nil, err
	}
	return values, nil
}

// RPopLPush 从源列表右侧弹出一个元素，并将其推入目标列表的左侧
// 这是一个原子操作
// 返回被移动的元素值
func (r *Client) RPopLPush(ctx context.Context, source, destination string) (value string, ok bool, err error) {
	value, err = r.UniversalClient.RPopLPush(ctx, source, destination).Result()
	if err != nil {
		if err == redis.Nil {
			slog.InfoContext(ctx, "Redis RPopLPush source empty", slog.String("source", source), slog.String("destination", destination))
			return "", false, nil
		}
		slog.ErrorContext(ctx, "Redis RPopLPush error", slog.String("source", source), slog.String("destination", destination), slog.String("error", err.Error()))
		return "", false, err
	}
	return value, true, nil
}

// RPush 向列表右侧插入值
// 返回插入后列表的长度
func (r *Client) RPush(ctx context.Context, key string, values ...any) (length int64, err error) {
	length, err = r.UniversalClient.RPush(ctx, key, values...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis RPush error", slog.String("key", key), slog.String("error", err.Error()))
		return
	}
	return length, nil
}

// RPushX 仅当列表存在时，向列表右侧插入值
// 如果列表不存在，操作不会执行
// 返回插入后列表的长度
func (r *Client) RPushX(ctx context.Context, key string, values ...any) (length int64, err error) {
	length, err = r.UniversalClient.RPushX(ctx, key, values...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis RPushX error", slog.String("key", key), slog.String("error", err.Error()))
		return
	}
	return length, nil
}

// LMove 将列表中的元素从一个位置移动到另一个位置
// source 和 destination 可以是同一个列表或不同的列表
// srcpos 和 destpos 为 "LEFT" 或 "RIGHT"，表示源位置和目标位置
// 返回被移动的元素值
func (r *Client) LMove(ctx context.Context, source, destination, srcpos, destpos string) (value string, ok bool, err error) {
	value, err = r.UniversalClient.LMove(ctx, source, destination, srcpos, destpos).Result()
	if err != nil {
		if err == redis.Nil {
			slog.InfoContext(ctx, "Redis LMove source empty",
				slog.String("source", source),
				slog.String("destination", destination),
				slog.String("srcpos", srcpos),
				slog.String("destpos", destpos))
			return "", false, nil
		}
		slog.ErrorContext(ctx, "Redis LMove error",
			slog.String("source", source),
			slog.String("destination", destination),
			slog.String("srcpos", srcpos),
			slog.String("destpos", destpos),
			slog.String("error", err.Error()))
		return "", false, err
	}
	return value, true, nil
}

// BLMove 阻塞式将列表中的元素从一个位置移动到另一个位置
// 如果源列表为空，会阻塞直到有元素可用或超时
// source 和 destination 可以是同一个列表或不同的列表
// srcpos 和 destpos 为 "LEFT" 或 "RIGHT"，表示源位置和目标位置
// timeout 为阻塞超时时间
// 返回被移动的元素值
func (r *Client) BLMove(ctx context.Context, source, destination, srcpos, destpos string, timeout time.Duration) (value string, ok bool, err error) {
	value, err = r.UniversalClient.BLMove(ctx, source, destination, srcpos, destpos, timeout).Result()
	if err != nil {
		if err == redis.Nil {
			slog.InfoContext(ctx, "Redis BLMove timeout or empty",
				slog.String("source", source),
				slog.String("destination", destination),
				slog.String("srcpos", srcpos),
				slog.String("destpos", destpos))
			return "", false, nil
		}
		slog.ErrorContext(ctx, "Redis BLMove error",
			slog.String("source", source),
			slog.String("destination", destination),
			slog.String("srcpos", srcpos),
			slog.String("destpos", destpos),
			slog.String("error", err.Error()))
		return "", false, err
	}
	return value, true, nil
}
