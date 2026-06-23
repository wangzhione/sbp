package filedir

import (
	"context"
	"log/slog"
	"path/filepath"
)

// RelPath 根据 basepath 计算 targpath 的相对路径, 并统一转换为 slash 分隔符.
//
// 参数:
//   - ctx: 仅用于错误日志上下文, 不参与路径计算.
//   - basepath: 基准路径, 可以是相对路径或绝对路径.
//   - targpath: 目标路径, 可以是相对路径或绝对路径.
//
// 返回:
//   - relPath: targpath 相对 basepath 的路径; targpath 等于 basepath 时为 ".".
//   - err: filepath.Rel 返回的错误; 例如 Windows 下不同盘符路径无法计算相对路径.
//
// 注意:
//   - filepath.Rel 只按路径规则计算, 不检查路径是否真实存在.
//   - filepath.ToSlash 会把系统路径分隔符转换为 "/", 便于配置、URL-like 路径或跨平台保存.
func RelPath(ctx context.Context, basepath, targpath string) (relPath string, err error) {
	relPath, err = filepath.Rel(basepath, targpath)
	if err != nil {
		slog.ErrorContext(ctx, "filepath.Rel", "error", err, "basepath", basepath, "targpath", targpath)
		return
	}
	relPath = filepath.ToSlash(relPath)
	return
}

// AbsPath 将 path 解析为绝对路径.
//
// 参数:
//   - ctx: 仅用于错误日志上下文, 不参与路径计算.
//   - path: 待解析路径; 可以是相对路径、绝对路径、"." 或 "..".
//
// 返回:
//   - abspath: 解析后的绝对路径; 例如 "." 通常会解析为当前工作目录.
//   - err: filepath.Abs 返回的错误; 失败时会记录原始 path 和当前 abspath.
//
// 注意:
//   - filepath.Abs 会先清理路径, 并在相对路径前拼接当前工作目录.
//   - 当前工作目录来自进程运行时环境, 因此不同启动目录会得到不同结果.
//   - 该函数只负责路径解析, 不检查目标文件或目录是否真实存在.
func AbsPath(ctx context.Context, path string) (abspath string, err error) {
	abspath, err = filepath.Abs(path)
	if err != nil {
		slog.ErrorContext(ctx, "filepath.Abs", "error", err, "path", path, "abspath", abspath)
		return
	}
	return
}

// - filepath.Dir(path)  	-> 目录
// - filepath.Base(path) 	-> 文件名
// - filepath.Ext(filename) -> 扩展名
// - strings.TrimSuffix(文件名, 扩展名) -> 不带扩展名的文件名
