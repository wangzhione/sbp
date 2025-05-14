package filedir

// 操作 disk 磁盘部分, util 工具

import (
	"fmt"
	"sort"
	"testing"

	"github.com/shirou/gopsutil/v4/disk"
)

func TestGetDiskUsage(t *testing.T) {
	usedPercent, err := GetDiskUsage(ctx, ".")
	if err != nil {
		t.Fatalf("❌ GetDiskUsage failed: %v", err)
	}

	if usedPercent <= 0 || usedPercent > 100 {
		t.Errorf("❌ unexpected usage percent: %.2f", usedPercent)
	} else {
		t.Logf("✅ 磁盘使用率: %.2f%%", usedPercent)
	}
}

type DiskInfo struct {
	Mountpoint  string
	Fstype      string
	Total       uint64
	Free        uint64
	Used        uint64
	UsedPercent float64
}

func formatSize(size uint64) string {
	const (
		_          = iota
		KB float64 = 1 << (10 * iota)
		MB
		GB
		TB
		PB
	)

	s := float64(size)

	switch {
	case s >= PB:
		return fmt.Sprintf("%.2f PB", s/PB)
	case s >= TB:
		return fmt.Sprintf("%.2f TB", s/TB)
	case s >= GB:
		return fmt.Sprintf("%.2f GB", s/GB)
	case s >= MB:
		return fmt.Sprintf("%.2f MB", s/MB)
	case s >= KB:
		return fmt.Sprintf("%.2f KB", s/KB)
	default:
		return fmt.Sprintf("%d B", size)
	}
}

func Test_Disk(t *testing.T) {
	partitions, err := disk.Partitions(true)
	if err != nil {
		panic(err)
	}

	var infos []DiskInfo
	for _, p := range partitions {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil || usage.Total == 0 {
			continue
		}

		infos = append(infos, DiskInfo{
			Mountpoint:  p.Mountpoint,
			Fstype:      usage.Fstype,
			Total:       usage.Total,
			Free:        usage.Free,
			Used:        usage.Used,
			UsedPercent: usage.UsedPercent,
		})
	}

	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Total > infos[j].Total
	})

	for _, d := range infos {
		fmt.Printf("挂载点: %-12s | 总容量: %-9s | 可用: %-9s | 使用率: %6.2f%% | FS: %s\n",
			d.Mountpoint,
			formatSize(d.Total),
			formatSize(d.Free),
			d.UsedPercent,
			d.Fstype,
		)
	}
}
