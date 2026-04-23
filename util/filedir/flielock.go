package filedir

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

// ErrFileLockBusy 表示锁文件已被其他进程 or goroutine 占用。
var ErrFileLockBusy = errors.New("file lock busy")

// TryFileLock 尝试为 path 获取一个跨进程锁。
// 当前实现基于 path+".lock" 锁文件，创建成功即视为拿到锁。
func TryFileLock(ctx context.Context, path string) (lock *FlieLock, err error) {
	lockpath := path + ".lock"

	err = os.MkdirAll(filepath.Dir(lockpath), 0o755)
	if err != nil {
		slog.ErrorContext(ctx, "os.MkdirAll error", "error", err, "path", path, "lockpath", lockpath)
		return nil, err
	}

	file, err := os.OpenFile(lockpath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0o664)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil, ErrFileLockBusy
		}
		slog.ErrorContext(ctx, "os.OpenFile lock error", "error", err, "path", path, "lockpath", lockpath)
		return nil, err
	}

	lock = &FlieLock{
		File: file,
		Path: lockpath,
	}

	// 留下锁持有者信息，便于排查僵尸锁文件。
	slog.InfoContext(ctx, "file lock acquired",
		"path", path, "lockpath", lockpath,
		"pid", os.Getpid(), "time", time.Now().Format(time.RFC3339Nano))

	return
}

// FileLock 持续等待直到拿到锁。
func FileLock(ctx context.Context, path string, sleep ...time.Duration) (lock *FlieLock, err error) {
	wait := 10 * time.Millisecond
	if len(sleep) > 0 && sleep[0] > 0 {
		wait = sleep[0]
	}

	for {
		lock, err = TryFileLock(ctx, path)
		if err == nil {
			return
		}
		if !errors.Is(err, ErrFileLockBusy) {
			return
		}

		time.Sleep(wait)
	}
}

// WithFileLock 获取锁并执行 fn，结束后释放锁。
func WithFileLock(ctx context.Context, path string, fn func() error, sleep ...time.Duration) error {
	lock, err := FileLock(ctx, path, sleep...)
	if err != nil {
		return err
	}

	err = fn()
	unlockerr := lock.Unlock(ctx)
	if err != nil {
		return err
	}
	return unlockerr
}

// FlieLock 表示一个已持有的锁文件句柄。
type FlieLock struct {
	*os.File
	Path string
}

func (our *FlieLock) Unlock(ctx context.Context) error {
	if our == nil {
		return nil
	}

	if our.File != nil {
		err := our.File.Close()
		if err != nil {
			slog.ErrorContext(ctx, "lock file close error", "error", err, "path", our.Path)
			return err
		}
		our.File = nil
	}

	if our.Path == "" {
		return nil
	}

	err := os.Remove(our.Path)
	if err != nil && !os.IsNotExist(err) {
		slog.ErrorContext(ctx, "os.Remove lock file error", "error", err, "path", our.Path)
		return err
	}

	our.Path = ""
	return nil
}
