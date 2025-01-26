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

var Options = slog.HandlerOptions{
	AddSource: true,
	Level:     slog.LevelInfo, // 默认开启 slog.LevelInfo, 具体业务可以 init 通过配置日志等级
}

func init() {
	logs := slog.New(&ContextHandler{slog.NewTextHandler(os.Stdout, &Options)})
	slog.SetDefault(logs)
}
