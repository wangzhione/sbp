package casu

import (
	"strconv"

	"github.com/wangzhione/sbp/util/casu/high"
)

// INT int or uint numbers type
type INT interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~uintptr
}

// FormatINT int to string 默认都是 10 进制数, fast quickly
func FormatINT[T INT](i T) string {
	if 0 <= i && i < nSmalls {
		return small(int(i))
	}
	return format10(uint64(i), i < 0)
}

// ParseINT string to int 默认都是 10 进制, 内部吃掉 error
func ParseINT[T INT](s string) (i T) {
	i, _ = ParseINTE[T](s)
	return
}

// ParseINTE string to int 默认都是 10 进制, 返回给外部需要处理的 Error
func ParseINTE[T INT](s string) (i T, err error) {
	if s == "" {
		return
	}

	if s[0] == '-' || s[0] == '+' {
		var v int64
		v, err = high.ParseInt([]byte(s), 64)
		if err != nil {
			return
		}
		i = T(v)
		return
	}

	var u uint64
	u, err = high.ParseUint([]byte(s), 64)
	if err != nil {
		return
	}
	i = T(u)

	// fix : 18446744073709551615
	if i < 0 || uint64(i) != u {
		return 0, strconv.ErrRange
	}
	return
}

// ParseBool 无异常 ParseBool 版本, 可以配合 strconv.FormatBool 互相转换
// 其中 strconv.FormatBool returns "true" or "false" according to the value of b.
func ParseBool(s string) bool {
	switch s {
	case "1", "T", "t",
		"TRUE", "true", "True",
		"truE", "trUe", "trUE",
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

var a = strconv.Itoa
