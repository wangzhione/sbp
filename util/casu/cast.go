package casu

import (
	"strconv"
)

// INT int or uint numbers type
type INT interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~uintptr
}

// FormatINT int to string 默认都是 10 进制数
func FormatINT[T INT](i T) string {
	if 0 <= i && i < nSmalls {
		return small(int(i))
	}
	_, s := formatBits(nil, uint64(i), 10, i < 0, false)
	return s
}

// ParseINT string to int 默认都是 10 进制, 内部吃掉 error
func ParseINT[T INT](s string) T {
	if s == "" {
		return 0
	}

	if s[0] == '-' || s[0] == '+' {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0
		}
		return T(v)
	}

	u, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return T(u)
}

// ParseINTE string to int 默认都是 10 进制, 返回给外部需要处理的 Error
func ParseINTE[T INT](s string) (i T, err error) {
	if s == "" {
		return
	}

	if s[0] == '-' || s[0] == '+' {
		var v int64
		v, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			return
		}
		i = T(v)
		return
	}

	var u uint64
	u, err = strconv.ParseUint(s, 10, 64)
	if err != nil {
		return
	}
	i = T(u)
	return
}

// ParseBool 无异常 ParseBool 版本, 可以配合 strconv.FormatBool 互相转换
// 其中 strconv.FormatBool returns "true" or "false" according to the value of b.
func ParseBool(s string) bool {
	switch s {
	case "1", "t", "T",
		"true", "TRUE", "True",
		"truE",
		"trUe", "trUE",
		"tRue", "tRuE", "tRUe", "tRUE",
		"TruE", "TrUe", "TrUE", "TRue", "TRuE", "TRUe":
		return true
	}
	return false
}

// FormatFloat float to string, 这是个商业业务代码, 不是科学代码, 业务场景不应该出现 float32
func FormatFloat[T ~float64 | ~float32](f T) string {
	// The special precision -1 uses the smallest number of digits
	// necessary such that ParseFloat will return f exactly.
	return strconv.FormatFloat(float64(f), 'f', -1, 64)
}

// ParseFloat string to float64, 业务上不应该出现 float32, 如果需要自行 float32(casu.ParseFloat(string))
func ParseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
