package idhash

import (
	"crypto/md5"
	"encoding/base64"
	"testing"
)

func TestBase64Encode(t *testing.T) {
	// Error MD5 Base64 example

	content := "0123456789"

	sign := MD5(content)
	if sign != "781e5e245d69b566979b86e28d23f2c7" {
		t.Fatal("MD5 fatal", content, sign)
	}

	base64sign := Base64Encode(sign)
	if base64sign != "NzgxZTVlMjQ1ZDY5YjU2Njk3OWI4NmUyOGQyM2YyYzc=" {
		t.Fatal("Base64Encode fatal", content, sign, base64sign)
	}

	t.Log("Error Success", content, "->", sign, "->", base64sign)
}

func TestMD5Base64Encode(t *testing.T) {
	// True MD5 Base64 example

	content := "0123456789"

	digest := md5.Sum([]byte(content))

	base64sign := base64.StdEncoding.EncodeToString(digest[:])
	if base64sign != "eB5eJF1ptWaXm4bijSPyxw==" {
		t.Fatal("Base64Encode fatal", content, base64sign)
	}

	t.Log("True Success", content, "->", base64sign)

	md5sign, err := Base64MD5Decode(base64sign)
	if err != nil {
		t.Fatal("Base64MD5Decode is fatal", base64sign, err)
	}
	if md5sign != "781e5e245d69b566979b86e28d23f2c7" {
		t.Fatal("MD5 fatal", content, md5sign)
	}
	t.Log("Success Base64MD5Decode", content, md5sign)
}
