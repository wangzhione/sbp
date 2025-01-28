package mysql

import (
	"context"
	"log/slog"
	"time"
)

// Hooks satisfies the sqlhook.Hooks interface
type Hooks struct{}

var KeyHooks any = "begin-time"

// Before hook will print the query with it's args and return the context with the timestamp
func (h *Hooks) Before(ctx context.Context, query string, args ...any) (context.Context, error) {
	begin := time.Now()

	slog.InfoContext(ctx, "MySQL before", "begin", begin.UnixNano(), "query", query, "args", args)

	return context.WithValue(ctx, KeyHooks, begin), nil
}

// After hook will get the timestamp registered on the Before hook and print the elapsed time
func (h *Hooks) After(ctx context.Context, query string, args ...any) (context.Context, error) {
	begin, ok := ctx.Value(KeyHooks).(time.Time)
	if ok {
		since := time.Since(begin)
		if since >= time.Second {
			slog.WarnContext(ctx, "MySQL After Warn slow", "since", time.Since(begin))
		} else {
			slog.InfoContext(ctx, "MySQL After", "since", time.Since(begin))
		}
	}
	return ctx, nil
}
