package rediser

import (
	"context"
	"log/slog"
	"time"

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

// HDel 删除哈希表中的一个或多个字段
// 返回被删除字段的数量
func (r *Client) HDel(ctx context.Context, key string, fields ...string) (count int64, err error) {
	count, err = r.UniversalClient.HDel(ctx, key, fields...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HDel error", slog.String("key", key), slog.Any("fields", fields), slog.String("error", err.Error()))
		return
	}
	return count, nil
}

// HExists 检查哈希表中指定字段是否存在
// 返回 true 表示字段存在，false 表示不存在
func (r *Client) HExists(ctx context.Context, key, field string) (exists bool, err error) {
	exists, err = r.UniversalClient.HExists(ctx, key, field).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HExists error", slog.String("key", key), slog.String("field", field), slog.String("error", err.Error()))
		return
	}
	return exists, nil
}

// HGetAll 获取哈希表中所有的字段和值
// 返回一个 map，键为字段名，值为字段值
func (r *Client) HGetAll(ctx context.Context, key string) (result map[string]string, err error) {
	result, err = r.UniversalClient.HGetAll(ctx, key).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HGetAll error", slog.String("key", key), slog.String("error", err.Error()))
		return
	}
	return result, nil
}

// HGetDel 获取哈希表中指定字段的值，然后删除该字段
// 返回被删除字段的值列表
func (r *Client) HGetDel(ctx context.Context, key string, fields ...string) (values []string, err error) {
	values, err = r.UniversalClient.HGetDel(ctx, key, fields...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HGetDel error", slog.String("key", key), slog.Any("fields", fields), slog.String("error", err.Error()))
		return
	}
	return values, nil
}

// HGetEX 获取哈希表中指定字段的值，并设置过期时间
// 返回字段值的列表
func (r *Client) HGetEX(ctx context.Context, key string, fields ...string) (values []string, err error) {
	values, err = r.UniversalClient.HGetEX(ctx, key, fields...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HGetEX error", slog.String("key", key), slog.Any("fields", fields), slog.String("error", err.Error()))
		return
	}
	return values, nil
}

// HGetEXWithArgs 获取哈希表中指定字段的值，并设置过期时间（带参数选项）
// options 包含过期时间等配置选项
// 返回字段值的列表
func (r *Client) HGetEXWithArgs(ctx context.Context, key string, options *redis.HGetEXOptions, fields ...string) (values []string, err error) {
	values, err = r.UniversalClient.HGetEXWithArgs(ctx, key, options, fields...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HGetEXWithArgs error", slog.String("key", key), slog.Any("fields", fields), slog.String("error", err.Error()))
		return
	}
	return values, nil
}

// HIncrBy 将哈希表中指定字段的值增加指定的整数增量
// 返回增加后的新值
func (r *Client) HIncrBy(ctx context.Context, key, field string, incr int64) (val int64, err error) {
	val, err = r.UniversalClient.HIncrBy(ctx, key, field, incr).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HIncrBy error", slog.String("key", key), slog.String("field", field), slog.Int64("incr", incr), slog.String("error", err.Error()))
		return
	}
	return val, nil
}

// HIncrByFloat 将哈希表中指定字段的值增加指定的浮点数增量
// 返回增加后的新值
func (r *Client) HIncrByFloat(ctx context.Context, key, field string, incr float64) (val float64, err error) {
	val, err = r.UniversalClient.HIncrByFloat(ctx, key, field, incr).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HIncrByFloat error", slog.String("key", key), slog.String("field", field), slog.Float64("incr", incr), slog.String("error", err.Error()))
		return
	}
	return val, nil
}

// HKeys 获取哈希表中所有字段名
// 返回字段名的字符串切片
func (r *Client) HKeys(ctx context.Context, key string) (fields []string, err error) {
	fields, err = r.UniversalClient.HKeys(ctx, key).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HKeys error", slog.String("key", key), slog.String("error", err.Error()))
		return
	}
	return fields, nil
}

// HLen 获取哈希表中字段的数量
// 返回字段的数量
func (r *Client) HLen(ctx context.Context, key string) (count int64, err error) {
	count, err = r.UniversalClient.HLen(ctx, key).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HLen error", slog.String("key", key), slog.String("error", err.Error()))
		return
	}
	return count, nil
}

// HMGet 获取哈希表中一个或多个字段的值
// 返回字段值的切片，顺序与输入的字段顺序一致
func (r *Client) HMGet(ctx context.Context, key string, fields ...string) (values []any, err error) {
	values, err = r.UniversalClient.HMGet(ctx, key, fields...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HMGet error", slog.String("key", key), slog.Any("fields", fields), slog.String("error", err.Error()))
		return
	}
	return values, nil
}

// HMSet 同时设置哈希表中多个字段的值
// values 是字段名和值的交替序列，例如: "field1", "value1", "field2", "value2"
// 返回是否设置成功
func (r *Client) HMSet(ctx context.Context, key string, values ...any) (ok bool, err error) {
	ok, err = r.UniversalClient.HMSet(ctx, key, values...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HMSet error", slog.String("key", key), slog.Any("values", values), slog.String("error", err.Error()))
		return
	}
	return ok, nil
}

// HSetEX 设置哈希表中字段的值，并设置过期时间
// fieldsAndValues 是字段名和值的交替序列
// 返回设置的字段数量
func (r *Client) HSetEX(ctx context.Context, key string, fieldsAndValues ...string) (count int64, err error) {
	count, err = r.UniversalClient.HSetEX(ctx, key, fieldsAndValues...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HSetEX error", slog.String("key", key), slog.Any("fieldsAndValues", fieldsAndValues), slog.String("error", err.Error()))
		return
	}
	return count, nil
}

// HSetEXWithArgs 设置哈希表中字段的值，并设置过期时间（带参数选项）
// options 包含过期时间等配置选项
// fieldsAndValues 是字段名和值的交替序列
// 返回设置的字段数量
func (r *Client) HSetEXWithArgs(ctx context.Context, key string, options *redis.HSetEXOptions, fieldsAndValues ...string) (count int64, err error) {
	count, err = r.UniversalClient.HSetEXWithArgs(ctx, key, options, fieldsAndValues...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HSetEXWithArgs error", slog.String("key", key), slog.Any("fieldsAndValues", fieldsAndValues), slog.String("error", err.Error()))
		return
	}
	return count, nil
}

// HSetNX 仅当哈希表中指定字段不存在时，才设置该字段的值
// 返回 true 表示设置成功（字段不存在），false 表示字段已存在
func (r *Client) HSetNX(ctx context.Context, key, field string, value any) (ok bool, err error) {
	ok, err = r.UniversalClient.HSetNX(ctx, key, field, value).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HSetNX error", slog.String("key", key), slog.String("field", field), slog.String("error", err.Error()))
		return
	}
	return ok, nil
}

// HScan 扫描哈希表中的字段和值
// cursor 是游标，从 0 开始；match 是匹配模式；count 是每次扫描的数量
// 返回新的游标和扫描到的字段值对
func (r *Client) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) (newCursor uint64, keys []string, err error) {
	keys, newCursor, err = r.UniversalClient.HScan(ctx, key, cursor, match, count).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HScan error", slog.String("key", key), slog.Uint64("cursor", cursor), slog.String("match", match), slog.Int64("count", count), slog.String("error", err.Error()))
		return
	}
	return newCursor, keys, nil
}

// HScanNoValues 扫描哈希表中的字段（不返回值）
// cursor 是游标，从 0 开始；match 是匹配模式；count 是每次扫描的数量
// 返回新的游标和扫描到的字段名列表
func (r *Client) HScanNoValues(ctx context.Context, key string, cursor uint64, match string, count int64) (newCursor uint64, keys []string, err error) {
	keys, newCursor, err = r.UniversalClient.HScanNoValues(ctx, key, cursor, match, count).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HScanNoValues error", slog.String("key", key), slog.Uint64("cursor", cursor), slog.String("match", match), slog.Int64("count", count), slog.String("error", err.Error()))
		return
	}
	return newCursor, keys, nil
}

// HVals 获取哈希表中所有字段的值
// 返回所有字段值的字符串切片
func (r *Client) HVals(ctx context.Context, key string) (values []string, err error) {
	values, err = r.UniversalClient.HVals(ctx, key).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HVals error", slog.String("key", key), slog.String("error", err.Error()))
		return
	}
	return values, nil
}

// HRandField 从哈希表中随机获取指定数量的字段名
// count 是要获取的字段数量，如果为负数则可能返回重复的字段
// 返回随机字段名的字符串切片
func (r *Client) HRandField(ctx context.Context, key string, count int) (fields []string, err error) {
	fields, err = r.UniversalClient.HRandField(ctx, key, count).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HRandField error", slog.String("key", key), slog.Int("count", count), slog.String("error", err.Error()))
		return
	}
	return fields, nil
}

// HRandFieldWithValues 从哈希表中随机获取指定数量的字段名和值
// count 是要获取的字段数量，如果为负数则可能返回重复的字段
// 返回随机字段名和值的键值对切片
func (r *Client) HRandFieldWithValues(ctx context.Context, key string, count int) (keyValues []redis.KeyValue, err error) {
	keyValues, err = r.UniversalClient.HRandFieldWithValues(ctx, key, count).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HRandFieldWithValues error", slog.String("key", key), slog.Int("count", count), slog.String("error", err.Error()))
		return
	}
	return keyValues, nil
}

// HStrLen 获取哈希表中指定字段值的字符串长度
// 返回字段值的字符串长度（字节数）
func (r *Client) HStrLen(ctx context.Context, key, field string) (length int64, err error) {
	length, err = r.UniversalClient.HStrLen(ctx, key, field).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HStrLen error", slog.String("key", key), slog.String("field", field), slog.String("error", err.Error()))
		return
	}
	return length, nil
}

// HExpire 设置哈希表中指定字段的过期时间（秒级精度）
// expiration 是过期时长
// fields 是要设置过期时间的字段列表
// 返回成功设置过期时间的字段数量
func (r *Client) HExpire(ctx context.Context, key string, expiration time.Duration, fields ...string) (counts []int64, err error) {
	counts, err = r.UniversalClient.HExpire(ctx, key, expiration, fields...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HExpire error", slog.String("key", key), slog.Duration("expiration", expiration), slog.Any("fields", fields), slog.String("error", err.Error()))
		return
	}
	return counts, nil
}

// HExpireWithArgs 设置哈希表中指定字段的过期时间（秒级精度，带参数选项）
// expiration 是过期时长
// expirationArgs 包含过期时间的配置选项
// fields 是要设置过期时间的字段列表
// 返回成功设置过期时间的字段数量
func (r *Client) HExpireWithArgs(ctx context.Context, key string, expiration time.Duration, expirationArgs redis.HExpireArgs, fields ...string) (counts []int64, err error) {
	counts, err = r.UniversalClient.HExpireWithArgs(ctx, key, expiration, expirationArgs, fields...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HExpireWithArgs error", slog.String("key", key), slog.Duration("expiration", expiration), slog.Any("fields", fields), slog.String("error", err.Error()))
		return
	}
	return counts, nil
}

// HPExpire 设置哈希表中指定字段的过期时间（毫秒级精度）
// expiration 是过期时长（毫秒）
// fields 是要设置过期时间的字段列表
// 返回成功设置过期时间的字段数量
func (r *Client) HPExpire(ctx context.Context, key string, expiration time.Duration, fields ...string) (counts []int64, err error) {
	counts, err = r.UniversalClient.HPExpire(ctx, key, expiration, fields...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HPExpire error", slog.String("key", key), slog.Duration("expiration", expiration), slog.Any("fields", fields), slog.String("error", err.Error()))
		return
	}
	return counts, nil
}

// HPExpireWithArgs 设置哈希表中指定字段的过期时间（毫秒级精度，带参数选项）
// expiration 是过期时长（毫秒）
// expirationArgs 包含过期时间的配置选项
// fields 是要设置过期时间的字段列表
// 返回成功设置过期时间的字段数量
func (r *Client) HPExpireWithArgs(ctx context.Context, key string, expiration time.Duration, expirationArgs redis.HExpireArgs, fields ...string) (counts []int64, err error) {
	counts, err = r.UniversalClient.HPExpireWithArgs(ctx, key, expiration, expirationArgs, fields...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HPExpireWithArgs error", slog.String("key", key), slog.Duration("expiration", expiration), slog.Any("fields", fields), slog.String("error", err.Error()))
		return
	}
	return counts, nil
}

// HExpireAt 设置哈希表中指定字段在指定时间点过期（秒级精度）
// tm 是过期的时间点
// fields 是要设置过期时间的字段列表
// 返回成功设置过期时间的字段数量
func (r *Client) HExpireAt(ctx context.Context, key string, tm time.Time, fields ...string) (counts []int64, err error) {
	counts, err = r.UniversalClient.HExpireAt(ctx, key, tm, fields...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HExpireAt error", slog.String("key", key), slog.Time("time", tm), slog.Any("fields", fields), slog.String("error", err.Error()))
		return
	}
	return counts, nil
}

// HExpireAtWithArgs 设置哈希表中指定字段在指定时间点过期（秒级精度，带参数选项）
// tm 是过期的时间点
// expirationArgs 包含过期时间的配置选项
// fields 是要设置过期时间的字段列表
// 返回成功设置过期时间的字段数量
func (r *Client) HExpireAtWithArgs(ctx context.Context, key string, tm time.Time, expirationArgs redis.HExpireArgs, fields ...string) (counts []int64, err error) {
	counts, err = r.UniversalClient.HExpireAtWithArgs(ctx, key, tm, expirationArgs, fields...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HExpireAtWithArgs error", slog.String("key", key), slog.Time("time", tm), slog.Any("fields", fields), slog.String("error", err.Error()))
		return
	}
	return counts, nil
}

// HPExpireAt 设置哈希表中指定字段在指定时间点过期（毫秒级精度）
// tm 是过期的时间点
// fields 是要设置过期时间的字段列表
// 返回成功设置过期时间的字段数量
func (r *Client) HPExpireAt(ctx context.Context, key string, tm time.Time, fields ...string) (counts []int64, err error) {
	counts, err = r.UniversalClient.HPExpireAt(ctx, key, tm, fields...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HPExpireAt error", slog.String("key", key), slog.Time("time", tm), slog.Any("fields", fields), slog.String("error", err.Error()))
		return
	}
	return counts, nil
}

// HPExpireAtWithArgs 设置哈希表中指定字段在指定时间点过期（毫秒级精度，带参数选项）
// tm 是过期的时间点
// expirationArgs 包含过期时间的配置选项
// fields 是要设置过期时间的字段列表
// 返回成功设置过期时间的字段数量
func (r *Client) HPExpireAtWithArgs(ctx context.Context, key string, tm time.Time, expirationArgs redis.HExpireArgs, fields ...string) (counts []int64, err error) {
	counts, err = r.UniversalClient.HPExpireAtWithArgs(ctx, key, tm, expirationArgs, fields...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HPExpireAtWithArgs error", slog.String("key", key), slog.Time("time", tm), slog.Any("fields", fields), slog.String("error", err.Error()))
		return
	}
	return counts, nil
}

// HPersist 移除哈希表中指定字段的过期时间，使其永久存在
// fields 是要移除过期时间的字段列表
// 返回成功移除过期时间的字段数量
func (r *Client) HPersist(ctx context.Context, key string, fields ...string) (counts []int64, err error) {
	counts, err = r.UniversalClient.HPersist(ctx, key, fields...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HPersist error", slog.String("key", key), slog.Any("fields", fields), slog.String("error", err.Error()))
		return
	}
	return counts, nil
}

// HExpireTime 获取哈希表中指定字段的过期时间（秒级时间戳）
// fields 是要查询的字段列表
// 返回每个字段的过期时间戳（秒），-1 表示没有设置过期时间，-2 表示字段不存在
func (r *Client) HExpireTime(ctx context.Context, key string, fields ...string) (timestamps []int64, err error) {
	timestamps, err = r.UniversalClient.HExpireTime(ctx, key, fields...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HExpireTime error", slog.String("key", key), slog.Any("fields", fields), slog.String("error", err.Error()))
		return
	}
	return timestamps, nil
}

// HPExpireTime 获取哈希表中指定字段的过期时间（毫秒级时间戳）
// fields 是要查询的字段列表
// 返回每个字段的过期时间戳（毫秒），-1 表示没有设置过期时间，-2 表示字段不存在
func (r *Client) HPExpireTime(ctx context.Context, key string, fields ...string) (timestamps []int64, err error) {
	timestamps, err = r.UniversalClient.HPExpireTime(ctx, key, fields...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HPExpireTime error", slog.String("key", key), slog.Any("fields", fields), slog.String("error", err.Error()))
		return
	}
	return timestamps, nil
}

// HTTL 获取哈希表中指定字段的剩余存活时间（秒）
// fields 是要查询的字段列表
// 返回每个字段的剩余存活时间（秒），-1 表示没有设置过期时间，-2 表示字段不存在
func (r *Client) HTTL(ctx context.Context, key string, fields ...string) (ttls []int64, err error) {
	ttls, err = r.UniversalClient.HTTL(ctx, key, fields...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HTTL error", slog.String("key", key), slog.Any("fields", fields), slog.String("error", err.Error()))
		return
	}
	return ttls, nil
}

// HPTTL 获取哈希表中指定字段的剩余存活时间（毫秒）
// fields 是要查询的字段列表
// 返回每个字段的剩余存活时间（毫秒），-1 表示没有设置过期时间，-2 表示字段不存在
func (r *Client) HPTTL(ctx context.Context, key string, fields ...string) (ttls []int64, err error) {
	ttls, err = r.UniversalClient.HPTTL(ctx, key, fields...).Result()
	if err != nil {
		slog.ErrorContext(ctx, "Redis HPTTL error", slog.String("key", key), slog.Any("fields", fields), slog.String("error", err.Error()))
		return
	}
	return ttls, nil
}
