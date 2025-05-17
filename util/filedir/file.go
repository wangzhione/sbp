package filedir

import (
	"archive/zip"
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
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

func GetModTime(ctx context.Context, path string) (time.Time, error) {
	fs, err := os.Stat(path)
	if err != nil {
		slog.ErrorContext(ctx, "os.Stat error", "path", path, "error", err)
		return time.Time{}, err
	}

	// ModTime() 用于返回文件或目录的 最后修改时间（mtime）
	return fs.ModTime(), nil
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
		slog.ErrorContext(ctx, "FileList error", "error", err, "dirname", dirname)
	}
	return
}

// AddFileToZip 将指定文件以 relPath 写入 zipWriter
func AddFileToZip(ctx context.Context, zipWriter *zip.Writer, filePath string, relPath string) error {
	// 打开文件
	f, err := os.Open(filePath)
	if err != nil {
		slog.ErrorContext(ctx, "Open file failed", "error", err, "file", filePath)
		return err
	}
	defer f.Close() // 确保文件关闭，防止 fd 泄漏

	// 创建 zip 条目
	fw, err := zipWriter.Create(relPath)
	if err != nil {
		slog.ErrorContext(ctx, "Create zip file entry failed", "error", err, "file", filePath)
		return err
	}

	// 写入文件内容
	if _, err := io.Copy(fw, f); err != nil {
		slog.ErrorContext(ctx, "Copy file content failed", "error", err, "file", filePath)
		return err
	}

	return nil
}

/*

// filepath.Walk 实际案例
err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {})

func ResponseWriterZipDir(ctx context.Context, w http.ResponseWriter, zipname string, dir string) {
	// 设置响应头，提前发送 attachment 信息
	w.Header().Set("Content-Disposition", "attachment; filename="+zipname)
	w.Header().Set("Content-Type", "application/zip")

	// 创建 ZIP writer，直接写入 HTTP 响应流
	zipWriter := zip.NewWriter(w)
	defer func() {
		// 尝试关闭 zipWriter，刷新缓存
		if err := zipWriter.Close(); err != nil {
			// 如果 zipWriter 关闭失败
			slog.ErrorContext(ctx, "Failed to finalize zipWriter", "error", err, "zipname", zipname)
			return
		}
	}()

	// 遍历目录
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			slog.ErrorContext(ctx, "Walk error", "error", err, "path", path)
			return err
		}

		// 构造 zip 文件中的相对路径（保留目录结构）
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			slog.ErrorContext(ctx, "Path Rel error", "error", err, "path", path)
			return err
		}

		// 忽略根路径（.）
		if relPath == "." {
			return nil
		}

		if info.IsDir() {
			// 是目录则创建空目录条目
			_, err := zipWriter.Create(relPath + "/")
			if err != nil {
				slog.ErrorContext(ctx, "Create zip folder entry failed", "error", err, "folder", relPath)
				return err
			}
			return nil
		}

		// 文件处理
		return AddFileToZip(ctx, zipWriter, path, relPath)
	})
	// 如果 Walk 过程中出错
	if err != nil {
		slog.ErrorContext(ctx, "Failed to walk and zip directory", "error", err, "dir", dir, "zipname", zipname)
		// 注意：不能再使用 http.Error，因为响应头已发，部分内容已写入
		return
	}
}
*/

// WriteFileIfNotExists 写入文件, 必须不存在才会写入
func WriteFileIfNotExists(ctx context.Context, path string, content []byte) (err error) {
	found, err := Exist(ctx, path)
	if err != nil {
		return
	}

	// 文件存在不再处理
	if found {
		return
	}

	return WriteFile(ctx, path, content)
}

func WriteFile(ctx context.Context, path string, content []byte) (err error) {
	// init 目录
	err = CreateDir(ctx, path)
	if err != nil {
		return
	}

	// 创建文件
	err = os.WriteFile(path, content, os.ModePerm)
	if err != nil {
		slog.ErrorContext(ctx, "os.WriteFile error", "error", err, "path", path, "len(content)", len(content))
	}
	return err
}
