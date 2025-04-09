package chain

import (
	"context"
	"log/slog"
	"os"
)

// Fatal slog.ErrorContext + os.Exit(-1)
func Fatal(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, args...)
	os.Exit(-1)
}
