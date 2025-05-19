package groupgo

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/wangzhione/sbp/util/filedir"
	"github.com/wangzhione/sbp/util/httpip"
)

// groupgodemo.go :: DownloadTask 适合自行 Ctrl + C -> Ctrl + V 用于实际并发下载业务中.
// 需要注意是 context 生命周期, 因为有些 http 服务 call 结束适合, context 会被取消
// 还要注意, 下载业务并不是多线程, 多进程安全的. 你应该审视你的业务为什么下载还需要多进程安全进行浪费?
//

type DownloadTask struct {
	URL     string            // 待下载 url
	Path    string            // 目标地址
	Headers map[string]string // http download head

	Log   bool // 是否打开打点日志, 默认不打开
	Force bool // 是否强制更新下载, 默认不强制下载
}

func (task *DownloadTask) Check() error {
	if task.URL == "" || task.Path == "" {
		return errors.New("error: download task param empty")
	}
	return nil
}

// DownloadTask 表示单个下载任务, 这里这种模式类似对象函数
type DownloadGroup struct {
	Task          []DownloadTask
	MaxConcurrent int // group go max 并发
}

func (down *DownloadGroup) Check(ctx context.Context) error {
	if down.MaxConcurrent <= 0 {
		down.MaxConcurrent = 16 // 拍脑门, 魔法数字
	}

	// 参数 check
	for i := range down.Task {
		if err := down.Task[i].Check(); err != nil {
			slog.ErrorContext(ctx, "down.Task[i].Check() error", "error", err, "task", down.Task[i])
			return err
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
		group.Go(func(ctx context.Context) (taskerr error) {
			// 如果目标文件已存在，直接跳过
			found, err := filedir.Exist(ctx, task.Path)
			if err != nil {
				return
			}
			if found && !task.Force {
				// 文件存在, 并且不需要强制下载 直接返回
				return
			}

			if task.Log {
				start := time.Now()
				slog.InfoContext(ctx, "Download task start",
					"uri", task.URL,
					"path", task.Path,
					"force", task.Force,
				)

				defer func() {
					duration := time.Since(start)
					// 这是个 demo 库, 至少介绍 groupgo 用法
					slog.InfoContext(ctx, "Download task end",
						"uri", task.URL,
						"path", task.Path,
						"duration", duration.Seconds(),
						"force", task.Force,
						"reason", taskerr,
					)
				}()
			}

			return httpip.Download(ctx, task.URL, task.Path, task.Headers)
		})
	}

	return group.Wait()
}
