package chain

import (
	"os"
	"path/filepath"
	"strings"
)

var ExePath = os.Args[0]

var ExeName = filepath.Base(ExePath)

// ExeDir 获取可执行文件所在目录, 结尾不带 '/'
var ExeDir = filepath.Dir(ExePath)

func Hostname() string {
	// 获取容器的 hostname（通常是容器的短 ID）
	hostname, err := os.Hostname()
	if err == nil {
		return hostname
	}

	return UUID()
}

var DefaultRotatingFile string

func init() {
	ext := filepath.Ext(ExeName)
	ExeName = strings.TrimSuffix(ExeName, ext)

	// {exe path dir}/logs/{exe name}-{hostname}.log
	DefaultRotatingFile = filepath.Join(ExeDir, "logs", ExeName+"-"+Hostname()+".log")
}
