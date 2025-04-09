package high

import (
	"strconv"
	"testing"
)

func FuzzAtoi(f *testing.F) {
	// Add some seed inputs for coverage
	seeds := []string{
		"123", "-123", "+456", "0", "0000123",
		"999999999", "18446744073709551615", // max uint64
		"9223372036854775807",  // max int64
		"-9223372036854775808", // min int64
		"abc", "", "++1", "--2", "12a34",
	}

	for _, seed := range seeds {
		f.Add([]byte(seed))
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		got, err := Atoi(data)
		str := string(data)
		// 这个投巧, 属于蒸馏 标准库 strconv.Atoi 结果
		expect, convErr := strconv.Atoi(str)

		if convErr != nil {
			if err == nil && len(data) != 0 {
				t.Errorf("Atoi(%q) = %d, expected error but got none", data, got)
			}
			return
		}

		if err != nil {
			t.Errorf("Atoi(%q) = error %v, expected %d", data, err, expect)
		}
		if got != expect {
			t.Errorf("Atoi(%q) = %d, want %d", data, got, expect)
		}
	})
}

func FuzzParseInt(f *testing.F) {
	seeds := []string{
		"0", "-1", "1", "+1", "12345", "-12345",
		"9223372036854775807", "-9223372036854775808",
		"18446744073709551615", "abc", "", "+-123", "99999999999999999999999999",
	}

	for _, seed := range seeds {
		f.Add([]byte(seed), 64)
	}

	f.Fuzz(func(t *testing.T, data []byte, bitSize int) {
		if bitSize < 0 || bitSize > 64 {
			return
		}
		got, err := ParseInt(data, bitSize)
		expect, convErr := strconv.ParseInt(string(data), 10, bitSize)

		if convErr != nil {
			if err == nil && len(data) != 0 {
				t.Errorf("ParseInt(%q, %d) = %d, expected error but got none", data, bitSize, got)
			}
			return
		}

		if err != nil {
			t.Errorf("ParseInt(%q, %d) = error %v, expected %d", data, bitSize, err, expect)
		}
		if got != expect {
			t.Errorf("ParseInt(%q, %d) = %d, want %d", data, bitSize, got, expect)
		}
	})
}

func FuzzParseUint(f *testing.F) {
	seeds := []string{
		"0", "123", "18446744073709551615", // max uint64
		"99999999999999999999999999999999",
		"-123", "+123", "abc", "", "12a34",
	}

	for _, seed := range seeds {
		f.Add([]byte(seed), 64)
	}

	f.Fuzz(func(t *testing.T, data []byte, bitSize int) {
		if bitSize < 0 || bitSize > 64 {
			return
		}
		got, err := ParseUint(data, bitSize)
		expect, convErr := strconv.ParseUint(string(data), 10, bitSize)

		if convErr != nil {
			if err == nil && len(data) != 0 {
				t.Errorf("ParseUint(%q, %d) = %d, expected error but got none", data, bitSize, got)
			}
			return
		}

		if err != nil {
			t.Errorf("ParseUint(%q, %d) = error %v, expected %d", data, bitSize, err, expect)
		}
		if got != expect {
			t.Errorf("ParseUint(%q, %d) = %d, want %d", data, bitSize, got, expect)
		}
	})
}
