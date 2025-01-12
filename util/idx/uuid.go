package idx

import (
	"encoding/hex"
	"log/slog"
	"runtime/debug"

	"github.com/google/uuid"
)

func UUID() (id string) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("sbpkg id UUID() panic", "recover", r, "stack", debug.Stack())
			// 填充默认的 id
			id = "00000000000000000000000000000000"
			return
		}
	}()

	var v4 = uuid.New()
	// 依赖 google uuid random uuid v4 算法构建, 其算法内部存在 panic, 默认会屏蔽吃掉 panic
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
