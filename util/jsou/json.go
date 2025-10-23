// Package jsou provides utility functions for working with JSON data, including marshaling, unmarshaling, file operations, and debugging helpers.
package jsou

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
)

// 这个库继承自 "encoding/json", 特殊情况存在 panic, 依赖程序最外层去捕获

// String 结构体转换为 JSON 字符串
func String(obj any) string {
	data, _ := json.Marshal(obj)
	return string(data)
}

// Unmarshal 将 JSON 字符串解析为结构体（泛型）
func Unmarshal[T any](data string) (obj T, err error) {
	err = json.Unmarshal([]byte(data), &obj)
	return
}

// To 将一个类型的值转换为(另)一个类型的值（泛型）; 类似 Simple DeepCopy
func To[T any](src any) (dst T, err error) {
	data, err := json.Marshal(src)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &dst)
	return
}

// WriteFile 尝试将 obj 转成 json 格式, 然后输出到 dst 目标文件中
func WriteFile(dst string, obj any) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil
	}

	// 所有者 (owner)	6 → rw-	可读可写
	// 所在组 (group)	6 → rw-	可读可写
	// 其他人 (others)	4 → r--	只读
	return os.WriteFile(dst, data, 0o664)
}

// ReadFile 读取 src 文件, 尝试生成 json T 对象
func ReadFile[T any](src string) (obj T, err error) {
	data, err := os.ReadFile(src)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &obj)
	return
}

// ReadWriteFile 1. 读取 src 文件, 尝试生成 json T obj 对象; 2. 尝试将 obj 转成 json 格式, 然后输出到 dst destination（目的地）目标文件中;
func ReadWriteFile[T any](src, dst string) (err error) {
	// 打开源文件
	source, err := os.Open(src)
	if err != nil {
		return
	}
	defer source.Close()

	var obj T
	err = json.NewDecoder(source).Decode(&obj)
	if err != nil {
		return
	}

	return WriteFile(dst, obj)
}

// Valid 判断字符串 是否为合法 json
// 当你想要用 []byte 当成参数时候, 默认你是有一定选择能力开放人员, 这时候可以自行选定 json.Valid ...
func Valid(data string) bool {
	return json.Valid([]byte(data))
}

// DEBUG json + fmt printf 简单打印测试, args[0] 可以传入 io.Writer, 不传入默认 os.Stdout
func DEBUG(args ...any) {
	if len(args) == 0 {
		return
	}

	i := 0
	var w io.Writer = os.Stdout
	if writer, ok := args[i].(io.Writer); ok {
		w = writer
		i++
	}
	for ; i < len(args); i++ {
		fmt.Fprintln(w, "")

		arg := args[i]
		if arg == nil {
			fmt.Fprintln(w, "DEBUG nil\nnil")
			continue
		}

		t := reflect.TypeOf(arg)
		if t.PkgPath() != "" {
			fmt.Fprintf(w, "DEBUG %s.%s\n", t.PkgPath(), t.Name())
		} else {
			fmt.Fprintf(w, "DEBUG %s\n", t.Name())
		}

		// 尝试格式化 JSON
		data, err := json.MarshalIndent(arg, "", "\t")
		if err != nil {
			fmt.Fprintf(w, "%+v\n", arg) // 备用输出，防止 JSON 失败时无法查看数据
		} else {
			fmt.Fprintln(w, string(data))
		}
	}

	fmt.Fprintln(w, "")
}
