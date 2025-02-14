package jsou

import (
	"encoding/json"
)

// Valid 判断字符串是否为合法 JSON
func Valid(stj string) bool {
	return json.Valid([]byte(stj))
}

// String 结构体转换为 JSON 字符串
func String(v any) string {
	data, _ := json.Marshal(v)
	return string(data)
}

// Unmarshal 将 JSON 字符串解析为结构体（泛型）
func Unmarshal[T any](stj string) (*T, error) {
	var obj T
	err := json.Unmarshal([]byte(stj), &obj)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

// Map JSON 字符串转为 map[string]interface{}
func Map(stj string) (result map[string]interface{}, err error) {
	err = json.Unmarshal([]byte(stj), &result)
	return
}

// Slice JSON 字符串转为 []interface{}
func Slice(stj string) (result []interface{}, err error) {
	err = json.Unmarshal([]byte(stj), &result)
	return
}
