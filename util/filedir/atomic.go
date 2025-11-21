package filedir

import (
	"os"
	"path/filepath"
)

// AtomicSyncWriteFile 以 “写临时文件 + 原子 rename” 的方式安全写文件：
// 1. 确保同目录下写入临时文件
// 2. 写完并 fsync 文件内容
// 3. 关闭文件句柄（兼容 Windows）
// 4. 原子重命名替换目标文件
// 5. 尝试 fsync 目录，增强崩溃安全
//
// 参数 perm：
//   - 如果 perm != 0，则使用 perm 作为新文件权限
//   - 如果 perm == 0 且旧文件已存在，则继承旧文件的权限
//   - 否则使用系统默认权限（受 umask 影响）
func AtomicSyncWriteFile(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	// 确保目录存在（如果你不想自动建目录，可以去掉这一段）
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	// 在目标目录下创建临时文件
	temp, err := os.CreateTemp(dir, base+".*.tmp")
	if err != nil {
		return err
	}
	name := temp.Name()

	// 确保最后清理临时文件（无论成功或失败）
	// - 如果 rename 成功，这个路径已不存在，Remove 会静默失败
	// - 如果中途出错，临时文件会被删除，避免目录堆垃圾
	defer os.Remove(name)

	// 决定最终权限：
	// 1. 调用方显式传了 perm -> 用它
	// 2. 否则如果原文件存在 -> 继承原文件权限
	// 3. 否则保留 CreateTemp 的默认权限（受 umask 影响）
	finalPerm := perm
	if finalPerm == 0 {
		if fi, err := os.Stat(path); err == nil {
			finalPerm = fi.Mode().Perm()
		}
	}
	if finalPerm != 0 {
		if err := temp.Chmod(finalPerm); err != nil {
			temp.Close()
			return err
		}
	}

	// 循环写入，防止短写
	written := 0
	for written < len(data) {
		n, err := temp.Write(data[written:])
		if err != nil {
			temp.Close()
			return err
		}
		written += n
	}

	// fsync 文件内容，确保数据落盘
	if err := temp.Sync(); err != nil {
		temp.Close()
		return err
	}

	// 关闭文件句柄（特别是 Windows 上，未关闭无法 rename/删除）
	if err := temp.Close(); err != nil {
		return err
	}

	// 原子 rename 替换目标文件：
	// - 要么看到旧文件，要么看到新文件
	// - 不会出现半截内容
	if err := os.Rename(name, path); err != nil {
		return err
	}

	// 尽力 fsync 目录，保证目录项更新也落盘
	// 某些平台/文件系统可能不支持对目录 fsync，这种情况下的错误可以忽略
	if dirFile, err := os.Open(dir); err == nil {
		_ = dirFile.Sync() // 忽略错误（例如不支持目录 fsync）
		dirFile.Close()
	}

	return nil
}
