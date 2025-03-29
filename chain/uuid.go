package chain

import (
	"crypto/rand"
	"io"
)

const hextable = "0123456789abcdef"

// UUID 的全称是 Universally Unique Identifier, "通用唯一标识符" or "全球唯一标识符"
func UUID() string {
	// // A UUID is a 128 bit (16 byte) Universal Unique IDentifier as defined in RFC9562.
	var uuid [16]byte
	_, err := io.ReadFull(rand.Reader, uuid[:]) // random function
	if err != nil {
		// 特定低版本 linux 内核 rand 会出错, 一旦出错, Go 运行时默认退出 ...
		// fatal("crypto/rand: failed to read random data (see https://go.dev/issue/66821): " + err.Error())
		// panic("unreachable") // To be sure.

		// never case, 兜底 用默认 uuid 串返回
		return "00000000000000000000000000000000"
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10

	var data [32]byte
	// "00000000-0000-0000-0000-000000000000" {8}-{4}-{4}-{4}-{12}
	for i, v := range uuid {
		data[i<<1] = hextable[v>>4]       // high 4 bit
		data[(i<<1)+1] = hextable[v&0x0f] //  low 4 bit
	}

	return string(data[:])
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
