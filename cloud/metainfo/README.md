
metainfo
========

该库提供了一种在 go 语言的 `context.Context` 中保存用于跨服务传递的元信息的统一接口。

**元信息**被设计为字符串的键值对，并且键是**大小写敏感的**。

数据传递会在整个服务调用链上一直传递，直到被丢弃。

该库被设计成针对 `context.Context` 进行操作的接口集合，而具体的元数据在网络上传输的形式和方式，应由支持该库的框架来实现。通常，终端用户不应该关注其具体传输形式，而应该仅依赖该库提供的抽象接口。

框架支持指南
------------

如果要在某个框架里引入 metainfo 并对框架的用户提供支持，需要满足如下的条件：

1. 框架使用的传输协议应该支持元信息的传递（例如 HTTP header、thrift 的 header transport 等）。
2. 当框架作为服务端接收到元信息后，需要将元信息添加到 `context.Context` 对象里。随后，进入用户的代码逻辑


API 参考
-------

**注意**

1. 出于兼容性和普适性，元信息的形式为字符串的 key value 对。
2. 空串作为 key 或者 value 都是无效的。
3. 由于 context 的特性，程序对 metainfo 的增删改只会对拥有相同的 context 或者其子 context 的代码可见。

**常量**

metainfo 包提供了几个常量字符串前缀，用于无 context（例如网络传输）的场景下标记元信息的类型。

典型的业务代码通常不需要用到这些前缀。支持 metainfo 的框架也可以自行选择在传输时用于区分元信息类别的方式。

- `PrefixPersistent`

**方法**

- `TransferForward(ctx context.Context) context.Context`
    -  persistent 数据等传递
- `GetPersistentValue(ctx context.Context, k string) (string, bool)`
    - 从 context 里获取指定 key 的 persistent 数据。
- `GetAllPersistentValues(ctx context.Context) map[string]string`
    - 从 context 里获取所有 persistent 数据。
- `RangePersistentValues(ctx context.Context, f func(k, v string) bool)`
    - 从 context 里基于 f 过滤获取 persistent 数据。
- `WithPersistentValue(ctx context.Context, k string, v string) context.Context`
    - 向 context 里添加一个 persistent 数据。
- `DelPersistentValue(ctx context.Context, k string) context.Context`
    - 从 context 里删除指定的 persistent 数据。
