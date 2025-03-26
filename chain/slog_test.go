package chain

import (
	"log/slog"
	"os"
	"testing"
)

var ctx = Context()

func TestContextHandler_Handle(t *testing.T) {
	slog.Debug("你好", "你好", "你好")
	slog.Info("你好", "你好", "你好")

	slog.WarnContext(ctx, "测试 warn", "123", "value 123")

	t.Log(EnableLevel)
}

func TestInitSLog(t *testing.T) {
	t.Log(os.Args[0])

	InitSLog()

	slog.WarnContext(ctx, "测试 warn", "123", "value 123")
}

func TestInitSLogRotatingFile(t *testing.T) {
	t.Log(os.Args[0])

	path := "logs/log.log"
	// var path string

	InitSLogRotatingFile(&Logger{Filename: path})

	slog.WarnContext(ctx, "测试 warn", "123", "value 123")

	slog.Debug("你好", "你好", "你好")
	slog.Info("你好", "你好", "你好")
}
