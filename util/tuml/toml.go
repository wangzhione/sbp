// Package tuml provides utility functions for working with TOML files, including marshaling, unmarshaling, file reading, and writing.
package tuml

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

// Unmarshal 将 TOML 字符串解析为结构体（泛型）
func Unmarshal[R any, P ~string | ~[]byte](data P) (obj R, err error) {
	err = toml.Unmarshal([]byte(data), &obj)
	return
}

// ReadFile 读取 src 文件, 尝试生成 json T 对象
func ReadFile[R any](patha string) (obj R, err error) {
	data, err := os.ReadFile(patha)
	if err != nil {
		return
	}

	err = toml.Unmarshal(data, &obj)
	return
}

// WriteFile 尝试将 obj 转成 json 格式, 然后输出到 patha 目标文件中
func WriteFile(patha string, obj any) error {
	data, err := toml.Marshal(obj)
	if err != nil {
		return nil
	}

	// 所有者 (owner)	6 → rw-	可读可写
	// 所在组 (group)	6 → rw-	可读可写
	// 其他人 (others)	4 → r--	只读
	return os.WriteFile(patha, data, 0o664)
}

// Valid 判断字符串 or []byte 是否为合法 json
func Valid[P ~string | ~[]byte](data P) bool {
	var v any
	return toml.Unmarshal([]byte(data), &v) == nil
}
