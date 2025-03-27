# tasks

tasks is a simple goroutine pool which aims to reuse goroutines and limit the number of goroutines.

> tasks 是一个简单的 goroutine 池，旨在复用 goroutine，并**限制 goroutine 的数量**。

## example

**[optional] Step 0 : main.init update tasks.PanicHandler** 

```Go
// first register global panic handler
tasks.PanicHandler = func (ctx context.Context, cover any) {
    // ctx is tasks.Go func context, cover = recover()
}
```

**Step 1 : Let's Go**

```Go
// package: new pool 
var o = tasks.NewPool(8)

// func: run Go
// ctx 参照 chain.CopyTrace 去脱离 context cancel or 脱敏
// https://github.com/wangzhione/sbp/blob/master/chain/trace.go#L30-L44
o.Go(ctx, func(ctx context.Context) {
    // Your business
})
```
