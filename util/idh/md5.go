package idh

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

func MD5(data []byte) string {
	// 小写 16 进制
	digest := md5.Sum(data)
	return hex.EncodeToString(digest[:])
}

// MD5String string to md5 sign
func MD5String(s string) string {
	return MD5([]byte(s))
}

func MD5File(path string) (sign string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}

	d := md5.New()
	// io.Copy 类似 32K copy buffer 读取直到 读取到 EOF, 然后成功的话 err == nil 并返回
	_, err = io.Copy(d, file)
	if err != nil {
		return
	}

	sign = hex.EncodeToString(d.Sum(nil))
	return
}
