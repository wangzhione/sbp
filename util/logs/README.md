# logs

`logs` is the package which provides logging functions for all the utilities in sbpkg.

`logs` 是对 Go 官方的 slog 简单事列初始化. 使用需要 import

```
import _ "sbpkg/util/logs"

```

随后, slog 进行 InfoContext or WarnContext or ErrorContext. 也可以参照其设计, 在 main func 初始化, 用业务自己的自定义 slog
