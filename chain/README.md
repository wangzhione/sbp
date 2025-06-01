# slog

采用 Go 官方的 `slog` 轮子直接实例化, 整体业务开发会更加简单和高效.

## uuid

`chain.UUID()` 是 uuid version 4 random 算法, 默认返回不带 '`-`' 风格的小写串

## trace

```
import "github.com/wangzhione/sbp/chain"


// step 1: First 注入 trace id

// WithContext add trace id to context
func WithContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, xRquestID, traceID)
}

// or 

// Request init http.Request and return request id
func Request(r *http.Request) (req *http.Request, requestID string) {
	// 获取或生成 requestID
	requestID = r.Header.Get(XRquestID)
	if requestID == "" {
		requestID = UUID()
	}
	// 注入 requestID 到 Context
	ctx := WithContext(r.Context(), requestID)

	req = r.WithContext(ctx)

	return
}


// step 2:  Second 获取 trace id
func GetTraceID(c context.Context) string 
```

## slog 用法 Ⅴ

```
import "github.com/wangzhione/sbp/chain"

chain.InitSLog()
```

import 后 可以无缝使用 slog 进行 InfoContext or WarnContext or ErrorContext. 也可以参照其内部代码, 在 main func 初始化, 用业务自己的自定义 slog. 其中 chain.XRquestID 各个环节交互唯一 trace key 串

也可以使用 `chain.InitSlogRotatingFile()` 进行日维度日志文件收集

***

