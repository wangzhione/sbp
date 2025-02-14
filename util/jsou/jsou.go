package jsou

import (
	"encoding/json"
	"fmt"
)

// Valid 判断字符串是否为合法 JSON
func Valid(stj string) bool {
	return json.Valid([]byte(stj))
}

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
