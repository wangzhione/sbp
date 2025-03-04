package casu

import (
	"testing"
)

func TestStringToFloat(t *testing.T) {
	f := 3.1415926
	s := FloatToString(f)
	t.Log(f, s)

	if f != StringToFloat(s) {
		t.Error("FloatToString, StringToFloat error")
	}
	if s != FloatToString(f) {
		t.Error("FloatToString, StringToFloat error")
	}
}

func TestStringToInt(t *testing.T) {
	var i8 int8 = 1
	var u16 uint16 = 2
	i := 3
	var u64 uint64 = 5

	if i8 != StringToInt[int8](IntToString(i8)) {
		t.Error("StringToInt, IntToString error i8", i8)
	}

	if u16 != StringToInt[uint16](IntToString(u16)) {
		t.Error("StringToInt, IntToString error u16", u16)
	}

	if i != StringToInt[int](IntToString(i)) {
		t.Error("StringToInt, IntToString error i", i)
	}

	if u64 != StringToInt[uint64](IntToString(u64)) {
		t.Error("StringToInt, IntToString error u64", u64)
	}
}
