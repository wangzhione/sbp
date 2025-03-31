
[mohae/deepcopy](https://github.com/mohae/deepcopy) -> [xieyuschen/deepcopy](https://github.com/xieyuschen/deepcopy) -> [wangzhione/deepcopy](https://github.com/wangzhione/sbp/tree/master/util/deepcopy)

# deepcopy

用法

```Go
func Clone[T any](src T) (dst T) {
	i, err := Copy(src)
	if err != nil {
		return
	}
	dst, _ = i.(T)
	return
}
```