package casu

import (
	"strconv"
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

func FuzzStringToIntE(f *testing.F) {
	testCases := []string{
		"123", "-456", "0", "18446744073709551615", "-9223372036854775808", "not_a_number", "", "++123", "--456", " 42", "- 42",
	}

	for _, tc := range testCases {
		f.Add(tc) // 添加初始测试样例
	}

	f.Fuzz(func(t *testing.T, s string) {
		var result int64
		var err error

		result, err = StringToIntE[int64](s)
		if err != nil {
			t.Logf("Expected error for input %q: %v", s, err)
		}

		// 验证转换结果是否符合 strconv 的标准行为
		if expected, convErr := strconv.ParseInt(s, 10, 64); convErr == nil {
			if result != expected {
				t.Errorf("Mismatch: expected %d, got %d for input %q", expected, result, s)
			}
		} else {
			if err == nil {
				t.Errorf("Expected error but got none for input %q", s)
			}
		}
	})
}
