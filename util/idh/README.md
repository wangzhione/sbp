# id hash 

id 和 hash 相关操作合集

## UUID

`idh.UUID()` 借助 `github.com/google/uuid` v4 random 算法, 默认返回不带 '`-`' 风格的小写串

## hash

see `md5.go` or `fenv/fnv.go`

`idh.MD5String()` 返回 md5 sign or `idh.HashString()` 返回 fnv hash 算法生成的 uint64

推荐直接使用, Hash 和 HashString 对 xxhash 一种包装器. xxhash 是 2012 年后行业主流 hash 算法, 优势是快
