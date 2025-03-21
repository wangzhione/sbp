package idhash

import (
	"encoding/hex"

	"github.com/google/uuid"
)

func UUID() string {
	// 依赖 google uuid random uuid v4 算法构建
	v4, err := uuid.NewRandom()
	if err != nil {
		// never case, 兜底 用默认 uuid 串返回
		return "00000000000000000000000000000000"
	}

	var dst [32]byte
	// "00000000-0000-0000-0000-000000000000" {8}-{4}-{4}-{4}-{12}
	hex.Encode(dst[:], v4[:4])
	hex.Encode(dst[8:12], v4[4:6])
	hex.Encode(dst[12:16], v4[6:8])
	hex.Encode(dst[16:20], v4[8:10])
	hex.Encode(dst[20:], v4[10:])
	return string(dst[:])
}

// 对于 UUID , 另一个思路, 借助 MySQL UUID_SHORT() 函数返回 int128 , 可惜是 Go 官方还没有支持相关 bigint 类型

/*
	ulonglong uuid_value;

	void uuid_short_init() {
		uuid_value = ((ulonglong)(server_id & 255) << 56) + ((ulonglong) server_startup_time_in_seconds << 24);
	}

	longlong uuid_short() {
		mysql_mutex_lock(&LOCK_uuid_generator);
		ulonglong val = uuid_value++;
		mysql_mutex_unlock(&LOCK_uuid_generator);

		return (longlong) val;
	}
*/
