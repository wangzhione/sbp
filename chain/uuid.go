package chain

import (
	"crypto/rand"
	"encoding/hex"
)

// UUID 的全称是 Universally Unique Identifier, "通用唯一标识符" or "全球唯一标识符"
func UUID() string {
	var id [16]byte         // A UUID is a 128 bit (16 byte) Universal Unique IDentifier as defined in RFC9562.
	_, _ = rand.Read(id[:]) // random function ; 细节 @see go/src/crypto/rand/rand.go

	id[6] = (id[6] & 0x0f) | 0x40 // Version 4
	id[8] = (id[8] & 0x3f) | 0x80 // Variant is 10

	var od [32]byte // "00000000 0000 0000 0000 000000000000" {8}{4}{4}{4}{12}
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
