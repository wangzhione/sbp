package filedir

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

// os.RemoveAll 删除文件 or 文件夹
// os.ReadFile data []byte -> ReadString text string = string(data)

// Exist 判断路径（文件或目录）是否存在
func Exist(ctx context.Context, path string) (exists bool, err error) {
	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // 路径不存在
		}
		slog.ErrorContext(ctx, "os.Stat error", "error", err, "path", path)
		return false, err // 其他错误（如权限问题）, 但对当前用户而言是不存在
	}
	return true, nil // 路径存在（无论是文件还是目录）
}

func IsExist(ctx context.Context, filename string) bool {
	exists, err := Exist(ctx, filename)
	if err == nil {
		// 这部分结果是 逻辑正确的, return true 就是存在, return false 表示不存在
		return exists
	}

	// err != nil
	// 这时候其实是不知道. 内部默认当 false 不存在, 让其业务自行再试试
	// not err IsNotExist 和 not err IsExist 业务上不是互为逆函数

	return false
}

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

// CopyWriter src file copy io.Writer
func CopyWriter(ctx context.Context, src string, writer io.Writer) error {
	// 打开源文件
	source, err := os.Open(src)
	if err != nil {
		slog.ErrorContext(ctx, "os.Open error", "error", err, "src", src)
		return err
	}
	defer source.Close()

	// 使用 io.Copy 进行高效复制
	_, err = io.Copy(writer, source)
	if err != nil {
		slog.ErrorContext(ctx, "io.Copy error", "error", err, "src", src)
	}
	return err
}

// CopyBodyFile resp.Body write to destination file; low api 会主动 body.Close
func CopyBodyFile(ctx context.Context, body io.ReadCloser, dst string) error {
	defer body.Close()

	// 创建目标文件
	destination, err := os.Create(dst)
	if err != nil {
		slog.ErrorContext(ctx, "os.Create error", "error", err, "dst", dst)
		return err
	}
	defer destination.Close()

	// 使用 io.Copy 进行高效复制
	_, err = io.Copy(destination, body)
	if err != nil {
		slog.ErrorContext(ctx, "io.Copy error", "error", err, "dst", dst)
	}
	return err
}

// CopyFile 复制文件 src 到 dst
func CopyFile(ctx context.Context, src, dst string) error {
	// 打开源文件
	source, err := os.Open(src)
	if err != nil {
		slog.ErrorContext(ctx, "os.Open error", "error", err, "src", src)
		return err
	}

	return CopyBodyFile(ctx, source, dst)
}

// FileList 收集完整的文件列表
func FileList(ctx context.Context, dirname string) (files []string, err error) {
	err = filepath.WalkDir(
		dirname,
		func(path string, dir os.DirEntry, direrr error) error {
			if direrr != nil {
				return direrr
			}

			// 只收集文件，跳过目录
			if dir.IsDir() {
				return nil
			}

			files = append(files, path)
			return nil
		},
	)
	if err != nil {
		slog.ErrorContext(ctx, "filepath.WalkDir error", "error", err, "dirname", dirname)
	}
	return
}

// OpenFile 打开文件, 低频率 api
func OpenFile(ctx context.Context, path string) (file *os.File, err error) {
	// 检查文件是否存在
	err = CreateDir(ctx, path)
	if err != nil {
		return
	}

	// os.OpenFile 内部有 runtime.SetFinalizer(f.file, (*file).close), 对象释放时候会 GC 1 close -> GC 2 free
	file, err = os.OpenFile(path, os.O_RDWR, 0o664)
	if err != nil {
		slog.ErrorContext(ctx, "os.OpenFile(path, os.O_RDWR, 0o664) error", "error", err, "path", path)
	}
	return
}
