package casu

import (
	"strconv"
	"testing"
)

func TestStringToFloat(t *testing.T) {
	f := 3.141592600
	s := FormatFloat(f)
	t.Log(f, s)

	if f != ParseFloat(s) {
		t.Error("FloatToString, StringToFloat error")
	}
	if s != FormatFloat(f) {
		t.Error("FloatToString, StringToFloat error")
	}
}

func TestStringToInt(t *testing.T) {
	var i8 int8 = 1
	var u16 uint16 = 2
	i := 3
	var u64 uint64 = 5

	if i8 != ParseINT[int8](FormatINT(i8)) {
		t.Error("StringToInt, IntToString error i8", i8)
	}

	if u16 != ParseINT[uint16](FormatINT(u16)) {
		t.Error("StringToInt, IntToString error u16", u16)
	}

	if i != ParseINT[int](FormatINT(i)) {
		t.Error("StringToInt, IntToString error i", i)
	}

	if u64 != ParseINT[uint64](FormatINT(u64)) {
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

		result, err = ParseINTE[int64](s)
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

// 假设 FormatINT 支持 int64 和 uint64（你可加其他类型）
func FuzzFormatINT(f *testing.F) {
	// 添加初始种子值
	f.Add(int64(0))
	f.Add(int64(42))
	f.Add(int64(-99999))
	f.Add(int64(1844674407370955161))

	f.Fuzz(func(t *testing.T, i int64) {
		// 调用你要测试的函数
		s := FormatINT(i)

		// 使用 strconv 验证一致性（默认十进制）
		want := strconv.FormatInt(i, 10)
		if s != want {
			t.Errorf("FormatINT(%d) = %q; want %q", i, s, want)
		}
	})
}
