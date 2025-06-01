package chain

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
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

// Exist 判断路径（文件或目录）是否存在
func Exist(path string) (exists bool, err error) {
	_, err = os.Stat(path)
	if err == nil {
		return true, nil // 路径存在（无论是文件还是目录）
	}

	if os.IsNotExist(err) {
		return false, nil // 路径不存在
	}
	return false, err // 其他错误（如权限问题）, 但对当前用户而言是不存在
}

// LogStartEnd Wrapper function to log start and end times, and measure duration
func LogStartEnd(ctx context.Context, name string, fn func(context.Context) error) (err error) {
	start := time.Now()
	slog.InfoContext(ctx, "["+name+"] - Start", "time", start.Format("2006-01-02 15:04:05.000000"))

	// Execute the wrapped function with context
	err = fn(ctx)

	end := time.Now()
	elapsed := end.Sub(start)
	slog.InfoContext(ctx, "["+name+"] - End", "elapsed", elapsed.Seconds(), "time", end.Format("2006-01-02 15:04:05.000000"))
	return
}
