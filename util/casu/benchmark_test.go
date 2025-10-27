package casu

import (
	"strconv"
	"testing"
)

// BenchmarkFormatINT 测试整数转字符串的性能
func BenchmarkFormatINT(b *testing.B) {
	testCases := []struct {
		name string
		val  int64
	}{
		{"Small", 42},
		{"Medium", 12345},
		{"Large", 9223372036854775807},
		{"Negative", -12345},
		{"Zero", 0},
		{"MaxInt64", 9223372036854775807},
		{"MinInt64", -9223372036854775808},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = FormatINT(tc.val)
			}
		})
	}
}

// BenchmarkStrconvFormatInt 测试标准库strconv.FormatInt的性能
func BenchmarkStrconvFormatInt(b *testing.B) {
	testCases := []struct {
		name string
		val  int64
	}{
		{"Small", 42},
		{"Medium", 12345},
		{"Large", 9223372036854775807},
		{"Negative", -12345},
		{"Zero", 0},
		{"MaxInt64", 9223372036854775807},
		{"MinInt64", -9223372036854775808},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = strconv.FormatInt(tc.val, 10)
			}
		})
	}
}

// BenchmarkStrconvItoa 测试标准库strconv.Itoa的性能
func BenchmarkStrconvItoa(b *testing.B) {
	testCases := []int{42, 12345, -12345, 0, 2147483647, -2147483648}

	for _, val := range testCases {
		b.Run(strconv.Itoa(val), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = strconv.Itoa(val)
			}
		})
	}
}

// BenchmarkParseINT 测试字符串转整数的性能（无错误处理）
func BenchmarkParseINT(b *testing.B) {
	testCases := []struct {
		name string
		val  string
	}{
		{"Small", "42"},
		{"Medium", "12345"},
		{"Large", "9223372036854775807"},
		{"Negative", "-12345"},
		{"Zero", "0"},
		{"MaxInt64", "9223372036854775807"},
		{"MinInt64", "-9223372036854775808"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = ParseINT[int64](tc.val)
			}
		})
	}
}

// BenchmarkParseINTE 测试字符串转整数的性能（带错误处理）
func BenchmarkParseINTE(b *testing.B) {
	testCases := []struct {
		name string
		val  string
	}{
		{"Small", "42"},
		{"Medium", "12345"},
		{"Large", "9223372036854775807"},
		{"Negative", "-12345"},
		{"Zero", "0"},
		{"MaxInt64", "9223372036854775807"},
		{"MinInt64", "-9223372036854775808"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = ParseINTE[int64](tc.val)
			}
		})
	}
}

// BenchmarkStrconvParseInt 测试标准库strconv.ParseInt的性能
func BenchmarkStrconvParseInt(b *testing.B) {
	testCases := []struct {
		name string
		val  string
	}{
		{"Small", "42"},
		{"Medium", "12345"},
		{"Large", "9223372036854775807"},
		{"Negative", "-12345"},
		{"Zero", "0"},
		{"MaxInt64", "9223372036854775807"},
		{"MinInt64", "-9223372036854775808"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = strconv.ParseInt(tc.val, 10, 64)
			}
		})
	}
}

// BenchmarkStrconvAtoi 测试标准库strconv.Atoi的性能
func BenchmarkStrconvAtoi(b *testing.B) {
	testCases := []string{"42", "12345", "-12345", "0", "2147483647", "-2147483648"}

	for _, val := range testCases {
		b.Run(val, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = strconv.Atoi(val)
			}
		})
	}
}

// BenchmarkSmallNumbers 专门测试小数字(0-99)的性能差异
func BenchmarkSmallNumbers(b *testing.B) {
	// 测试0-99的小数字
	for i := 0; i < 100; i++ {
		b.Run("FormatINT_"+strconv.Itoa(i), func(b *testing.B) {
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				_ = FormatINT(i)
			}
		})
	}
}

// BenchmarkSmallNumbersStrconv 测试标准库小数字性能
func BenchmarkSmallNumbersStrconv(b *testing.B) {
	// 测试0-99的小数字
	for i := 0; i < 100; i++ {
		b.Run("StrconvItoa_"+strconv.Itoa(i), func(b *testing.B) {
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				_ = strconv.Itoa(i)
			}
		})
	}
}

// BenchmarkUint64 测试uint64类型的性能
func BenchmarkUint64(b *testing.B) {
	testCases := []uint64{0, 42, 12345, 18446744073709551615}

	for _, val := range testCases {
		b.Run(strconv.FormatUint(val, 10), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = FormatINT(val)
			}
		})
	}
}

// BenchmarkStrconvFormatUint 测试标准库FormatUint性能
func BenchmarkStrconvFormatUint(b *testing.B) {
	testCases := []uint64{0, 42, 12345, 18446744073709551615}

	for _, val := range testCases {
		b.Run(strconv.FormatUint(val, 10), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = strconv.FormatUint(val, 10)
			}
		})
	}
}

// BenchmarkMemoryAllocation 测试内存分配情况
func BenchmarkMemoryAllocation(b *testing.B) {
	b.Run("FormatINT", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = FormatINT(12345)
		}
	})

	b.Run("StrconvFormatInt", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = strconv.FormatInt(12345, 10)
		}
	})
}

// BenchmarkConcurrent 测试并发性能
func BenchmarkConcurrent(b *testing.B) {
	b.Run("FormatINT_Concurrent", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = FormatINT(12345)
			}
		})
	})

	b.Run("StrconvFormatInt_Concurrent", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = strconv.FormatInt(12345, 10)
			}
		})
	})
}
