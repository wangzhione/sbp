package chain

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type daylogger struct {
	*os.File
	lasttime time.Time

	LogsDir string // LogsDir ★ 默认 log dir 在 {exe dir}/logs
}

func startdaylogger() error {
	our := &daylogger{ // our 类似跨函数闭包
		LogsDir: filepath.Join(ExeDir, "logs"),
	}

	err := os.MkdirAll(our.LogsDir, os.ModePerm)
	if err != nil {
		println("os.MkdirAll error", our.LogsDir)
		return err
	}

	if err := our.rotate(); err != nil {
		return err
	}
	go our.rotateloop()
	return nil
}

func (our *daylogger) rotate() error {
	now := time.Now()

	days := now.Format("20060102") // e.g. 20250522
	// {exe path dir}/logs/{exe name}-{20250522}-{hostname}.log
	filename := filepath.Join(our.LogsDir, ExeName+"-"+days+"-"+Hostname+".log")
	print("rotate day init log", Hostname, filename)

	if our.File != nil && our.Name() == filename {
		found, err := Exist(filename)
		if found || err != nil {
			return err
		}
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		println("daylogger os.OpenFile error", err.Error(), filename)
		return err
	}

	options := &slog.HandlerOptions{
		Level: EnableLevel,
	}

	stdoutandfile := io.MultiWriter(os.Stdout, file)

	var hourly slog.Handler
	if EnableText() {
		hourly = slog.NewTextHandler(stdoutandfile, options)
	} else {
		hourly = slog.NewJSONHandler(stdoutandfile, options)
	}

	slog.SetDefault(slog.New(&TraceHandler{hourly}))

	_ = our.Close() // os.OpenFile 有兜底 runtime.SetFinalizer(f.file, (*file).close) 😂
	our.File = file

	// 历史日志清理
	our.sevenday(now)

	return nil
}

func (our *daylogger) rotateloop() {
	for {
		now := time.Now()
		// 下一个整点, 计算需要 sleep 时间
		next := now.Truncate(time.Hour).Add(time.Hour)
		time.Sleep(next.Sub(now))

		_ = our.rotate() // 业务 println 打印 error 日志兜底
	}
}

func (our *daylogger) sevenday(now time.Time) {
	if now.Sub(our.lasttime) < DefaultCheckTime {
		// 时间间隔太小直接返回
		return
	}
	our.lasttime = now

	cutoff := now.Add(-DefaultCleanTime)
	// 尝试清理历史文件
	var files []string
	err := filepath.WalkDir(
		our.LogsDir,
		func(path string, dir os.DirEntry, direrr error) error {
			if direrr != nil {
				return direrr
			}

			// 只收集文件，跳过目录
			if dir.IsDir() {
				return nil
			}

			// fix `logs/materialefficiencytool-20250514-ms-2scj6hpg-1-6c44dcc954-rfnhf.log` bug
			// fix `logs/segmentclips-20250411-nb-1282427673004035712-9qrao4gnd4e8.log` bug

			// 正则：匹配 logs/... 中的 8 位数字段
			matches := Dre.FindStringSubmatch(path)
			if len(matches) < 2 {
				println("daylogger reD.FindStringSubmatch error", strings.Join(matches, " "), Hostname, path)
				files = append(files, path)
				return nil
			}

			// 提取中间的时间字符串
			timeStr := matches[1]
			if len(timeStr) > 8 {
				timeStr = timeStr[:8]
			}

			// 解析时间
			t, err := time.Parse("20060102", timeStr)
			if err != nil {
				println("daylogger filepath.WalkDir time.Parse error", err.Error(), Hostname, path)
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
		println("daylogger filepath.WalkDir error", err.Error(), our.LogsDir)
		return
	}

	for _, file := range files {
		err = os.Remove(file)
		if err != nil {
			println("daylogger os.Remove error", err.Error(), file)
		} else {
			println("daylogger os.Remove success", file)
		}
	}
}
