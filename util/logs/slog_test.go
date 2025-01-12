package logs

import (
	"log/slog"
	"sbpkg/util/idx"
	"testing"
)

func TestContextHandler_Handle(t *testing.T) {
	slog.Info("你好", "你好", "你好")

	ctx := idx.Context()
	slog.WarnContext(ctx, "测试 warn", "123", "value 123")
}
