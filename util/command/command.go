// Package command provides utilities for executing system commands and batch operations.
package command

import (
	"context"
	"log/slog"
	"os/exec"
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
	Serialization bool        // true 串行化; 默认 false 内部自行决定
	Options       []*exec.Cmd // 要执行的命令列表, exec.CommandContext(ctx, ...
}

func (b *BatchOption) Start() error {
	for _, opt := range b.Options {
		// 遇到错误会停下,
		// 因为有时候在错误情况下继续执行, 行为未知的, 还不如主动出错, 等待工程师接入
		if err := opt.Start(); err != nil {
			return err
		}
	}
	return nil
}

func (b *BatchOption) Wait() error {
	for _, opt := range b.Options {
		if err := opt.Wait(); err != nil {
			return err
		}
	}
	return nil
}

func (b *BatchOption) Run(ctx context.Context) error {
	if len(b.Options) == 0 {
		return nil
	}

	startTime := time.Now()
	slog.InfoContext(ctx, "BatchOption Run Begin", slog.Time("StartTime", startTime))

	defer func() {
		slog.InfoContext(ctx, "BatchOption Run End",
			slog.Float64("Duration", time.Since(startTime).Seconds()),
			slog.Time("StartTime", startTime),
		)
	}()

	if !b.Serialization {
		err := b.Start()
		if err != nil {
			slog.InfoContext(ctx, "BatchOption b.Start() error",
				slog.Time("StartTime", startTime), slog.Any("error", err))
		}
		err = b.Wait()
		if err != nil {
			slog.InfoContext(ctx, "BatchOption b.Wait() error",
				slog.Time("StartTime", startTime), slog.Any("error", err))
		}
		return err
	}

	for _, opt := range b.Options {
		// 遇到错误会停下,
		// 因为有时候在错误情况下继续执行, 行为未知的, 还不如主动出错, 等待工程师接入
		if err := opt.Start(); err != nil {
			slog.InfoContext(ctx, "BatchOption 2 b.Start() error",
				slog.Time("StartTime", startTime), slog.Any("error", err))
			return err
		}
		if err := opt.Wait(); err != nil {
			slog.InfoContext(ctx, "BatchOption 2 b.Wait() error",
				slog.Time("StartTime", startTime), slog.Any("error", err))
			return err
		}
	}
	return nil
}
