package filedir

import (
	"os"
	"path/filepath"
)

// CreateDir 根据 path 创建 dir
func CreateDir(path string) error {
	// 获取文件所在的目录路径
	dir := filepath.Dir(path)

	return os.MkdirAll(dir, 0o755)
}
