package logs

import (
	"context"
	"os"
)

// Logger is a logger interface that provides logging function with levels.
type Logger interface {
	Trace(ctx context.Context, format string, a ...any)
	Debug(ctx context.Context, format string, a ...any)
	Info(ctx context.Context, format string, a ...any)
	Notice(ctx context.Context, format string, a ...any)
	Warn(ctx context.Context, format string, a ...any)
	Error(ctx context.Context, format string, a ...any)
	Fatal(ctx context.Context, format string, a ...any)
}

// Level defines the priority of a log message.
// When a logger is configured with a level, any log message with a lower
// log level (smaller by integer comparison) will not be output.
type Level uint8

// The levels of logs.
// 单纯的业务开发, 你应该只需要考虑 Debug Info Warn Error 四类业务日志
const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelNotice
	LevelWarn
	LevelError
	LevelFatal
)

var level Level

// SetLevel sets the level of logs below which logs will not be output.
// The default log level is LevelTrace.
func SetLevel(v Level) {
	if LevelTrace <= v && v <= LevelFatal {
		level = v
	}
}

// Fatal calls the default logger's Fatal method and then os.Exit(1).
func Fatal(ctx context.Context, format string, a ...interface{}) {
	defaultLogger.Fatal(ctx, format, a...)
	os.Exit(1)
}

// Error calls the default logger's Error method.
func Error(ctx context.Context, format string, a ...interface{}) {
	if level > LevelError {
		return
	}
	defaultLogger.Error(ctx, format, a...)
}

// Warn calls the default logger's Warn method.
func Warn(ctx context.Context, format string, a ...interface{}) {
	if level > LevelWarn {
		return
	}
	defaultLogger.Warn(ctx, format, a...)
}

// Notice calls the default logger's Notice method.
func Notice(ctx context.Context, format string, a ...interface{}) {
	if level > LevelNotice {
		return
	}
	defaultLogger.Notice(ctx, format, a...)
}

// Info calls the default logger's Info method.
func Info(ctx context.Context, format string, a ...interface{}) {
	if level > LevelInfo {
		return
	}
	defaultLogger.Info(ctx, format, a...)
}

// Debug calls the default logger's Debug method.
func Debug(ctx context.Context, format string, a ...interface{}) {
	if level > LevelDebug {
		return
	}
	defaultLogger.Debug(ctx, format, a...)
}

// Trace calls the default logger's Trace method.
func Trace(ctx context.Context, format string, a ...interface{}) {
	if level > LevelTrace {
		return
	}
	defaultLogger.Trace(ctx, format, a...)
}
