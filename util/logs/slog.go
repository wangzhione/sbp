package logs

import (
	"context"
	"log/slog"
	"os"
	"sbpkg/util/idx"
)

type ContextHandler struct {
	slog.Handler
}

func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	traceID := idx.GetTraceID(ctx)
	if len(traceID) > 0 {
		r.AddAttrs(slog.String(idx.TraceIDKey, traceID))
	}

	return h.Handler.Handle(ctx, r)
}

func init() {
	opts := slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logs := slog.New(&ContextHandler{slog.NewJSONHandler(os.Stdout, &opts)})
	slog.SetDefault(logs)
}
