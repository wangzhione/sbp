package safego

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/wangzhione/sbp/util/chain"
)

func Run(ctx context.Context, fn func() error) (err error) {
	defer func() {
		if cover := recover(); cover != nil {
			err = fmt.Errorf("panic recovered: %#v", cover)

			// 遇到启动不起来, 异常退出, 打印堆栈方便排除问题
			slog.ErrorContext(ctx, "Run panic error",
				slog.Any("error", cover),
				slog.String("stack", string(debug.Stack())), // 记录详细的堆栈信息
			)
		}
	}()

	return fn()
}

func Go(ctx context.Context, fn func(), keys ...any) {
	begin := time.Now()
	ctx = chain.CopyTrace(ctx, keys)
	slog.InfoContext(ctx, "Go Run Begin", slog.Time("begin", begin))
	go func() {
		defer func() {
			if cover := recover(); cover != nil {
				// 遇到启动不起来, 异常退出, 打印堆栈方便排除问题
				slog.ErrorContext(ctx, "Go panic error",
					slog.Any("error", cover),
					slog.String("stack", string(debug.Stack())), // 记录详细的堆栈信息
				)
			}

			end := time.Now()
			slog.InfoContext(ctx, "Go Run End",
				slog.Time("begin", begin),
				slog.Time("end", end),
				slog.Duration("cost", end.Sub(begin)))
		}()

		fn()
	}()
}

func ID() (goroutineid string) {
	// goroutine 123 [running]:
	// uintptr int64 10 进制最长 20 位; 9 + 1 + 20 + 1 = 31 位最长
	var buf [32]byte

	// If all is false, Stack formats the stack trace for the calling goroutine."
	n := runtime.Stack(buf[:], false)
	// 把 buf[:n] 里的内容按空白字符（空格、换行、制表符）拆分成多个字段（token）

	L := 10 // len("goroutine ") = 10
	for R := L + 1; R < n; R++ {
		if buf[R] == ' ' {
			// 包含 left 位置的元素 不包含 right 位置的元素
			// 切出来的是一个长度为 right - left 的切片
			goroutineid = string(buf[L:R])
			return
		}
	}

	return
}
