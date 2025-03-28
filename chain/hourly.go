package chain

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

type hourlylogger struct {
	*os.File
}

func starthourlylogger() error {
	our := &hourlylogger{} // our 类似跨函数闭包
	if err := our.rotate(); err != nil {
		return err
	}
	go our.rotateloop()
	return nil
}

// Exist 判断路径（文件或目录）是否存在
func Exist(path string) (exists bool, err error) {
	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // 路径不存在
		}
		return false, err // 其他错误（如权限问题）, 但对当前用户而言是不存在
	}
	return true, nil // 路径存在（无论是文件还是目录）
}

func (our *hourlylogger) rotate() error {
	hours := time.Now().Format("2006010215") // e.g. 2025032815

	// {exe path dir}/logs/{exe name}-{2025032815}-{hostname}.log
	filename := filepath.Join(LogsDir, ExeName+"-"+hours+"-"+Hostname()+".log")

	if our.File != nil && our.Name() == filename {
		found, err := Exist(filename)
		if found || err != nil {
			return err
		}
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		println("hourlylogger os.OpenFile error", err, filename)
		return err
	}

	options := &slog.HandlerOptions{
		Level: EnableLevel,
	}

	hourly := slog.NewJSONHandler(io.MultiWriter(os.Stdout, file), options)

	slog.SetDefault(slog.New(&TraceHandler{hourly}))

	_ = our.Close() // os.OpenFile 有兜底 runtime.SetFinalizer(f.file, (*file).close) 😂
	our.File = file

	// 历史日志清理
	our.sevenday()

	return nil
}

func (our *hourlylogger) rotateloop() {
	for {
		now := time.Now()
		// 下一个整点
		next := now.Truncate(time.Hour).Add(time.Hour)
		sleep := next.Sub(now)
		time.Sleep(sleep)

		_ = our.rotate() // 业务 println 打印 error 日志兜底
	}
}

var DefaultCleanDay = 15 // 15 天前, 有时候过 7 天假期, 回来 7 天日志没了 ...

func (our *hourlylogger) sevenday() {
}
