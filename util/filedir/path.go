package filedir

import (
	"context"
	"log/slog"
	"path/filepath"
)

// RelPath filepath.Rel + filepath.ToSlash
func RelPath(ctx context.Context, basepath, targpath string) (relPath string, err error) {
	relPath, err = filepath.Rel(basepath, targpath)
	if err != nil {
		slog.ErrorContext(ctx, "filepath.Rel", "error", err, "basepath", basepath, "targpath", targpath)
		return
	}
	relPath = filepath.ToSlash(relPath)
	return
}

// AbsPath 解析相对路径, 得到绝对路径. 例如 . -> os.Getwd()
func AbsPath(ctx context.Context, path string) (abspath string, err error) {
	abspath, err = filepath.Abs(path)
	if err != nil {
		slog.ErrorContext(ctx, "filepath.Abs", "error", err, "path", path, "abspath", abspath)
		return
	}
	return
}
