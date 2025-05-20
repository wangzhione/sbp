package safego

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"runtime/debug"
)

// goroutine 是有成本, 并且存在泄露或崩溃的风险!
// 你需要足够能力去驾驭它. 所以对于 下面模板代码应该形成肌肉记忆 or 了然于胸
/*
go func() {
	defer func() {
		if cover := recover(); cover != nil {
			// 遇到启动不起来, 异常退出, 打印堆栈方便排除问题
			slog.ErrorContext(ctx, "Go panic error",
				slog.Any("error", cover),
				slog.String("stack", string(debug.Stack())), // 记录详细的堆栈信息
			)
		}
	}()

	...
}()
*/

func Run(ctx context.Context, fn func(context.Context) error) (err error) {
	defer func() {
		if cover := recover(); cover != nil {
			err = fmt.Errorf("panic: safego.Run %#v", cover)

			// 遇到启动不起来, 异常退出, 打印堆栈方便排除问题
			slog.ErrorContext(ctx, "Run panic error",
				slog.Any("error", err),
				slog.String("stack", string(debug.Stack())), // 记录详细的堆栈信息
			)
		}
	}()

	return fn(ctx)
}

func Cover(ctx context.Context) {
	if cover := recover(); cover != nil {
		// 遇到启动不起来, 异常退出, 打印堆栈方便排除问题
		slog.ErrorContext(ctx, "recover go panic error",
			slog.Any("error", cover),
			slog.String("stack", string(debug.Stack())), // 记录详细的堆栈信息
		)
	}
}

func Go(c context.Context, fn func(context.Context)) {
	go func(ctx context.Context) {
		defer Cover(ctx)

		fn(ctx)
	}(c)
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
