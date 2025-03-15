package jsou

import (
	"encoding/json"
	"fmt"
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

// ReadFile 文件中读取生成 json 对象
func ReadFile[T any](filePath string) (obj T, err error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &obj)
	return
}

// WriteFile 把内容写入文件中
func WriteFile(filePath string, obj any) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil
	}

	return os.WriteFile(filePath, data, 0o644)
}

// ReadWriteFile src -> T json obj -> dst
func ReadWriteFile[T any](src, dst string) (err error) {
	// 打开源文件
	source, err := os.Open(src)
	if err != nil {
		return
	}
	defer source.Close()

	// 创建目标文件
	dest, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dest.Close()

	var obj T
	err = json.NewDecoder(source).Decode(&obj)
	if err != nil {
		return
	}

	return json.NewEncoder(dest).Encode(obj)
}

// Valid 判断字符串是否为合法 JSON
func Valid(stj string) bool {
	return json.Valid([]byte(stj))
}

// Map JSON 字符串转为 map[string]any 类似 Unmarshal[map[string]any](stj)
func Map(stj string) (obj map[string]any, err error) {
	err = json.Unmarshal([]byte(stj), &obj)
	return
}

// Slice JSON 字符串转为 []any
func Slice(stj string) (obj []any, err error) {
	err = json.Unmarshal([]byte(stj), &obj)
	return
}

// Debug json + printf 方便单元测试
func Debug(args ...any) {
	for _, arg := range args {
		fmt.Println()

		if arg == nil {
			fmt.Println("DEBUG nil\nnil")
			continue
		}

		t := reflect.TypeOf(arg)
		if t.PkgPath() != "" {
			fmt.Printf("DEBUG %s.%s\n", t.PkgPath(), t.Name())
		} else {
			fmt.Printf("DEBUG %s\n", t.Name())
		}

		// 尝试格式化 JSON
		data, err := json.MarshalIndent(arg, "", "\t")
		if err != nil {
			fmt.Printf("%#v\n", arg) // 备用输出，防止 JSON 失败时无法查看数据
		} else {
			fmt.Println(string(data))
		}
	}

	fmt.Println()
}
