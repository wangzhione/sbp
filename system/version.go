package system

import (
	"runtime"
	"runtime/debug"
)

// BuildVersion 完整编译版本
var BuildVersion string = runtime.Version()

// GitVersion 项目发布时候代码 git 版本信息 | git rev-parse HEAD
var GitVersion string

// GitCommitTime 最近一次提交时间（来自 vcs.time）
var GitCommitTime string

func init() {
	info, ok := debug.ReadBuildInfo()

	if !ok {
		println("debug.ReadBuildInfo() return no ok")
		return
	}

	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			GitVersion = setting.Value
		case "vcs.time":
			GitCommitTime = setting.Value
		}
	}
}
