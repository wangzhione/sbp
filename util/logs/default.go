package logs

import (
	"context"
	"fmt"
	"log"
	"os"

	"sbpkg/util/idx"
)

var _ Logger = (*localLogger)(nil)

// SetDefaultLogger sets the default logger.
// This is not concurrency safe, which means it should only be called during init.
func SetDefaultLogger(s Logger) {
	if s != nil {
		defaultLogger = s
	}
}

var defaultLogger Logger = &localLogger{
	logger: log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile|log.Lmicroseconds),
}

type localLogger struct {
	logger *log.Logger
}

func (l *localLogger) logf(format string, a ...any) {
	l.logger.Output(3, fmt.Sprintf(format, a...))
}

func (l *localLogger) Fatal(ctx context.Context, format string, a ...any) {
	l.logf(idx.GetTraceID(ctx)+" [Fatal] "+format, a...)
}

func (l *localLogger) Error(ctx context.Context, format string, a ...any) {
	l.logf(idx.GetTraceID(ctx)+" [Error] "+format, a...)
}

func (l *localLogger) Warn(ctx context.Context, format string, a ...any) {
	l.logf(idx.GetTraceID(ctx)+" [Warn] "+format, a...)
}

func (l *localLogger) Notice(ctx context.Context, format string, a ...any) {
	l.logf(idx.GetTraceID(ctx)+" [Notice] "+format, a...)
}

func (l *localLogger) Info(ctx context.Context, format string, a ...any) {
	l.logf(idx.GetTraceID(ctx)+" [Info] "+format, a...)
}

func (l *localLogger) Debug(ctx context.Context, format string, a ...any) {
	l.logf(idx.GetTraceID(ctx)+" [Debug] "+format, a...)
}

func (l *localLogger) Trace(ctx context.Context, format string, a ...any) {
	l.logf(idx.GetTraceID(ctx)+" [Trace] "+format, a...)
}
