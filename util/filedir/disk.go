package filedir

import (
	"context"
	"log/slog"

	"github.com/shirou/gopsutil/v4/disk"
)

// 操作 disk 磁盘部分, util 工具

// GetDiskUsage 返回 path 目录所在分区的已使用百分比（如 68.32）
func GetDiskUsage(ctx context.Context, path string) (usedPercent float64, err error) {
	usageStat, err := disk.Usage(path)
	if err != nil {
		slog.ErrorContext(ctx, "disk.Usage error", "path", path, "error", err)
		return 0, err
	}
	return usageStat.UsedPercent, nil
}
