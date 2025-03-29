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

// Handle add trace id @see https://github.com/golang/go/issues/73054#event-16988835247
func (h TraceHandler) Handle(ctx context.Context, r slog.Record) error {
	fs := runtime.CallersFrames([]uintptr{r.PC})
	f, _ := fs.Next()

	// {path}/{short package name}.{short func name} -> {short func name}
	i := len(f.Function) - 2
	for ; i >= 0 && f.Function[i] != '/' && f.Function[i] != '.'; i-- {
	}
	funcName := f.Function[i+1:]

	// add short source
	source := fmt.Sprintf("%s:%d:%s", filepath.Base(f.File), f.Line, funcName)
	r.AddAttrs(slog.String(slog.SourceKey, source))
	// context 依赖 WithContext(ctx, id) or Request(r)
	r.AddAttrs(slog.String(XRquestID, GetTraceID(ctx)))

	return h.Handler.Handle(ctx, r)
}
