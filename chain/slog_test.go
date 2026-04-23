package chain

import (
	"log/slog"
	"os"
	"testing"

	"github.com/wangzhione/sbp/system"
)

var ctx = Context()

func TestInitSLog(t *testing.T) {
	t.Log(os.Args[0])

	InitSLog()

	slog.WarnContext(ctx, "测试 warn", "123", "value 123", "234", "value 234", "456", "value 456", "789", "value 789")
}

func TestInitSLogRotatingFile(t *testing.T) {
	t.Log(system.ExeNameSuffixExt)

	InitSLogRotatingFile()

	slog.DebugContext(ctx, "测试 debug")
	slog.InfoContext(ctx, "测试 info", "123", "value 123")
	slog.WarnContext(ctx, "测试 warn", "123", "value 123", "234", "value 234")
	slog.ErrorContext(ctx, "测试 error", "123", "value 123", "234", "value 234", "456", "value 456")
}
