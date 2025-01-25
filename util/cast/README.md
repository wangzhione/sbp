# cast

[spf13/cast](https://github.com/spf13/cast) 库非常出名, 特别是 PHP 类似运行时语言转 Golang 开发很喜欢. 

但在 Golang 开发中, 推荐不要用, 除非遇到通信 RPC 双方跨语言协议适合尝试. 

平时可以用更加安全的 Go 内置转换. 其实大部分业务转换无外乎 int 和 string 直接交互, 或者 time 业务.

常见基本都能从 cast.go 找到解法

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

# time 时间业务

```Go
str := "2025-03-09T01:28:12.946770726Z"
time, error := time.Parse(time.RFC3339Nano, str)

result, error := time.ParseInLocation(time.RFC3339Nano, str, time.UTC)

time.Now().Unix().Format("2006-01-02 15:04:05")

time = time.Unix(sec, nsec)

timeLayout := "2006-01-02 15:04:05"                 // 转化所需模板
loc, _ := time.LoadLocation("Asia/Shanghai")        // 设置上海时区
```

复杂的 RPC 业务 可以使用 `spf13/cast.ToTime` 用于转换为 `go time` 类型. 普通业务流畅和安全可控简单用原生 api 就可以. 

## 补充
 
这里 `cast` 是对 `spf13/cast` 补充, 二者可以结合使用
