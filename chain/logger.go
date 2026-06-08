package chain

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/wangzhione/sbp/system"
)

type GetfileFn func(logdir string) (now time.Time, filename string)

// GetfileByDay 按天切割日志
// 生成的日志文件名: {exe path dir}/logs/{20250522}-{exe name}-{hostname}.log
func GetfileByDay(logdir string) (now time.Time, filename string) {
	now = time.Now()
	times := now.Format("20060102")
	// {exe path dir}/logs/{20250522}-{exe name}-{hostname}.log
	filename = filepath.Join(logdir, times+"-"+system.ExeName+"-"+system.Hostname+".log")
	println("GetfileByDay day init log", filename)
	return
}

func GetfileByHour(logdir string) (now time.Time, filename string) {
	now = time.Now()
	times := now.Format("2006010215")
	// {exe path dir}/logs/{2025032815}-{exe name}-{hostname}.log
	filename = filepath.Join(logdir, times+"-"+system.ExeName+"-"+system.Hostname+".log")
	println("GetfileByHour init log", filename)
	return
}

// Startlogger 启动一个 slog 实例
func Startlogger(logdir string, getfilefn GetfileFn, closecutoff bool) error {
	if logdir == "" {
		// Log Dir 默认日志目录 {exe dir}/logs
		logdir = filepath.Join(system.ExeDir, "logs")
	}
	if getfilefn == nil {
		// 默认按天切割日志
		getfilefn = GetfileByDay
	}

	// our 类似跨函数闭包
	our := &timelogger{
		LogDir:    logdir,
		getfilefn: getfilefn,
	}

	err := os.MkdirAll(our.LogDir, os.ModePerm)
	if err != nil {
		println("os.MkdirAll error", our.LogDir)
		return err
	}

	if err := our.rotate(); err != nil {
		return err
	}

	if !closecutoff {
		// 启动日志切割循环, 固定时间检查一次是否需要切割日志
		go our.rotateloop()
	}

	return nil
}

type timelogger struct {
	*os.File
	lasttime time.Time

	getfilefn GetfileFn
	LogDir    string // ★ 默认 log dir 在 {exe dir}/logs
}

func (our *timelogger) rotate() error {
	now, filename := our.getfilefn(our.LogDir)

	if our.File != nil && our.Name() == filename {
		found, err := system.Exist(filename)
		if found || err != nil {
			return err
		}
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o664)
	if err != nil {
		println("rotate os.OpenFile error", err.Error(), filename)
		return err
	}

	slog.SetDefault(slog.New(&TraceHandler{
		slog.NewJSONHandler(io.MultiWriter(os.Stdout, file), &slog.HandlerOptions{
			Level: EnableLevel,
		}),
	}))

	_ = our.File.Close() // os.OpenFile 有兜底 runtime.SetFinalizer(f.file, (*file).close) 😂
	our.File = file

	// 历史日志清理
	our.sevenday(now)

	return nil
}

func (our *timelogger) rotateloop() {
	for {
		now := time.Now()
		// 下一个整点, 计算需要 sleep 时间
		next := now.Truncate(time.Hour).Add(time.Hour)
		time.Sleep(next.Sub(now))
		_ = our.rotate()
	}
}

// DefaultCleanTime 默认 15 天前, 有时候过 7 天假期, 回来 7 天日志没了 ...
var DefaultCleanTime = -15 * 24 * time.Hour

// DefaultCheckTime sevenday 每次检查是否要清理历史日志时间间隔
var DefaultCheckTime = 7 * time.Hour

func (our *timelogger) sevenday(now time.Time) {
	if now.Sub(our.lasttime) < DefaultCheckTime {
		// 时间间隔太小直接返回
		return
	}
	our.lasttime = now

	cutoff := now.Add(DefaultCleanTime)
	// 尝试清理历史文件
	var files []string
	err := filepath.WalkDir(
		our.LogDir,
		func(path string, dir os.DirEntry, direrr error) error {
			if direrr != nil {
				return direrr
			}

			// 只收集文件，跳过目录
			if dir.IsDir() {
				return nil
			}

			filename := filepath.Base(path)

			// ★ 仅处理 .log 文件
			if filepath.Ext(filename) != ".log" {
				return nil
			}

			if len(filename) < 8 {
				// 文件名太短, 不符合格式
				return nil
			}

			// 提取开始的时间字符串
			timeStr := filename[:8]
			// 解析时间
			t, err := time.Parse("20060102", timeStr)
			if err != nil {
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
		println("sevenday filepath.WalkDir error", err.Error(), our.LogDir)
	}

	for i, file := range files {
		err = os.Remove(file)
		if err != nil {
			println("sevenday os.Remove error", i, file, err.Error())
			continue
		}
		println("sevenday os.Remove success", file)
	}
}
