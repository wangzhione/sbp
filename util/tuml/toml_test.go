package tuml

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFileMarshalError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.toml")

	// func 无法编码为 TOML, WriteFile 必须把 toml.Marshal 的错误返回给调用方.
	err := WriteFile(path, map[string]any{"bad": func() {}})
	if err == nil {
		t.Fatal("WriteFile should return toml.Marshal error")
	}

	if _, statErr := os.Stat(path); !os.IsNotExist(statErr) {
		t.Fatalf("WriteFile should not create file on marshal error, statErr=%v", statErr)
	}
}
