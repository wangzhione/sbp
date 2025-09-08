package chain

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
)

// EnableLevel 默认开启 slog.LevelDebug, 具体业务可以 init 通过配置日志等级
var EnableLevel slog.Level = slog.LevelDebug

func InitSLog() {
	slog.SetDefault(slog.New(&TraceHandler{
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: EnableLevel,
		}),
	}))
}

func InitSlogRotatingFile() error {
	// 默认是 start day logger;
	// 如果需要 hour logger, Please DefaultGetFile = GetfileByHour 随后 Call InitSlogRotatingFile()
	return Startlogger()
}

type TraceHandler struct {
	slog.Handler
}

var CodeKey = "code"

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

	// go run test   : e:\github.com\wangzhione\sbp\chain\slog_test.go:26:TestInitSlogRotatingFile
	// go debug test : slog_test.go:27:TestInitSlogRotatingFile
	source := fmt.Sprintf("%s:%d:%s", filepath.Base(frame.File), frame.Line, funcName)

	r.AddAttrs(
		// context 依赖 WithContext(ctx, {trace id}) or Request(r)
		slog.String(XRquestID, GetTraceID(ctx)),

		// short code source, 和 slog.HandlerOptions::AddSource 可以共存, 推荐 设置 AddSource = false
		slog.String(CodeKey, source),
	)

	return h.Handler.Handle(ctx, r)
}
