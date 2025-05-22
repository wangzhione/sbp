package chain

import (
	"testing"
	"time"
)

func Test_daylogger_sevenday(t *testing.T) {
	timeStr := "20250514"
	// 解析时间
	x, err := time.Parse("2006010215", timeStr)
	t.Log(x, err)

	x, err = time.Parse("20060102", timeStr)
	t.Log(x, err)

	timeStr = "2025051423"
	x, err = time.Parse("20060102", timeStr)
	t.Log(x, err)

	timeStr = "1234"
	if len(timeStr) > 8 {
		timeStr = timeStr[:8]
	}
	t.Log(timeStr[:])
}
