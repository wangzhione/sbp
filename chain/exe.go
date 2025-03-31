package chain

import (
	"os"
	"path/filepath"
	"strings"
)

var ExePath = os.Args[0]

var ExeName = filepath.Base(ExePath)

var ExeExt = filepath.Ext(ExeName)

var ExeNameSuffixExt = strings.TrimSuffix(ExeName, ExeExt)

// ExeDir 获取可执行文件所在目录, 结尾不带 '/'
var ExeDir = filepath.Dir(ExePath)

func hostname() string {
	// 获取容器的 hostname（通常是容器的短 ID）
	hostname, err := os.Hostname()
	if err == nil {
		return hostname
	}

	return UUID()
}

var Hostname = hostname()
