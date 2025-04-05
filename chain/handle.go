package chain

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
)

type TraceHandler struct {
	slog.Handler
}

// Handle add trace
// @see https://github.com/golang/go/issues/73054#event-16988835247
func (h TraceHandler) Handle(ctx context.Context, r slog.Record) error {
	frames := runtime.CallersFrames([]uintptr{r.PC})
	frame, _ := frames.Next()

	// {path}/{short package name}.{short func name} -> {short func name}
	i := len(frame.Function) - 2
	for i >= 0 && frame.Function[i] != '/' && frame.Function[i] != '.' {
		i--
	}
	funcName := frame.Function[i+1:]

	source := fmt.Sprintf("%s:%d.%s", filepath.Base(frame.File), frame.Line, funcName)

	r.AddAttrs(
		// context 依赖 WithContext(ctx, {trace id}) or Request(r)
		slog.String(XRquestID, GetTraceID(ctx)),

		// // short source, need slog.HandlerOptions::AddSource = false
		slog.String(slog.SourceKey, source),
	)

	return h.Handler.Handle(ctx, r)
}
