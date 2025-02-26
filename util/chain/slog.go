package chain

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/wangzhione/sbp/util/filedir"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ContextHandler struct {
	slog.Handler
}

func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	// context 需要在首次出现地方 注入 WithTraceID(&ctx) trace id
	traceID := GetTraceID(ctx)
	if len(traceID) > 0 {
		r.AddAttrs(slog.String(Key, traceID))
	}

	return h.Handler.Handle(ctx, r)
}

// EnableLevel 默认开启 slog.LevelDebug, 具体业务可以 init 通过配置日志等级
var EnableLevel slog.Level = slog.LevelDebug

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

func InitRotatingFileSLog(logger ...Logger) {
	switch len(logger) {
	case 0:
		logger = append(logger, Logger{
			// 自适应 path, 默认 {ExeDir}/logs/{ExeName}.log
			Filename:   filepath.Join(filedir.ExeDir, "logs", filedir.ExeName+".log"),
			MaxSize:    600, // 单位 MB ; 0 is 不按大小分割
			MaxBackups: 0,   // 不限制备份数量
			MaxAge:     7,   // 保留日志 7 天内的所有日志
			LocalTime:  true,
			Compress:   false, // 是否压缩旧日志文件, 默认不压缩
		})
	case 1:
		if len(logger[0].Filename) == 0 {
			logger[0].Filename = filepath.Join(filedir.ExeDir, "logs", filedir.ExeName+".log")
		}
	default:
	}

	for i := 0; i < len(logger); i++ {
		println(
			fmt.Sprintf(
				`{"index":%d, "time":"%s", "level":"DEBUG", "msg":"InitRotatingFileSLog", "path":"%s"}`,
				i+1,
				time.Now(),
				logger[i].Filename,
			),
		)
	}

	// lumberjack 会 mkdir + open file

	Options := &slog.HandlerOptions{
		AddSource: true,
		Level:     EnableLevel,
	}

	var multiWriter io.Writer
	if len(logger) > 1 {
		// 有的项目, 喜欢在普通 {project}.log 日志基础上, 构建 {project}.error.log
		multiWriter = io.MultiWriter(os.Stdout, &logger[0], &logger[1])
	} else {
		multiWriter = io.MultiWriter(os.Stdout, &logger[0])
	}

	var handler slog.Handler = slog.NewJSONHandler(multiWriter, Options)
	if os.Getenv("LOG_FORMAT") == "text" {
		handler = slog.NewTextHandler(multiWriter, Options)
	}

	logs := slog.New(&ContextHandler{handler})
	slog.SetDefault(logs)
}

func InitFileSLog(path string) {
	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		slog.Error("os.OpenFile panic error", "path", path, "error", err)
	}

	Options := &slog.HandlerOptions{
		AddSource: true,
		Level:     EnableLevel,
	}

	// 创建 MultiWriter，让日志同时写入文件和终端
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	var handler slog.Handler = slog.NewJSONHandler(multiWriter, Options)
	if os.Getenv("LOG_FORMAT") == "text" {
		handler = slog.NewTextHandler(multiWriter, Options)
	}

	logs := slog.New(&ContextHandler{handler})
	slog.SetDefault(logs)
}

func InitSLog() {
	Options := &slog.HandlerOptions{
		AddSource: true, // 启用日志源文件定位
		Level:     EnableLevel,
	}

	var handler slog.Handler = slog.NewJSONHandler(os.Stdout, Options)
	if os.Getenv("LOG_FORMAT") == "text" {
		handler = slog.NewTextHandler(os.Stdout, Options)
	}

	logs := slog.New(&ContextHandler{handler})
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
	slog.InfoContext(ctx, "["+name+"] - End", "elapsed", elapsed, "time", end.Format("2006-01-02 15:04:05.000000"))
}
