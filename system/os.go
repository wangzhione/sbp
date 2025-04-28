package system

import (
	"context"
	"log/slog"
	"runtime"
	"runtime/debug"
	"time"
)

// Linux 默认是服务部署的最终服务器, 方便利用 system.Linux 默认做一些特殊处理逻辑
const Linux bool = runtime.GOOS == "linux"

/*
 runtime.GOOS 是 Go 语言中的一个常量，用于获取当前操作系统的名称。它的枚举值包括但不限于：

 windows
 linux
 darwin (macOS)
 freebsd
 openbsd
 netbsd
 android
 ios
 js (用于 Go 编译为 JavaScript)
 plan9
 solaris
*/

// BeginTime 系统启动时间
var BeginTime = time.Now()

// End 主要用于 main 函数中 defer End(context) 操作, 纪录程序结束的行为
func End(ctx context.Context) {
	if cover := recover(); cover != nil {
		// 遇到启动不起来, 异常退出, 打印堆栈方便排除问题
		slog.ErrorContext(ctx, "main init panic error",
			slog.Any("error", cover),
			slog.Time("SystemBeginTime", BeginTime),
			slog.String("GOOS", runtime.GOOS),
			slog.String("BuildVersion", BuildVersion),
			slog.String("GitVersion", GitVersion),
			slog.String("stack", string(debug.Stack())), // 记录详细的堆栈信息
		)
	}

	end := time.Now()
	slog.InfoContext(ctx, "main init end ...",
		slog.Time("SystemBeginTime", BeginTime),
		slog.Float64("elapsed_hours", end.Sub(BeginTime).Hours()),
		slog.Time("EndTime", end),
		slog.String("GOOS", runtime.GOOS),
		slog.String("BuildVersion", BuildVersion),
		slog.String("GitVersion", GitVersion),
	)
}
