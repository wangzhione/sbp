package trace

import (
	"log/slog"
	"testing"
)

func TestContextHandler_Handle(t *testing.T) {
	slog.Debug("你好", "你好", "你好")
	slog.Info("你好", "你好", "你好")

	ctx := Context()
	slog.WarnContext(ctx, "测试 warn", "123", "value 123")
}

func TestEnableDebug(t *testing.T) {
	b := slog.Default().Enabled(Background, slog.LevelDebug)
	t.Log(b)

	x := slog.Default().Enabled(Background, slog.LevelInfo)
	t.Log(x)

	Options.Level = slog.LevelDebug
	t.Log(slog.Default().Enabled(Background, slog.LevelDebug))
}
