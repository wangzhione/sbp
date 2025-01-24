# cast

[spf13/cast](https://github.com/spf13/cast) 库非常出名, 特别是 PHP 类似运行时语言转 Golang 开发很喜欢. 

但在 Golang 开发中, 推荐不要用, 除非遇到通信 RPC 双方跨语言协议适合尝试. 

平时可以用更加安全的 Go 内置转换. 其实大部分业务转换无外乎 int 和 string 直接交互

## Go 常见转换

**string 转成 int：**

```Go
int, err := strconv.Atoi(string)
```

**string 转成 int64：**

```Go
int64, err := strconv.ParseInt(string, 10, 64)
```

**int 转成 string：**

```Go
string := strconv.Itoa(int)
```

**int64 转成 string：**

```Go
string := strconv.FormatInt(int64, 10)
```

方便也可以使用本库泛型版本 `cast.IntToString` 和 `cast.StringToInt` 具体看实现

## 补充
 
这里 `cast` 是对 `spf13/cast` 补充, 二者可以结合使用
