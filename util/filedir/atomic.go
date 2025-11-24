package filedir

import (
	"os"
	"path/filepath"
)

// FSyncWriteFile 原子性地将 data 写入指定路径的文件中。
func FSyncWriteFile(path string, data []byte, perm ...os.FileMode) error {
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	// 自动创建目录（可选）
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	// 创建临时文件：
	//   pattern = base + ".*.temp"
	//   其中 "*" 会被自动替换成随机字符串，形成如下格式的文件名：
	//
	//       <base>.<random>.temp
	//
	//   例如：
	//       base = "config.json"
	//       生成的文件可能是：
	//       "config.json.123456789.temp"
	//
	// 临时文件会被创建在目标目录 dir 中，确保 rename 原子性。
	// 临时文件将使用系统默认权限（受 umask 影响）
	temp, err := os.CreateTemp(dir, base+".*.temp")
	if err != nil {
		return err
	}
	name := temp.Name()

	defer os.Remove(name)

	// 若调用方显式传入 perm（例如 0644）则设置权限
	if len(perm) > 0 && perm[0] != 0 {
		if err := temp.Chmod(perm[0]); err != nil {
			temp.Close()
			return err
		}
	}

	// 写入循环，避免短写
	written := 0
	for written < len(data) {
		n, err := temp.Write(data[written:])
		if err != nil {
			temp.Close()
			return err
		}
		written += n
	}

	// fsync 文件内容
	if err := temp.Sync(); err != nil {
		temp.Close()
		return err
	}

	// Windows rename 之前需要先关闭句柄
	if err := temp.Close(); err != nil {
		return err
	}

	// 原子替换
	if err := os.Rename(name, path); err != nil {
		return err
	}

	// fsync 目录项（尽力而为）
	if dirFile, err := os.Open(dir); err == nil {
		_ = dirFile.Sync()
		dirFile.Close()
	}

	return nil
}
