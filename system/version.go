package system

import (
	"runtime"
	"runtime/debug"
)

// Linux 默认是服务部署的最终服务器, 方便利用 system.Linux 默认做一些特殊处理逻辑
const Linux bool = runtime.GOOS == "linux"

/*
 runtime.GOOS 是 Go 语言中的一个常量，用于获取当前操作系统的名称。它的枚举值包括但不限于：

 windows
 linux
 darwin (macOS)
 freebsd
 openbsd
 netbsd
 android
 ios
 js (用于 Go 编译为 JavaScript)
 plan9
 solaris
*/

const Windows bool = runtime.GOOS == "windows"

const Darwin bool = runtime.GOOS == "darwin"

// BuildVersion 完整编译版本
var BuildVersion string = runtime.Version()

// GitVersion 项目发布时候代码 git 版本信息 | git rev-parse HEAD
// 依赖下面类型的 build 编译方式
// $env:CGO_ENABLED="0"; $env:GOOS="linux"; $env:GOARCH="amd64"; go build -trimpath -buildvcs=true -o {target} .
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
