package filedir

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
)

// CreateDir 根据 path 创建 dir
func CreateDir(ctx context.Context, path string) error {
	// 如果文件不存在，尝试创建文件所在的目录
	dir := filepath.Dir(path)

	// 检查文件是否存在
	found, _ := Exist(ctx, path)

	// 已经存在, 直接返回
	if found {
		return nil
	}

	// 确保目录存在，如果不存在则创建; 0o777	rwxrwxrwx	全执行+读写权限
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		slog.ErrorContext(ctx, "os.MkdirAll error", "error", err, "path", path, "dir", dir)
	}
	return err
}

// MkdirAll 通过 os.MkdirAll 创建 0o777 懒人目录
func MkdirAll(ctx context.Context, paths ...string) (err error) {
	for i, path := range paths {
		erri := os.MkdirAll(path, os.ModePerm)
		if erri != nil {
			slog.ErrorContext(ctx, "os.MkdirAll error", "error", erri, "path", path, "i", i)

			// 只纪录第一次出现错误, 方便用户排查原因
			if err == nil {
				err = erri
			}
		}
	}

	return
}
