package groupgo

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/wangzhione/sbp/util/httpip"
)

type DownloadTask struct {
	URL  string // 待下载 url
	Path string // 目标地址

	Headers map[string]string // http download head
}

func (task *DownloadTask) Check() error {
	if task.URL == "" || task.Path == "" {
		return errors.New("error: download task param empty")
	}
	return nil
}

// DownloadTask 表示单个下载任务, 这里这种模式类似对象函数
type DownloadGroup struct {
	Task []DownloadTask

	Headers       map[string]string // http download head
	MaxConcurrent int               // group go max 并发
}

func (down *DownloadGroup) Check(ctx context.Context) error { // Check and init
	if down.MaxConcurrent <= 0 {
		down.MaxConcurrent = 16 // 二八芳龄, 很多魔法数字, 猛拍脑门
	}

	// 参数 check
	for i := range down.Task {
		if err := down.Task[i].Check(); err != nil {
			slog.ErrorContext(ctx, "down.Task[i].Check() error", "error", err, "task", down.Task[i])
			return err
		}

		if down.Headers != nil && down.Task[i].Headers == nil {
			down.Task[i].Headers = down.Headers
		}
	}

	return nil
}

// Download 批量并发下载，使用 groupgo 管理 goroutine，限制最大并发数, 默认同步下载
func (down *DownloadGroup) Download(ctx context.Context) (err error) {
	if len(down.Task) == 0 {
		return
	}

	if err = down.Check(ctx); err != nil {
		return
	}

	// 尝试阻塞下载模式
	group := NewGroup(ctx, down.MaxConcurrent)

	for _, task := range down.Task {
		group.Go(func(ctx context.Context) error {
			start := time.Now()
			slog.InfoContext(ctx, "Download task start",
				"uri", task.URL,
				"path", task.Path,
			)

			taskerr := httpip.DownloadIfNotExists(ctx, task.URL, task.Path, task.Headers)

			duration := time.Since(start)

			// 下载是比较重操作, 增加多一点业务打点, 这是个 demo 库, 至少介绍 groupgo 用法
			slog.InfoContext(ctx, "Download task end",
				"uri", task.URL,
				"path", task.Path,
				"duration", duration.Seconds(),
				"reason", taskerr,
			)

			return taskerr
		})
	}

	return group.Wait()
}
