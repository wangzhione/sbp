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

// OpenFile 打开文件
func OpenFile(path string) (file *os.File, err error) {
	// 检查文件是否存在
	_, err = os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}

		// 如果文件不存在，尝试创建文件所在的目录
		dir := filepath.Dir(path)

		// 确保目录存在，如果不存在则创建
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return
		}
	}

	// os.OpenFile 内部有 runtime.SetFinalizer(f.file, (*file).close), 对象释放时候会 GC 1 close -> GC 2 free
	return os.OpenFile(path, os.O_RDWR, 0o644)
}

// Exist 粗略检查文件是否存在
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
