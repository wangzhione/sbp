package chain

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	ExePath          = os.Args[0]                          // ExePath 获取可执行文件路径(相对路径 or 绝对路径)
	ExeDir           = filepath.Dir(ExePath)               // ExeDir 获取可执行文件所在目录, 结尾不带 '/'
	ExeName          = filepath.Base(ExePath)              // ExeName 获取不带路径的可执行文件名
	ExeExt           = filepath.Ext(ExeName)               // ExeExt 获取可执行文件名的扩展名
	ExeNameSuffixExt = strings.TrimSuffix(ExeName, ExeExt) // ExeNameSuffixExt 获取可执行文件名, 不包含扩展名
)

var Hostname = func() string {
	// 获取容器的 hostname（通常是容器的短 ID）
	hostname, err := os.Hostname()
	if err == nil {
		return hostname
	}

	return UUID()
}()

// Exist 判断路径（文件或目录）是否存在
func Exist(path string) (exists bool, err error) {
	_, err = os.Stat(path)
	if err == nil {
		return true, nil // 路径存在（无论是文件还是目录）
	}

	if os.IsNotExist(err) {
		return false, nil // 路径不存在
	}
	return false, err // 其他错误（如权限问题）, 但对当前用户而言是不存在
}

// LogStartEnd Wrapper function to log start and end times, and measure duration
func LogStartEnd(ctx context.Context, name string, fn func(context.Context) error) (err error) {
	start := time.Now()
	slog.InfoContext(ctx, "["+name+"] - Start", "time", start.Format("2006-01-02 15:04:05.000000"))

	// Execute the wrapped function with context
	err = fn(ctx)

	end := time.Now()
	elapsed := end.Sub(start)
	slog.InfoContext(ctx, "["+name+"] - End", "elapsed", elapsed.Seconds(), "time", end.Format("2006-01-02 15:04:05.000000"))
	return
}

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
