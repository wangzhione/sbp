package tuml

import (
	"context"
	"log/slog"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// Unmarshal 将 JSON 字符串解析为结构体（泛型）
func Unmarshal[T any](ctx context.Context, str string) (obj T, err error) {
	err = toml.Unmarshal([]byte(str), &obj)
	if err != nil {
		slog.ErrorContext(ctx, "toml.Unmarshal error", "error", err, "toml", str)
	}
	return
}

// ReadFile 读取 src 文件, 尝试生成 json T 对象
func ReadFile[T any](ctx context.Context, patha string) (obj T, err error) {
	data, err := os.ReadFile(patha)
	if err != nil {
		slog.ErrorContext(ctx, "os.ReadFile error", "error", err, "patha", patha)
		return
	}

	err = toml.Unmarshal(data, &obj)
	return
}

// WriteFile 尝试将 obj 转成 json 格式, 然后输出到 dst 目标文件中
func WriteFile(ctx context.Context, dst string, obj any) error {
	data, err := toml.Marshal(obj)
	if err != nil {
		slog.ErrorContext(ctx, "toml.Marshal error", "error", err, "dst", dst, "obj", obj)
		return nil
	}

	// 所有者 (owner)	6 → rw-	可读可写
	// 所在组 (group)	6 → rw-	可读可写
	// 其他人 (others)	4 → r--	只读
	err = os.WriteFile(dst, data, 0o664)
	if err != nil {
		slog.ErrorContext(ctx, "os.WriteFile 0o644 error", "error", err, "dst", dst)
	}
	return err
}
