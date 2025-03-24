package command

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Run exec.CommandContext 帮助操作
func Run(ctx context.Context, bin string, args ...string) (err error) {
	err = exec.CommandContext(ctx, bin, args...).Run()
	if err != nil {
		slog.ErrorContext(ctx, "exec.CommandContext error", "error", err, "bin", bin, "args", args)
	}
	return
}

type BatchOption struct {
	C       context.Context // context, 自行 chain.CopyTrace or context.WithTimeout
	Options []*Option       // 要执行的命令列表
}

func (b *BatchOption) Start() error {
	for _, opt := range b.Options {
		// 遇到错误会停下,
		// 因为有时候在错误情况下继续执行, 行为未知的, 还不如主动出错, 等待工程师接入
		if err := opt.Start(b.C); err != nil {
			opt.cmd = nil
			return err
		}
	}

	return nil
}

func (b *BatchOption) Wait() error {
	for _, opt := range b.Options {
		if opt.cmd == nil {
			break // b.Start 启动
		}
		if err := opt.Wait(); err != nil {
			return err
		}
	}
	return nil
}

func (b *BatchOption) Run() error {
	if err := b.Start(); err != nil {
		return err
	}
	return b.Wait()
}

// Option 定义命令的执行参数, Run 等同于 Start + Wait
type Option struct {
	// Path is the path of the command to run.
	//
	// This is the only field that must be set to a non-zero
	// value. If Path is relative, it is evaluated relative
	// to Dir.
	Path string

	// [optional] Args holds command line arguments, including the command as Args[0].
	// If the Args field is empty or nil, Run uses {Path}.
	//
	// In typical use, both Path and Args are set by calling Command.
	Args []string

	// [optional] Env specifies the environment of the process.
	// Each entry is of the form "key=value".
	// If Env is nil, the new process uses the current process's
	// environment.
	// If Env contains duplicate environment keys, only the last
	// value in the slice for each duplicate key is used.
	// As a special case on Windows, SYSTEMROOT is always added if
	// missing and not explicitly set to the empty string.
	//
	// See also the Dir field, which may set PWD in the environment.
	Env []string

	// [optional] Dir specifies the working directory of the command.
	// If Dir is the empty string, Run runs the command in the
	// calling process's current directory.
	//
	// On Unix systems, the value of Dir also determines the
	// child process's PWD environment variable if not otherwise
	// specified. A Unix process represents its working directory
	// not by name but as an implicit reference to a node in the
	// file tree. So, if the child process obtains its working
	// directory by calling a function such as C's getcwd, which
	// computes the canonical name by walking up the file tree, it
	// will not recover the original value of Dir if that value
	// was an alias involving symbolic links. However, if the
	// child process calls Go's [os.Getwd] or GNU C's
	// get_current_dir_name, and the value of PWD is an alias for
	// the current directory, those functions will return the
	// value of PWD, which matches the value of Dir.
	Dir string

	Stdlog         bool // [optional] 是否需要 standard 标准 日志 stdout + stderr 内容
	Stdout, Stderr bytes.Buffer

	StartTime time.Time       // Start 开始时间
	c         context.Context // [必填] Start 方法传入的 context.Context 参数
	cmd       *exec.Cmd       // [必填] cmd 构建对象
}

func (o *Option) Command() string {
	var cmd bytes.Buffer

	cmd.WriteString(o.Path)
	for _, arg := range o.Args {
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

// Start 在当前 Go 进程中创建一个子进程，并执行指定的命令，而不等待其执行完毕（即非阻塞）。
// 它只是启动这个命令，并让它运行在后台，由 cmd.Wait() 来负责等待它结束。
//
// Timeout time.Duration // [optional] 超时时间
//
//	if o.Timeout > 0 {
//	    var cancel context.CancelFunc
//	    ctx, cancel = context.WithTimeout(ctx, o.Timeout)
//	    defer cancel()
//	}
func (o *Option) Start(ctx context.Context) (err error) {
	o.StartTime = time.Now()
	o.c = ctx

	slog.InfoContext(ctx, "Option Start Begin",
		slog.Time("StartTime", o.StartTime),
		slog.String("Command", o.Command()),
		slog.Any("Env", o.Env),
		slog.String("Dir", o.Dir),
	)

	o.cmd = exec.CommandContext(ctx, o.Path, o.Args...)
	o.cmd.Env = o.Env
	o.cmd.Dir = o.Dir
	if o.Stdlog {
		o.cmd.Stdout = io.MultiWriter(os.Stdout, &o.Stdout)
		o.cmd.Stderr = io.MultiWriter(os.Stderr, &o.Stderr)
	}

	err = o.cmd.Start()
	if err != nil {
		slog.ErrorContext(ctx, "Option Start o.cmd.Start()",
			slog.Any("error", err),
			slog.Time("StartTime", o.StartTime),
		)
	}

	return
}

func (o *Option) Wait() (err error) {
	err = o.cmd.Wait()
	if err != nil {
		// 工程代码就是如此, 健壮性代码远远多于普通正式运行的逻辑
		slog.ErrorContext(o.c, "Option Run exec.CommandContext error",
			slog.Any("error", err),
			slog.Time("StartTime", o.StartTime),
			slog.String("Command", o.Command()),
			slog.Any("Env", o.Env),
			slog.String("Dir", o.Dir),
		)
	}

	slog.InfoContext(o.c, "Option Run End",
		slog.Duration("Duration", time.Since(o.StartTime)),
		slog.Time("StartTime", o.StartTime),
	)

	return
}

// Run starts the specified command and waits for it to complete.
//
// The returned error is nil if the command runs, has no problems
// copying stdin, stdout, and stderr, and exits with a zero exit
// status.
//
// If the command starts but does not complete successfully, the error is of
// type [*ExitError]. Other error types may be returned for other situations.
//
// If the calling goroutine has locked the operating system thread
// with [runtime.LockOSThread] and modified any inheritable OS-level
// thread state (for example, Linux or Plan 9 name spaces), the new
// process will inherit the caller's thread state.
func (o *Option) Run(ctx context.Context) error {
	if err := o.Start(ctx); err != nil {
		return err
	}
	return o.Wait()
}
