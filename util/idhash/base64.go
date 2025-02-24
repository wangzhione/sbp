package idhash

import (
	"encoding/base64"
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
func Base64Decode(input string) string {
	decodedBytes, _ := base64.StdEncoding.DecodeString(input)
	return string(decodedBytes)
}

// Base64EncodeURL encodes a string to URL-safe Base64
func Base64EncodeURL(input string) string {
	return base64.URLEncoding.EncodeToString([]byte(input))
}

// Base64DecodeURL decodes a URL-safe Base64-encoded string
func Base64DecodeURL(input string) string {
	decodedBytes, _ := base64.URLEncoding.DecodeString(input)
	return string(decodedBytes)
}
