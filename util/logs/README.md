# logs

`logs` is the package which provides logging functions for all the utilities in sbpkg.

`logs` [最初版本](https://github.com/wangzhione/sbpkg/commit/31ce0c165f3aef210926dc5f9ba5f7f08adb0b35#diff-9dc24e6d44b4f20a2c5d287b7560dc229700661a1eeb5f18c15fe84864687d80) 是个简单手写 log demo. 后面主要是采用 Go 官方的 `slog` 轮子直接实例化. 使用需要 import

```
import _ "sbpkg/util/logs"
```

然后 可以无缝使用 slog 进行 InfoContext or WarnContext or ErrorContext. 也可以参照其内部代码, 在 main func 初始化, 用业务自己的自定义 slog.

***

## `道常无为而无不为` 
