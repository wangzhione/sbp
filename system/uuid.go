package system

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// UUID 生成 RFC 9562 UUIDv7
//
// 返回格式：32 位小写 hex，无横线
// 示例：018f2e4b7c2d7a4c9c8e2f3a4b5c6d7e
//
// UUIDv7 layout:
//
//	0                   47 48   51 52          63 64 65 66              127
//	|     unix_ts_ms      | ver |    rand_a      |var|      rand_b        |
//
// - unix_ts_ms: 48-bit Unix 毫秒时间戳，大端序
// - ver:        4-bit，固定 0b0111
// - rand_a:     12-bit 随机数
// - var:        2-bit，固定 0b10
// - rand_b:     62-bit 随机数
func UUID() string {
	var id [16]byte
	var od [32]byte

	// 先填满随机数，后面会覆盖 timestamp/version/variant 对应 bit
	_, _ = rand.Read(id[:])

	// UUIDv7 前 48 bit 是 Unix Epoch milliseconds，big-endian
	ms := uint64(time.Now().UnixMilli())

	id[0] = byte(ms >> 40)
	id[1] = byte(ms >> 32)
	id[2] = byte(ms >> 24)
	id[3] = byte(ms >> 16)
	id[4] = byte(ms >> 8)
	id[5] = byte(ms)

	// Version 7: byte 6 的高 4 bit 设置为 0111
	id[6] = (id[6] & 0x0f) | 0x70

	// Variant: byte 8 的高 2 bit 设置为 10
	id[8] = (id[8] & 0x3f) | 0x80

	hex.Encode(od[:], id[:])
	return string(od[:])
}

// 对于 UUID , 另一个拓展思路, 借助 MySQL UUID_SHORT() 函数返回 int128

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
