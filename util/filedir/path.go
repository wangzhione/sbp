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
