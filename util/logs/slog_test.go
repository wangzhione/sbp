package logs

import (
	"log/slog"
	"testing"

	"sbp/util/idh"
)

func TestContextHandler_Handle(t *testing.T) {
	slog.Info("你好", "你好", "你好")

	ctx := idh.Context()
	slog.WarnContext(ctx, "测试 warn", "123", "value 123")
}
