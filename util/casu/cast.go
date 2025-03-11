package casu

import (
	"log/slog"
	"strconv"
)

// IUNumber int or uint numbers type
type IUNumber interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~uintptr
}

// IntToString int to string 默认都是 10 进制数
func IntToString[T IUNumber](i T) string {
	if 0 <= i && i < nSmalls {
		return small(int(i))
	}
	_, s := formatBits(nil, uint64(i), 10, i < 0, false)
	return s
}

// StringToInt string to int 默认都是 10 进制, 内部吃掉 error 业务上会打印日志
func StringToInt[T IUNumber](s string) T {
	if s == "" {
		return 0
	}

	if s[0] == '-' || s[0] == '+' {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			slog.Error("ParseInt to int error", "value", s, "reason", err)
			return 0
		}
		return T(v)
	}

	u, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		slog.Error("ParseUint to int error", "value", s, "reason", err)
		return 0
	}
	return T(u)
}

// StringToIntE string to int 默认都是 10 进制, 需要处理 Error
func StringToIntE[T IUNumber](s string) (i T, err error) {
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

// StringToBool 无异常 ParseBool 版本, 可以配合 FormatBool 互相转换
// 其中 FormatBool returns "true" or "false" according to the value of b.
func StringToBool(s string) bool {
	if len(s) > 4 || len(s) < 1 {
		return false
	}

	switch s[0] {
	case '1', 't', 'T':
		return true
	}
	return false
}

// FloatToString float to string, 这是个商业业务代码, 不是科学代码, 业务场景不应该出现 float32
func FloatToString[T ~float64 | ~float32](f T) string {
	return strconv.FormatFloat(float64(f), 'f', -1, 64)
}

func StringToFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		slog.Error("StringToFloat to float64 error", "value", s, "reason", err)
	}
	return f
}
