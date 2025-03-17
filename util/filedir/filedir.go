package filedir

import (
	"io"
	"os"
	"path/filepath"
)

// CreateDir 根据 path 创建 dir
func CreateDir(path string) error {
	// 如果文件不存在，尝试创建文件所在的目录
	dir := filepath.Dir(path)

	// 检查文件是否存在
	_, err := os.Stat(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	// 确保目录存在，如果不存在则创建
	return os.MkdirAll(dir, os.ModePerm)
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
	exists, _ := ExistE(filename)
	return exists
}

// ExistE 判断路径（文件或目录）是否存在
func ExistE(filepath string) (exists bool, err error) {
	_, err = os.Stat(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // 路径不存在
		}
		return false, err // 其他错误（如权限问题）
	}
	return true, nil // 路径存在（无论是文件还是目录）
}

// CopyWriter src file copy io.Writer
func CopyWriter(src string, writer io.Writer) error {
	// 打开源文件
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	// 使用 io.Copy 进行高效复制
	_, err = io.Copy(writer, source)
	return err
}

// CopyBodyFile resp.Body write file low api
func CopyBodyFile(body io.ReadCloser, dst string) error {
	defer body.Close()

	// 创建目标文件
	dest, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dest.Close()

	// 使用 io.Copy 进行高效复制
	_, err = io.Copy(dest, body)
	return err
}

// CopyFile 复制文件 src 到 dst
func CopyFile(src, dst string) error {
	// 打开源文件
	source, err := os.Open(src)
	if err != nil {
		return err
	}

	return CopyBodyFile(source, dst)
}

func CopyFileSync(src, dst string) error {
	// 打开源文件
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	// 创建目标文件
	dest, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dest.Close()

	// 使用 io.Copy 进行高效复制
	_, err = io.Copy(dest, source)
	if err != nil {
		return err
	}

	// 确保数据写入磁盘
	return dest.Sync()
}

// FileList 收集完整的文件列表
func FileList(dirname string) (files []string, err error) {
	err = filepath.WalkDir(dirname, func(path string, d os.DirEntry, err error) error {
		if err == nil {
			// 只收集文件，跳过目录
			if !d.IsDir() {
				files = append(files, path)
			}
		}
		return err
	})
	return
}

// ReadString os.ReadFile []byte -> string
func ReadString(filename string) (text string, err error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}
	text = string(data)
	return
}
