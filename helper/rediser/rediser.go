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

// Set 设置 key-value
// Set Redis `SET key value [expiration]` command.
// Use expiration for `SETEx`-like behavior.
//
// Zero expiration means the key has no expiration time.
// KeepTTL is a Redis KEEPTTL option to keep existing TTL, it requires your redis-server version >= 6.0,
// otherwise you will receive an error: (error) ERR syntax error.
func (r *Client) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	err := r.UniversalClient.Set(ctx, key, value, expiration).Err()
	if err != nil {
		slog.ErrorContext(ctx, "Redis Set error", slog.String("key", key), slog.String("error", err.Error()))
	}
	return err
}

// Get 获取 key 的值
func (r *Client) Get(ctx context.Context, key string) (value string, ok bool, err error) {
	value, err = r.UniversalClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			slog.InfoContext(ctx, "Redis Get key not found", slog.String("key", key), slog.String("error", err.Error()))
			return "", false, nil
		}
		slog.ErrorContext(ctx, "Redis Get key error", slog.String("key", key), slog.String("error", err.Error()))
		return
	}

	// 找到了, 设置 ok = true
	return value, true, nil
}

// Del 删除 key
func (r *Client) Del(ctx context.Context, key string) error {
	err := r.UniversalClient.Del(ctx, key).Err()
	if err != nil {
		slog.ErrorContext(ctx, "Redis Del error", slog.String("key", key), slog.String("error", err.Error()))
	}
	return err
}

// Exists 检查 key 是否存在
func (r *Client) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.UniversalClient.Exists(ctx, key).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis Exists error", slog.String("key", key), slog.String("error", err.Error()))
		return false, err
	}
	return count > 0, nil
}

// Expire 设置 key 过期时间
func (r *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	err := r.UniversalClient.Expire(ctx, key, expiration).Err()
	if err != nil {
		slog.ErrorContext(ctx, "Redis Expire error", slog.String("key", key), slog.String("error", err.Error()))
	}
	return err
}

// TTL 获取 key 剩余存活时间
func (r *Client) TTL(ctx context.Context, key string) (ttl time.Duration, err error) {
	ttl, err = r.UniversalClient.TTL(ctx, key).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis TTL error", slog.String("key", key), slog.String("error", err.Error()))
	}
	return
}

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

// Incr 原子递增 key 对应的数值, 返回递增后的值
func (r *Client) Incr(ctx context.Context, key string) (val int64, err error) {
	val, err = r.UniversalClient.Incr(ctx, key).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis Incr failed", slog.String("key", key), slog.String("error", err.Error()))
	}
	return
}

// IncrBy 递增指定值, 返回递增后的值
func (r *Client) IncrBy(ctx context.Context, key string, increment int64) (val int64, err error) {
	val, err = r.UniversalClient.IncrBy(ctx, key, increment).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis IncrBy failed", slog.String("key", key), slog.Int64("increment", increment), slog.String("error", err.Error()))
	}
	return
}

// Decr 原子递减 key 对应的数值, 返回递减后的值
func (r *Client) Decr(ctx context.Context, key string) (val int64, err error) {
	val, err = r.UniversalClient.Decr(ctx, key).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis Decr failed", slog.String("key", key), slog.String("error", err.Error()))
	}
	return
}

// DecrBy 递减指定值, 返回递减后的值
// 如果 Redis 中 key 不存在, 执行 DECR 或 DECRBY 命令时, Redis 会自动创建 key 并初始化为 0, 然后执行递减操作。
func (r *Client) DecrBy(ctx context.Context, key string, decrement int64) (val int64, err error) {
	val, err = r.UniversalClient.DecrBy(ctx, key, decrement).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis DecrBy failed", slog.String("key", key), slog.Int64("decrement", decrement), slog.String("error", err.Error()))
	}
	return
}

// GetSet 设置新值, 并返回旧值
func (r *Client) GetSet(ctx context.Context, key string, value any) (old string, err error) {
	old, err = r.UniversalClient.GetSet(ctx, key, value).Result()
	if err != nil {
		if err == redis.Nil {
			slog.InfoContext(ctx, "Redis GetSet key not found", slog.String("key", key))
			// 吃掉首次设置, 返回 redis.Nil 情况
			return "", nil
		}
		slog.ErrorContext(ctx, "Redis GetSet failed", slog.String("key", key), slog.String("error", err.Error()))
	}
	return
}

// SetNX 只有 key 不存在时才会设置值
func (r *Client) SetNX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error) {
	success, err := r.UniversalClient.SetNX(ctx, key, value, expiration).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis SetNX failed", slog.String("key", key), slog.String("error", err.Error()))
	}
	return success, err
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
