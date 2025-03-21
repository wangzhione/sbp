# id hash 

id 和 hash 相关操作合集

## hash

see `md5.go` or `fenv/fnv.go`

`idhash.MD5()` 返回 md5 sign or `fnv.HashString()` 返回 fnv hash 算法生成的 uint64

推荐直接使用 or 优先使用, Hash 和 HashString 对 xxhash 一种包装器. xxhash 是 2012 年后行业主流 hash 算法, 优势是快
