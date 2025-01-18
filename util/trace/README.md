# trace

`logs` [最初版本](https://github.com/wangzhione/sbp/commit/31ce0c165f3aef210926dc5f9ba5f7f08adb0b35#diff-9dc24e6d44b4f20a2c5d287b7560dc229700661a1eeb5f18c15fe84864687d80) 是个简单手写 log demo. 

后面主要是采用 Go 官方的 `slog` 轮子直接实例化, 整体业务开发会更加简单和高效.

**用法**

```
import _ "sbp/util/trace"
```

impor 后 可以无缝使用 slog 进行 InfoContext or WarnContext or ErrorContext. 也可以参照其内部代码, 在 main func 初始化, 用业务自己的自定义 slog. 

***

## trace

```
import "sbp/util/trace"


// First 注入 trace id
func WithTraceID(c *context.Context) string

// Second 获取 trace id
func GetTraceID(c context.Context) string 
```


## `道常无为而无不为` 
