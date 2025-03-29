package chain

import (
	"testing"
)

func Test_hourlylogger_sevenday(t *testing.T) {
	DefaultCleanTime = 0

	our := &hourlylogger{} // our 类似跨函数闭包
	if err := our.rotate(); err != nil {
		t.Log("Error", err)
	} else {
		t.Log("Success")
	}
}
