// Package system provides utilities for detecting the operating system and related logic.
package system

import (
	"runtime"
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
