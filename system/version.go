package system

import (
	"runtime"
	"runtime/debug"
)

// BuildVersion 完整编译版本
var BuildVersion string = runtime.Version()

// GitVersion 项目发布时候代码 git 版本信息 | git rev-parse HEAD
var GitVersion string

func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		println("debug.ReadBuildInfo() return not ok")
	} else {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				GitVersion = setting.Value
				break
			}
		}
	}
}
