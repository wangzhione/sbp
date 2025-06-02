package chain

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"
)

func Test_search(t *testing.T) {
	text := `
logs/abc-materialefficiencytool-2025051404-ms-2scj6hpg-1-6c44dcc954-rfnhf.log
logs/segmentclips-2025041115-nb-1282427673004035712-9qrao4gnd4e8.log
`

	// 正则：匹配 logs/... 中的 10 位数字段
	re := regexp.MustCompile(`logs/(?:[^/-]+-)*(\d{10,12})-`)
	matches := re.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		if len(match) > 1 {
			fmt.Println(match[1]) // 输出 2025051404、2025041115 等
		}
	}

	path := `F:\zhi.wang\cutoff\logs\segmentclips.exe-2025050617-MSI.log`
	// 正则：匹配 logs/... 中的 10 位数字段
	onematches := Dre.FindStringSubmatch(path)
	if len(onematches) < 2 {
		println("hourlylogger reD.FindStringSubmatch error", strings.Join(onematches, " "), Hostname, path)
	}
	println(onematches[1])
}

func Test_filename(t *testing.T) {
	filename := `logs/segmentclips-2025041115-nb-1282427673004035712-9qrao4gnd4e8.log`

	path := strings.TrimSuffix(filename, "nb-1282427673004035712-9qrao4gnd4e8.log")

	t.Log(filename, path)

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
		return
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
		return
	}

	// 提取中间的时间字符串
	timeStr := path[start+1 : end]

	// 解析时间
	xtx, err := time.Parse("2006010215", timeStr)
	if err != nil {
		println("hourlylogger filepath.WalkDir time.Parse error", err.Error(), path)
		return
	}

	t.Log("Success", xtx)
}
