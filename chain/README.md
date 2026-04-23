# chain

`chain` 主要提供两类能力:

- 基于 `context.Context` / `http.Request` 的 trace 透传
- 基于 Go 官方 `log/slog` 的默认日志初始化, 并自动补充 `X-Request-Id` 与调用位置

## 1. trace

核心常量:

```go
const XRquestID = "X-Request-Id"
```

常用函数:

```go
func Context() context.Context
func WithContext(ctx context.Context, traceID string) context.Context
func GetTraceID(ctx context.Context) string
func TraceID(ctx context.Context) string
func CopyTrace(ctx context.Context) context.Context
func CopyContext(ctx context.Context, keys ...any) context.Context
func Request(r *http.Request, headers ...string) (*http.Request, string)
```

说明:

- `Context()` 返回一个带随机 trace id 的基础 context
- `WithContext()` 手动注入 trace id
- `GetTraceID()` 只取值, 没有就返回空串
- `TraceID()` 保证返回非空 trace id, 取不到时会新生成
- `CopyTrace()` / `CopyContext()` 会脱离原 context 的 timeout / cancel, 适合异步任务继续携带 trace
- `Request()` 优先从传入的 header 列表里取 trace id, 否则再读 `X-Request-Id`, 最后兜底生成新的 id

### 示例: 手动注入 trace

```go
package main

import (
	"context"
	"log/slog"

	"github.com/wangzhione/sbp/chain"
	"github.com/wangzhione/sbp/system"
)

func main() {
	chain.InitSLog()

	ctx := chain.WithContext(context.Background(), system.UUID())

	slog.InfoContext(ctx, "hello chain")
}
```

### 示例: 从 HTTP 请求中接入 trace

```go
package main

import (
	"log/slog"
	"net/http"

	"github.com/wangzhione/sbp/chain"
)

func handler(w http.ResponseWriter, r *http.Request) {
	req, traceID := chain.Request(r)

	slog.InfoContext(req.Context(), "request in", "traceID", traceID)
	w.WriteHeader(http.StatusOK)
}
```

### 示例: 异步任务复制 trace

```go
go func(ctx context.Context) {
	ctx = chain.CopyTrace(ctx)
	slog.InfoContext(ctx, "async worker start")
}(ctx)
```

## 2. slog

### `InitSLog`

初始化默认 `slog` 到标准输出:

```go
chain.InitSLog()
```

输出时会自动补两类字段:

- `X-Request-Id`: 从 `ctx` 里提取 trace id
- `code`: 调用位置, 格式类似 `slog_test.go:26:TestInitSLogRotatingFile`

日志等级默认来自:

```go
var EnableLevel slog.Level = slog.LevelDebug
```

如果业务有配置中心或命令行参数, 可以先改 `chain.EnableLevel`, 再初始化日志。

### `InitSLogRotatingFile`

初始化按时间切分的文件日志, 默认行为:

- 输出到标准输出 + 日志文件
- 默认日志目录: `{exe dir}/logs`
- 默认按天切割
- 默认清理 15 天前日志

```go
if err := chain.InitSLogRotatingFile(); err != nil {
	slog.Error("chain.InitSLogRotatingFile error", "error", err)
	chain.InitSLog()
}
```

相关变量:

```go
var LogsDir = filepath.Join(system.ExeDir, "logs")
var DefaultGetFile = GetfileByDay
var DefaultCleanTime = -15 * 24 * time.Hour
var DefaultCheckTime = 7 * time.Hour
```

如果想按小时切割:

```go
chain.DefaultGetFile = chain.GetfileByHour
err := chain.InitSLogRotatingFile()
```

如果不想启动后台轮转协程:

```go
err := chain.InitSLogRotatingFile(true)
```

## 3. 关于 trace id

`chain` 包本身不提供 `chain.UUID()`。

当前 trace id 默认来自:

```go
system.UUID()
```

它返回的是不带 `-` 的 32 位小写 uuid v4 风格字符串, 例如:

```text
22ba3cffc8de4a2d9dc8a95d09ed03e1
```

## 4. 推荐初始化方式

```go
package main

import (
	"context"
	"log/slog"

	"github.com/wangzhione/sbp/chain"
)

func main() {
	if err := chain.InitSLogRotatingFile(); err != nil {
		chain.InitSLog()
		slog.ErrorContext(context.Background(), "InitSLogRotatingFile error", "error", err)
	}

	ctx := chain.Context()
	slog.InfoContext(ctx, "service start")
}
```
