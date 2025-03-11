package filedir

import (
	"os"
	"path/filepath"
	"strings"
)

var ExePath = os.Args[0]

var ExeName = filepath.Base(ExePath)

// ExeDir 获取可执行文件所在目录, 结尾不带 '/'
var ExeDir = filepath.Dir(ExePath)

func init() {
	ext := filepath.Ext(ExeName)
	ExeName = strings.TrimSuffix(ExeName, ext)
}
