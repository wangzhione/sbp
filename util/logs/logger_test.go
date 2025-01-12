package logs

import (
	"testing"

	"sbpkg/util/idx"
)

func TestSetLevel(t *testing.T) {
	var ctx = idx.Context()

	Debug(ctx, "debug")
	Info(ctx, "%s 重要信息", "8964")
	Warn(ctx, "警告信息 %d", 1010)
	Error(ctx, "Error")
}
