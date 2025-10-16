// Package idhash provides utility functions for generating MD5 hashes of strings, byte slices, and files.
package idhash

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// MD5 string to md5 sign
func MD5[T ~[]byte | ~string](data T) string {
	// 小写 16 进制
	digest := md5.Sum([]byte(data))
	return hex.EncodeToString(digest[:])
}

func MD5File(path string) (sign string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	d := md5.New()
	// io.Copy 类似 32K copy buffer 读取直到 读取到 EOF, 然后成功的话 err == nil 并返回
	_, err = io.Copy(d, file)
	if err != nil {
		return
	}

	sign = hex.EncodeToString(d.Sum(nil))
	return
}
