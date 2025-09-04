package chain

import (
	"crypto/rand"
	"encoding/hex"
)

// UUID 的全称是 Universally Unique Identifier, "通用唯一标识符" or "全球唯一标识符"
func UUID() string {
	// // A UUID is a 128 bit (16 byte) Universal Unique IDentifier as defined in RFC9562.
	var uuid [16]byte
	_, _ = rand.Read(uuid[:]) // random function ; 细节 @see go/src/crypto/rand/rand.go

	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10

	var out [32]byte
	// "00000000-0000-0000-0000-000000000000" {8}-{4}-{4}-{4}-{12}
	hex.Encode(out[:], uuid[:])
	return string(out[:])
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
