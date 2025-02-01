package chain

import (
	"context"
	"log/slog"
	"os"
	"time"
)

type ContextHandler struct {
	slog.Handler
}

func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	// context 需要在首次出现地方 注入 WithTraceID(&ctx) trace id
	traceID := GetTraceID(ctx)
	if len(traceID) > 0 {
		r.AddAttrs(slog.String(Key, traceID))
	}

	return h.Handler.Handle(ctx, r)
}

// EnableLevel 默认开启 slog.LevelDebug, 具体业务可以 init 通过配置日志等级
var EnableLevel slog.Level = slog.LevelDebug

func init() {
	Options := slog.HandlerOptions{
		AddSource: true,
		Level:     EnableLevel,
	}
	logs := slog.New(&ContextHandler{slog.NewTextHandler(os.Stdout, &Options)})
	slog.SetDefault(logs)
}

// LogStartEnd Wrapper function to log start and end times, and measure duration
func LogStartEnd(ctx context.Context, name string, fn func(context.Context)) {
	start := time.Now()
	slog.InfoContext(ctx, "["+name+"] - Start", "time", start.Format("2006-01-02 15:04:05.000000"))

	// Execute the wrapped function with context
	fn(ctx)

	end := time.Now()
	elapsed := end.Sub(start)
	slog.InfoContext(ctx, "["+name+"] - End", "elapsed", elapsed, "time", end.Format("2006-01-02 15:04:05.000000"))
}
