# cache

Google 的 Guava Cache Java 库非常好用, 但在 Go 上几乎没有看见这样的库. 

这里希望有个支持 LRU + TTL + Max Limit 实战项目用的 simple beautiful cache 库, 也在考虑和寻找相关成熟库

目前选中 [ristretto](https://github.com/dgraph-io/ristretto) 库, 这个像个 Go 库

## ristretto example

**Installing**

```Shell
go get github.com/dgraph-io/ristretto/v2@latest
```

**example**

```Go
cache, err := ristretto.NewCache(&ristretto.Config[string, string]{
    NumCounters: 1e7,     // 用于统计频率
    MaxCost:     1 << 30, // 最大内存成本为 1 GB
    BufferItems: 64,      // 写入缓冲区大小
})
if err != nil {
    panic(err)
}

// 假设我们缓存一些字符串
str := "This is some data"
cache.Set("key", str, int64(len(str))) // 设置成本为字符串长度
cache.Wait() // 确保写入成功

value, found := cache.Get("key")
if found {
    fmt.Println("Cache hit:", value)
}
```

**NumCounters**

1 **主要作用：**

- NumCounters 是用来维护缓存中条目访问频率的计数器数量。
- 它直接影响驱逐（eviction）的准确性和缓存的命中率（hit ratio）。

2 **推荐设置：**

- 如果你的缓存预期能存储 1,000,000 条目，那么 NumCounters 应该设置为 10,000,000（即 10 倍）。
- 这种通过经验设置的值能够提升频率统计的精度，从而改进缓存驱逐策略。

3 **内存使用：**

- 每个计数器大约需要 3 字节（4 位计数器 × 4 副本 + 大约 1 字节的布隆过滤器）。
- NumCounters 的值会被内部自动调整为 最近的 2 的幂，因此实际内存开销可能略高于公式计算。

**BufferItems**

1 **作用：**

- BufferItems 决定了 Get 操作的缓冲区大小。
- 它是为了缓解高并发访问时的争用（contention）。

2 **默认值：**

- 默认设置为 64，在绝大多数场景下，这个值已经能提供良好的性能。

3 **调整建议：**

- 如果在高并发场景下发现 Get 性能下降，可以尝试将该值按 64 的增量逐步增大。
- 但大多数情况下，你不需要更改这个值。

***

多看代码注释, 优先知道怎么用, 怎么用的安全用的漂亮. 不要迷在源码里面走不出来.
