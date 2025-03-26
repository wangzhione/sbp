package chain

import (
	"context"
	"io"
	"log/slog"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

// EnableLevel 默认开启 slog.LevelDebug, 具体业务可以 init 通过配置日志等级
var EnableLevel slog.Level = slog.LevelDebug

func InitSLog() {
	options := &slog.HandlerOptions{
		Level: EnableLevel,
	}

	var handler slog.Handler = slog.NewJSONHandler(os.Stdout, options)
	if os.Getenv("LOG_FORMAT") == "text" {
		handler = slog.NewTextHandler(os.Stdout, options)
	}

	logs := slog.New(&TraceHandler{handler})
	slog.SetDefault(logs)
}

/*
	// lumberjack 会 mkdir + open file
	logger := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    600, // 单位 MB ; 0 is 不按大小分割
		MaxBackups: 0,   // 不限制备份数量
		MaxAge:     7,   // 保留日志 7 天内的所有日志
		LocalTime:  true,
		Compress:   false, // 是否压缩旧日志文件, 默认不压缩
	}
*/

type Logger = lumberjack.Logger

// InitSLogRotatingFile 需要自行管理 logger close 操作
func InitSLogRotatingFile(args ...*Logger) {
	var logger *Logger

	switch len(args) {
	case 0:
		logger = &Logger{
			MaxSize:    600, // 单位 MB ; 0 is 不按大小分割
			MaxBackups: 0,   // 不限制备份数量
			MaxAge:     7,   // 保留日志 7 天内的所有日志
			LocalTime:  true,
			Compress:   false, // 是否压缩旧日志文件, 默认不压缩
		}
	case 1:
		logger = args[0]
	case 2:
		logger = args[0]
	default:
		panic("len(args) > 1")
	}

	if len(logger.Filename) == 0 {
		logger.Filename = DefaultRotatingFile
	}

	// lumberjack 会 mkdir + open file

	options := &slog.HandlerOptions{
		Level: EnableLevel,
	}

	var multiWriter io.Writer
	if len(args) == 2 {
		// 有的项目, 喜欢在普通 {project}.log 日志基础上, 构建 {project}.error.log
		multiWriter = io.MultiWriter(os.Stdout, logger, args[1])
	} else {
		multiWriter = io.MultiWriter(os.Stdout, logger)
	}

	var handler slog.Handler = slog.NewJSONHandler(multiWriter, options)
	if os.Getenv("LOG_FORMAT") == "text" {
		handler = slog.NewTextHandler(multiWriter, options)
	}

	logs := slog.New(&TraceHandler{handler})
	slog.SetDefault(logs)
}

// LogStartEnd Wrapper function to log start and end times, and measure duration
func LogStartEnd(ctx context.Context, name string, fn func(context.Context)) {
	start := time.Now()
	slog.InfoContext(ctx, "["+name+"] - Start", "time", start.Format("2006-01-02 15:04:05.000000"))

	// Execute the wrapped function with context
	fn(ctx)

	end := time.Now()
	elapsed := end.Sub(start)
	slog.InfoContext(ctx, "["+name+"] - End", "elapsed", elapsed.Seconds(), "time", end.Format("2006-01-02 15:04:05.000000"))
}
