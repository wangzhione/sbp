package jsou

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
)

// String 结构体转换为 JSON 字符串
func String(obj any) string {
	data, _ := json.Marshal(obj)
	return string(data)
}

// Unmarshal 将 JSON 字符串解析为结构体（泛型）
func Unmarshal[T any](stj string) (obj T, err error) {
	err = json.Unmarshal([]byte(stj), &obj)
	return
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

// ReadWriteFile
// 1. 读取 src 文件, 尝试生成 json T obj 对象;
// 2. 尝试将 obj 转成 json 格式, 然后输出到 dst destination（目的地）目标文件中;
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

// Valid 判断字符串 or []byte 是否为合法 json
func Valid[T ~string | ~[]byte](dj T) bool {
	return json.Valid([]byte(dj))
}

// Map json 字符串 or []byte 数据集转为 map[string]any 类似 Unmarshal[map[string]any](dj)
func Map[T ~string | ~[]byte](dj T) (obj map[string]any, err error) {
	err = json.Unmarshal([]byte(dj), &obj)
	return
}

// Slice json 字符串 or []byte 数据集转为 []any
func Slice[T ~string | ~[]byte](dj T) (obj []any, err error) {
	err = json.Unmarshal([]byte(dj), &obj)
	return
}

// !!! 用深拷贝时候, 认真思考下, 能否真的有必要, 如果不用这种隐式深拷贝是否也可以 !!!

// DeepCopy json 深拷贝
func DeepCopy[T any](src T) (dst T, err error) {
	data, err := json.Marshal(src)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &dst)
	return
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
