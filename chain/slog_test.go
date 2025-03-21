package chain

import (
	"log/slog"
	"os"
	"testing"
)

func TestContextHandler_Handle(t *testing.T) {
	slog.Debug("你好", "你好", "你好")
	slog.Info("你好", "你好", "你好")

	ctx := Context()
	slog.WarnContext(ctx, "测试 warn", "123", "value 123")
}

func TestEnableDebug(t *testing.T) {
	t.Log(EnableLevel)
}

func TestInitSLogRotatingFile(t *testing.T) {
	t.Log(os.Args[0])

	path := "logs/log.log"
	// var path string

	InitSLogRotatingFile(&Logger{Filename: path})

	ctx := Context()
	slog.WarnContext(ctx, "测试 warn", "123", "value 123")

	slog.Debug("你好", "你好", "你好")
	slog.Info("你好", "你好", "你好")
}
