
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

## 介绍 

deep copy 覆盖各种语法层面 case 写正确很困难, 这里也这是投机取消, 遇到循环引用, 和深度过深的 any 直接中断报错. 

**!!!当前想用 deep copy 时候一定要去思考, 能不能不用!!!**

很难找到一个案例, 非用不可, 未来会重新思考是否有保留的必要. 如果临时一用一次, json marshal unmarshal 未尝不可. 是吧
