package jsou

import (
	"encoding/json"
	"fmt"
	"os"
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
func Debug(obj any, prefix ...any) {
	fmt.Println()
	if len(prefix) > 0 {
		fmt.Print(prefix...)
	}
	fmt.Printf("JSON DEBUG <%T>\n", obj)
	data, err := json.MarshalIndent(obj, "", "\t")
	if err != nil {
		fmt.Printf("jsou.Debug error: %+v\n", err)
	}
	fmt.Printf("%s\n\n", data)
}
