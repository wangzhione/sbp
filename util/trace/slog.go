package trace

import (
	"context"
	"log/slog"
	"os"
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
	opts := slog.HandlerOptions{
		Level: EnableLevel,
	}
	logs := slog.New(&ContextHandler{slog.NewJSONHandler(os.Stdout, &opts)})
	slog.SetDefault(logs)
}

func EnableDebug() bool {
	return slog.Default().Enabled(context.Background(), slog.LevelDebug)
}
