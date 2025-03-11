package idhash

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

//
// base64 库语法糖帮助库, 如需要获取详细 error 推荐 base64 DecodeString
//

// Base64Encode encodes a string to standard Base64
func Base64Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

// Base64Decode decodes a Base64-encoded string.
// 如需要 error 请:
//
//	decodedBytes, err := base64.StdEncoding.DecodeString(input)
//	if err != nil { ... return }
//	[]data : decodedBytes -> ...
func Base64Decode(output string) string {
	decodedBytes, _ := base64.StdEncoding.DecodeString(output)
	return string(decodedBytes)
}

// Base64EncodeURL encodes a string to URL-safe Base64
func Base64EncodeURL(input string) string {
	return base64.URLEncoding.EncodeToString([]byte(input))
}

// Base64DecodeURL decodes a URL-safe Base64-encoded string
func Base64DecodeURL(output string) string {
	decodedBytes, _ := base64.URLEncoding.DecodeString(output)
	return string(decodedBytes)
}

// Base64MD5 input -> md5 [16]byte -> base64 string
func Base64MD5(input string) string {
	digest := md5.Sum([]byte(input))
	return base64.StdEncoding.EncodeToString(digest[:])
}

// ErrBase64MD5Size base64 解码后 size error
var ErrBase64MD5Size = errors.New("error: base64 md5 size")

// Base64MD5Decode base64 string -> md5 [16]byte -> input
func Base64MD5Decode(output string) (string, error) {
	digest, err := base64.StdEncoding.DecodeString(output)
	if err != nil {
		return "", err
	}
	if len(digest) != md5.Size {
		return "", ErrBase64MD5Size
	}
	return hex.EncodeToString(digest), nil
}
