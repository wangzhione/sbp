package command

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"
	"time"

	"github.com/wangzhione/sbp/util/jsou"
)

// Run exec.CommandContext 帮助操作
func Run(ctx context.Context, bin string, args ...string) (err error) {
	err = exec.CommandContext(ctx, bin, args...).Run()
	if err != nil {
		slog.ErrorContext(ctx, "exec.CommandContext error", "error", err, "bin", bin, "args", args)
	}
	return
}

// Options 定义命令的执行参数, Run 等同于 Start + Wait
type Options struct {
	Bin     string        // [必填] 命令名
	Args    []string      // [可选] 命令参数
	Work    string        // [可选] 工作目录
	Timeout time.Duration // [可选] 超时时间

	Stdlog bool   // [可选] 是否打印日志 stdout + stderr
	Stdout string // [可选] 和 ShowLog 绑定, true 会纪录 stdout 内容
	Stderr string // [可选] 和 ShowLog 绑定, true 会纪录 stderr 内容

	NotCommandLog bool // [可选] 默认 false, 打印最终的 command log
}

func (opts *Options) Command() string {
	var cmd bytes.Buffer

	cmd.WriteString(opts.Bin)
	for _, arg := range opts.Args {
		cmd.WriteString(" ")

		// 如果参数里有空格或特殊符号，加引号
		if strings.ContainsAny(arg, " \t\n\"'\\") {
			cmd.WriteString("\"")
			cmd.WriteString(arg)
			cmd.WriteString("\"")
		} else {
			cmd.WriteString(arg)
		}
	}

	return cmd.String()
}

func (opts *Options) Run(ctx context.Context) (err error) {
	start := time.Now()

	defer func() {
		duration := time.Since(start)

		// 工程代码就是如此, 健壮性代码远远多于普通正式运行的逻辑
		if cover := recover(); cover != nil {
			err = fmt.Errorf("Options panic error recovered: %#v", cover)
			slog.ErrorContext(ctx, "Options Run panic error",
				slog.Any("error", cover),
				slog.String("stack", string(debug.Stack())),
			)
		}

		if err != nil {
			slog.ErrorContext(ctx, "Options Run exec.CommandContext error",
				slog.Any("error", err),
				slog.String("Options", jsou.String(opts)),
			)
		}

		if !opts.NotCommandLog {
			slog.InfoContext(ctx, "Options Run End",
				slog.Duration("Duration", duration),
				slog.Time("StartTime", start),
			)
		}
	}()

	if !opts.NotCommandLog {
		slog.InfoContext(ctx, "Options Run Start",
			slog.String("Command", opts.Command()),
			slog.Time("StartTime", start),
		)
	}

	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, opts.Bin, opts.Args...)

	if opts.Work != "" {
		cmd.Dir = opts.Work
	}

	if opts.Stdlog {
		var stdout, stderr bytes.Buffer
		cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)

		err = cmd.Run()
		opts.Stdout = stdout.String()
		opts.Stderr = stderr.String()
	} else {
		err = cmd.Run()
	}
	return
}
