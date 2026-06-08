package system

import (
	"os"
	"path/filepath"
	"strings"
)

var (
	ExePath          = os.Args[0]                          // ExePath 获取可执行文件路径(相对路径 or 绝对路径)
	ExeDir           = filepath.Dir(ExePath)               // ExeDir 获取可执行文件所在目录, 结尾不带 '/'
	ExeName          = filepath.Base(ExePath)              // ExeName 获取不带路径的可执行文件名
	ExeExt           = filepath.Ext(ExeName)               // ExeExt 获取可执行文件名的扩展名
	ExeNameSuffixExt = strings.TrimSuffix(ExeName, ExeExt) // ExeNameSuffixExt 获取可执行文件名, 不包含扩展名
)

// Hostname 获取主机名 or 容器短 ID
var Hostname = func() string {
	// 获取容器的 hostname（通常是容器的短 ID）
	hostname, err := os.Hostname()
	if err != nil {
		println("os.Hostname error", err.Error())
		// 失败时候用 UUID 作为主机标识
		return UUID()
	}
	return hostname
}()

// Exist 判断路径（文件或目录）是否存在
func Exist(path string) (exists bool, err error) {
	_, err = os.Stat(path)
	if err == nil {
		// 路径存在（无论是文件还是目录）
		return true, nil
	}
	if os.IsNotExist(err) {
		// 路径不存在
		return false, nil
	}
	// 其他错误（如权限问题）, 但对当前用户而言是不存在
	return false, err
}
