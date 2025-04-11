package chain

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type hourlylogger struct {
	*os.File
	lasttime time.Time
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
	if err == nil {
		return true, nil // 路径存在（无论是文件还是目录）
	}

	if os.IsNotExist(err) {
		return false, nil // 路径不存在
	}
	return false, err // 其他错误（如权限问题）, 但对当前用户而言是不存在
}

var Hostnamelog = Hostname + ".log"

func (our *hourlylogger) rotate() error {
	now := time.Now()
	hours := now.Format("2006010215") // e.g. 2025032815

	// {exe path dir}/logs/{exe name}-{2025032815}-{hostname}.log
	filename := filepath.Join(LogsDir, ExeName+"-"+hours+"-"+Hostnamelog)

	if our.File != nil && our.Name() == filename {
		found, err := Exist(filename)
		if found || err != nil {
			return err
		}
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		println("hourlylogger os.OpenFile error", err.Error(), filename)
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
	our.sevenday(now)

	return nil
}

func (our *hourlylogger) rotateloop() {
	for {
		now := time.Now()
		// 下一个整点, 计算需要 sleep 时间
		next := now.Truncate(time.Hour).Add(time.Hour)
		time.Sleep(next.Sub(now))

		_ = our.rotate() // 业务 println 打印 error 日志兜底
	}
}

var DefaultCleanTime = 15 * 24 * time.Hour // 默认 15 天前, 有时候过 7 天假期, 回来 7 天日志没了 ...

const DefaultCheckTime = 7 * time.Hour // sevenday 每次检查是否要清理历史日志时间间隔

// LogsDir ★ 默认 log dir 在 {exe dir}/logs
var LogsDir = filepath.Join(ExeDir, "logs")

func (our *hourlylogger) sevenday(now time.Time) {
	if now.Sub(our.lasttime) < DefaultCheckTime {
		// 时间间隔太小直接返回
		return
	}
	our.lasttime = now

	cutoff := now.Add(-DefaultCleanTime)
	// 尝试清理历史文件
	var files []string
	err := filepath.WalkDir(
		LogsDir,
		func(path string, dir os.DirEntry, direrr error) error {
			if direrr != nil {
				return direrr
			}

			// 只收集文件，跳过目录
			if dir.IsDir() {
				return nil
			}

			// fix `logs/segmentclips-2025041115-nb-1282427673004035712-9qrao4gnd4e8.log` bug
			path = strings.TrimSuffix(path, Hostnamelog)

			// {exe path dir}/logs/{exe name}-{2025032815}-{hostname}.log
			// 从后往前找两个 '-' 的位置
			// 第一次循环，从后往前找第一个 '-'（end）
			end := -1
			for i := len(path) - 1; i >= 0; i-- {
				if path[i] == '-' {
					end = i
					break
				}
			}
			if end == -1 {
				return nil
			}

			// 第二次循环，从 end-2 开始往前找第二个 '-'（start）
			start := -1
			for i := end - 2; i >= 0; i-- {
				if path[i] == '-' {
					start = i
					break
				}
			}
			if start == -1 {
				return nil
			}

			// 提取中间的时间字符串
			timeStr := path[start+1 : end]

			// 解析时间
			t, err := time.Parse("2006010215", timeStr)
			if err != nil {
				println("hourlylogger filepath.WalkDir time.Parse error", err.Error(), path)
				return nil
			}

			// 判断是否超过待删除时间, 如果没有超过直接返回
			if !t.Before(cutoff) {
				return nil
			}

			if path == our.Name() {
				// 特殊 case, 保留当前输出文件
				return nil
			}

			files = append(files, path)
			return nil
		},
	)
	if err != nil {
		println("hourlylogger filepath.WalkDir error", err.Error(), LogsDir)
		return
	}

	for _, file := range files {
		err = os.Remove(file)
		println("hourlylogger os.Remove error", err.Error(), file)
	}
}
