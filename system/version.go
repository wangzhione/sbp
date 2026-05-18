// Package system provides helpers for runtime environment, process metadata,
// platform detection, and system-level identifiers.
package system

import (
	"runtime"
	"runtime/debug"
)

const (
	Aix       bool = runtime.GOOS == "aix"       // IBM AIX Unix 服务器平台
	Android   bool = runtime.GOOS == "android"   // Android 移动端或嵌入式平台
	Darwin    bool = runtime.GOOS == "darwin"    // macOS 桌面或服务器平台
	Dragonfly bool = runtime.GOOS == "dragonfly" // DragonFly BSD 系统平台
	Freebsd   bool = runtime.GOOS == "freebsd"   // FreeBSD 系统平台
	Hurd      bool = runtime.GOOS == "hurd"      // GNU Hurd 系统平台
	Illumos   bool = runtime.GOOS == "illumos"   // illumos/Solaris 衍生系统平台
	Ios       bool = runtime.GOOS == "ios"       // iOS 移动端平台
	Js        bool = runtime.GOOS == "js"        // JavaScript/WebAssembly 运行平台
	Linux     bool = runtime.GOOS == "linux"     // Linux 服务器或桌面平台
	Nacl      bool = runtime.GOOS == "nacl"      // Native Client 沙箱运行平台
	Netbsd    bool = runtime.GOOS == "netbsd"    // NetBSD 系统平台
	Openbsd   bool = runtime.GOOS == "openbsd"   // OpenBSD 系统平台
	Plan9     bool = runtime.GOOS == "plan9"     // Plan 9 系统平台
	Solaris   bool = runtime.GOOS == "solaris"   // Oracle Solaris 系统平台
	Wasip1    bool = runtime.GOOS == "wasip1"    // WASI Preview 1 WebAssembly 平台
	Windows   bool = runtime.GOOS == "windows"   // Windows 桌面或服务器平台
	Zos       bool = runtime.GOOS == "zos"       // IBM z/OS 大型机平台
)

// BuildGoVersion Go 编译器完整版本
var BuildGoVersion string = runtime.Version()

// GitVersion 项目发布时候代码 git 版本信息 | git rev-parse HEAD
// 依赖下面类型的 build 编译方式
// $env:CGO_ENABLED="0"; $env:GOOS="linux"; $env:GOARCH="amd64"; go build -trimpath -buildvcs=true -o {target} .
var GitVersion string

// GitLastCommitTime 最近一次提交时间（来自 vcs.time）
var GitLastCommitTime string

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
			GitLastCommitTime = setting.Value
		}
	}
}
