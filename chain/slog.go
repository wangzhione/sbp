package chain

import (
	"context"
	"log/slog"
	"os"
	"time"
)

var LOG_FORMAT string = "json" // "json" or "text"

func isText() bool {
	return LOG_FORMAT == "text" || os.Getenv("LOG_FORMAT") == "text"
}

// EnableLevel 默认开启 slog.LevelDebug, 具体业务可以 init 通过配置日志等级
var EnableLevel slog.Level = slog.LevelDebug

func InitSLog() {
	options := &slog.HandlerOptions{
		Level: EnableLevel,
	}

	var handler slog.Handler
	if isText() {
		handler = slog.NewTextHandler(os.Stdout, options)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, options)
	}

	logs := slog.New(&TraceHandler{handler})
	slog.SetDefault(logs)
}

func InitSlogRotatingFile() error {
	err := os.MkdirAll(LogsDir, os.ModePerm)
	if err != nil {
		println("os.MkdirAll error", LogsDir)
		return err
	}
	return starthourlylogger()
}

// LogStartEnd Wrapper function to log start and end times, and measure duration
func LogStartEnd(ctx context.Context, name string, fn func(context.Context) error) (err error) {
	start := time.Now()
	slog.InfoContext(ctx, "["+name+"] - Start", "time", start.Format("2006-01-02 15:04:05.000000"))

	// Execute the wrapped function with context
	err = fn(ctx)

	end := time.Now()
	elapsed := end.Sub(start)
	slog.InfoContext(ctx, "["+name+"] - End", "elapsed", elapsed.Seconds(), "time", end.Format("2006-01-02 15:04:05.000000"))
	return
}
