package chain

import (
	"context"
	"io"
	"log/slog"
	"os"
	"testing"
	"time"

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

func TestTraceHandlerHandleZeroPC(t *testing.T) {
	handler := TraceHandler{
		Handler: slog.NewTextHandler(io.Discard, nil),
	}
	record := slog.NewRecord(time.Time{}, slog.LevelInfo, "test", 0)

	if err := handler.Handle(context.Background(), record); err != nil {
		t.Fatal(err)
	}
}
