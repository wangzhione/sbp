# trace

`logs` [最初版本](https://github.com/wangzhione/sbp/commit/31ce0c165f3aef210926dc5f9ba5f7f08adb0b35#diff-9dc24e6d44b4f20a2c5d287b7560dc229700661a1eeb5f18c15fe84864687d80) 是个简单手写 log demo. 

后面主要是采用 Go 官方的 `slog` 轮子直接实例化, 整体业务开发会更加简单和高效.

**用法**

```
import "github.com/wangzhione/sbp/util/chain"

InitSLog()
```

impor 后 可以无缝使用 slog 进行 InfoContext or WarnContext or ErrorContext. 也可以参照其内部代码, 在 main func 初始化, 用业务自己的自定义 slog. 其中 chain.Key 各个环节交互唯一 key 串

***

## trace

```
import "sbp/util/chain"


// step 1: First 注入 trace id

// WithContext add trace id to context
func WithContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, key, id)
}

// or 

// Request init http.Request and return request id
func Request(r *http.Request) (req *http.Request, requestID string) {
	// 获取或生成 requestID
	requestID = r.Header.Get(Key)
	if requestID == "" {
		requestID = idhash.UUID()
	}
	// 注入 requestID 到 Context
	ctx := WithContext(r.Context(), requestID)

	req = r.WithContext(ctx)

	return
}


// step 2:  Second 获取 trace id
func GetTraceID(c context.Context) string 
```


## `道常无为而无不为` 
