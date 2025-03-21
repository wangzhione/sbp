package idhash

import "testing"

func TestMD5(t *testing.T) {
	sign := MD5("123456") // e10adc3949ba59abbe56e057f20f883e
	t.Log(sign)

	if sign != "e10adc3949ba59abbe56e057f20f883e" {
		t.Fatal("MD5 123456 error")
	}
}
