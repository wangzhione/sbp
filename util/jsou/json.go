// Package jsou provides utility functions for working with JSON data, including marshaling, unmarshaling, file operations, and debugging helpers.
package jsou

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
)

// 这个库继承自官方 "encoding/json", 特殊情况存在 panic, 依赖程序最外层去捕获 recover()。
//
// 注意：在以下特殊情况下可能触发 panic（非返回 error）：
// 1. Marshal() 参数类型不受支持：
//    例如 func、chan、complex、unsafe.Pointer 等类型，内部反射找不到编码器会直接 panic。
// 2. Unmarshal() 目标不是指针或为 nil 接口值：
//    Unmarshal() 必须接收“可写入”的目标，否则会直接 panic。
//    常见错误：
//       var v map[string]interface{}
//       json.Unmarshal(data, v)        // ❌ panic: Unmarshal(non-pointer map[string]interface {})
//    因为 v 只是一个值（非指针），函数内部无法修改它的内容。
//    正确写法应为：
//       json.Unmarshal(data, &v)       // ✅ 传入指针，允许修改内容
//    同理，如果目标是一个 nil 接口（例如 var v interface{}，未取地址），
//    也会 panic，因为无法动态设置其值，需传入 &v。
// 3. struct tag 非法（如 json:"a,b,c"）：
//    标准库解析 tag 失败会 panic，应确保 tag 符合 "name[,option]" 形式。
// 4. 结构体循环引用导致栈溢出：
//    结构体字段指向自身（或间接循环）时递归展开无限循环，最终触发 stack overflow。
// 5. 自定义 MarshalJSON / UnmarshalJSON 方法内部 panic：
//    若用户自定义序列化逻辑中主动 panic（或访问空指针）会向上传递。
// 6. 输入数据包含非法 UTF-8 字节或 Reader 实现 panic：
//    在字符串转换或流读取时可能触发 panic（常见于损坏或非文本输入）。

// 相关协议部分阅读 https://www.json.org/json-zh.html

// String 结构体转换为 JSON 字符串
func String(obj any) string {
	data, _ := json.Marshal(obj)
	return string(data)
}

// Unmarshal 将 JSON 字符串解析为结构体（泛型）
func Unmarshal[R any, P ~string | ~[]byte](data P) (obj R, err error) {
	err = json.Unmarshal([]byte(data), &obj)
	return
}

// To 将一个类型的值转换为(另)一个类型的值（泛型）; 类似 Simple DeepCopy
func To[R any](a any) (b R, err error) {
	data, err := json.Marshal(a)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &b)
	return
}

// Map json 字符串 数据集转为 map[string]any 类似 Unmarshal[map[string]any](data)
func Map[P ~string | ~[]byte](data P) (obj map[string]any, err error) {
	err = json.Unmarshal([]byte(data), &obj)
	return
}

// Array json 字符串 数据集转为 []any 类似 Unmarshal[[]any](data)
func Array[P ~string | ~[]byte](data P) (obj []any, err error) {
	err = json.Unmarshal([]byte(data), &obj)
	return
}

// WriteFile 尝试将 obj 转成 json 格式, 然后输出到 dst 目标文件中
func WriteFile(patha string, obj any) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil
	}

	// 所有者 (owner)	6 → rw-	可读可写
	// 所在组 (group)	6 → rw-	可读可写
	// 其他人 (others)	4 → r--	只读
	return os.WriteFile(patha, data, 0o664)
}

// ReadFile 读取 src 文件, 尝试生成 json T 对象
func ReadFile[R any](src string) (obj R, err error) {
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
func Valid[P ~string | ~[]byte](data P) bool {
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
