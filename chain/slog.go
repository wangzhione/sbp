package chain

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"
)

// EnableLevel 默认开启 slog.LevelDebug, 具体业务可以 init 通过配置日志等级
var EnableLevel slog.Level = slog.LevelDebug

func InitSLog() {
	options := &slog.HandlerOptions{
		Level: EnableLevel,
	}

	var handler slog.Handler
	if EnableText() {
		handler = slog.NewTextHandler(os.Stdout, options)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, options)
	}

	logs := slog.New(&TraceHandler{handler})
	slog.SetDefault(logs)
}

func InitSlogRotatingFile() error {
	return Startdaylogger()
}

// EnableText 日志给专业人士看的, 当前行业显学, 还是以 json 格式为主流.
// 设计上越独裁, 使用方越自由, 要么简单用, 要么不用
func EnableText() bool {
	return strings.EqualFold(os.Getenv("LOG_FORMAT"), "text")
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
