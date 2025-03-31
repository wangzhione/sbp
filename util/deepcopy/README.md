
[mohae/deepcopy](https://github.com/mohae/deepcopy) -> [xieyuschen/deepcopy](https://github.com/xieyuschen/deepcopy) -> [wangzhione/deepcopy](https://github.com/wangzhione/sbp/tree/master/util/deepcopy)

# deepcopy

用户举例

```Go
import "github.com/wangzhione/sbp/util/deepcopy"

func deepcopy.Clone[T any](ctx context.Context, src T) (dst T) {
	i, err := deepcopy.Copy(src)
	if err != nil {
		slog.ErrorContext(ctx, "deepcopy.Copy panic error", "error", err, "src", src)
		return
	}
	dst, _ = i.(T)
	return
}
```