package times

import (
	"testing"
	"time"
)

func TestShanghaiTimeLocation(t *testing.T) {
	t.Log("ShanghaiLoction", time.Now().In(ShanghaiLoction))

	t.Log("ShanghaiTimeString", NowString())
}
