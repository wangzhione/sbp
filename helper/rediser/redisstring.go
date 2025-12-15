package rediser

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

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

// Append 将 value 追加到 key 的值的末尾
// 如果 key 不存在, 则创建一个新的 key 并设置 value
func (r *Client) Append(ctx context.Context, key, value string) (length int64, err error) {
	length, err = r.UniversalClient.Append(ctx, key, value).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis Append failed", slog.String("key", key), slog.String("error", err.Error()))
	}
	return
}

// GetRange 获取 key 中字符串的子串, 范围由 start 和 end 指定 (包含两端)
// start 和 end 可以是负数, 表示从字符串末尾开始计算
func (r *Client) GetRange(ctx context.Context, key string, start, end int64) (value string, err error) {
	value, err = r.UniversalClient.GetRange(ctx, key, start, end).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis GetRange failed", slog.String("key", key), slog.Int64("start", start), slog.Int64("end", end), slog.String("error", err.Error()))
	}
	return
}

// GetEx 获取 key 的值, 并设置过期时间
// 如果 key 不存在, 返回 redis.Nil 错误
func (r *Client) GetEx(ctx context.Context, key string, expiration time.Duration) (value string, ok bool, err error) {
	value, err = r.UniversalClient.GetEx(ctx, key, expiration).Result()
	if err != nil {
		if err == redis.Nil {
			slog.InfoContext(ctx, "Redis GetEx key not found", slog.String("key", key))
			return "", false, nil
		}
		slog.ErrorContext(ctx, "Redis GetEx failed", slog.String("key", key), slog.String("error", err.Error()))
		return
	}
	return value, true, nil
}

// GetDel 获取 key 的值, 并删除该 key
// 如果 key 不存在, 返回 redis.Nil 错误
func (r *Client) GetDel(ctx context.Context, key string) (value string, ok bool, err error) {
	value, err = r.UniversalClient.GetDel(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			slog.InfoContext(ctx, "Redis GetDel key not found", slog.String("key", key))
			return "", false, nil
		}
		slog.ErrorContext(ctx, "Redis GetDel failed", slog.String("key", key), slog.String("error", err.Error()))
		return
	}
	return value, true, nil
}

// IncrByFloat 将 key 中存储的数字加上指定的浮点增量值, 返回递增后的值
// 如果 key 不存在, 会先初始化为 0 再执行操作
func (r *Client) IncrByFloat(ctx context.Context, key string, increment float64) (val float64, err error) {
	val, err = r.UniversalClient.IncrByFloat(ctx, key, increment).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis IncrByFloat failed", slog.String("key", key), slog.Float64("increment", increment), slog.String("error", err.Error()))
	}
	return
}

// MGet 批量获取多个 key 的值
// 返回一个切片, 顺序与输入的 keys 顺序一致
// 如果某个 key 不存在, 对应的值为空字符串
func (r *Client) MGet(ctx context.Context, keys ...string) (values []string, err error) {
	results, err := r.UniversalClient.MGet(ctx, keys...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis MGet failed", slog.Any("keys", keys), slog.String("error", err.Error()))
		return nil, err
	}
	// 将 []interface{} 转换为 []string, nil 值转换为空字符串
	values = make([]string, len(results))
	for i, v := range results {
		if v == nil {
			values[i] = ""
		} else {
			values[i] = v.(string)
		}
	}
	return values, nil
}

// MSet 批量设置多个 key-value 对
// values 参数应该是 key1, value1, key2, value2, ... 的形式
func (r *Client) MSet(ctx context.Context, values ...any) error {
	err := r.UniversalClient.MSet(ctx, values...).Err()
	if err != nil {
		slog.ErrorContext(ctx, "Redis MSet failed", slog.Any("values", values), slog.String("error", err.Error()))
	}
	return err
}

// MSetNX 批量设置多个 key-value 对, 仅当所有 key 都不存在时才设置
// 如果任何一个 key 已存在, 则所有 key 都不会被设置
// values 参数应该是 key1, value1, key2, value2, ... 的形式
// 返回 true 表示所有 key 都设置成功, false 表示至少有一个 key 已存在
func (r *Client) MSetNX(ctx context.Context, values ...any) (bool, error) {
	success, err := r.UniversalClient.MSetNX(ctx, values...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis MSetNX failed", slog.Any("values", values), slog.String("error", err.Error()))
	}
	return success, err
}

// SetEx 设置 key 的值, 并设置过期时间
// 这是一个原子操作, 等同于 SET key value EX seconds
func (r *Client) SetEx(ctx context.Context, key string, value any, expiration time.Duration) error {
	err := r.UniversalClient.SetEx(ctx, key, value, expiration).Err()
	if err != nil {
		slog.ErrorContext(ctx, "Redis SetEx failed", slog.String("key", key), slog.String("error", err.Error()))
	}
	return err
}

// SetXX 仅当 key 已存在时才设置值
// 如果 key 不存在, 操作不会执行, 返回 false
// 如果 key 存在, 设置成功返回 true
func (r *Client) SetXX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error) {
	success, err := r.UniversalClient.SetXX(ctx, key, value, expiration).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis SetXX failed", slog.String("key", key), slog.String("error", err.Error()))
	}
	return success, err
}

// SetRange 从 offset 位置开始, 用 value 覆盖 key 中存储的字符串
// 如果 key 不存在, 会先创建一个空字符串, 然后用 value 填充
// 返回修改后字符串的长度
func (r *Client) SetRange(ctx context.Context, key string, offset int64, value string) (length int64, err error) {
	length, err = r.UniversalClient.SetRange(ctx, key, offset, value).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis SetRange failed", slog.String("key", key), slog.Int64("offset", offset), slog.String("error", err.Error()))
	}
	return
}

// StrLen 获取 key 中存储的字符串的长度
// 如果 key 不存在, 返回 0
// 时间复杂度: O(1), Redis 内部维护了字符串长度信息, 无需遍历字符串
func (r *Client) StrLen(ctx context.Context, key string) (length int64, err error) {
	length, err = r.UniversalClient.StrLen(ctx, key).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis StrLen failed", slog.String("key", key), slog.String("error", err.Error()))
	}
	return
}
