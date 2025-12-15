// Package times provides utility functions for working with time.
package times

import (
	"fmt"
	"time"
)

// ShanghaiLoction 获取当前时间戳（东八区时间），格式化为易读的格式
// 使用 UTC+8 时区（东八区，即北京时间）
// LoadLocation 加载 "Asia/Shanghai" 时区，这是标准的东八区时区标识
var ShanghaiLoction *time.Location

func init() {
	var err error

	ShanghaiLoction, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		// 如果加载时区失败，记录错误但继续使用 UTC 时间
		// 这种情况极少发生，但为了健壮性需要处理
		fmt.Println("Failed to load Asia/Shanghai timezone, using UTC", err)
		ShanghaiLoction = time.Local
	}
}

// NowString 获取当前时间戳（东八区时间），格式化为易读的格式
func NowString() string {
	return time.Now().In(ShanghaiLoction).String()
}
