package chain

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetExeName(t *testing.T) {
	t.Log(ExeName)

	exePath, err := os.Executable()
	t.Log(exePath, err)

	t.Log(ExeDir)
}

// SplitPath 获取路径的目录部分和文件名部分
// - filepath.Dir(path)  	-> 目录
// - filepath.Base(path) 	-> 文件名
// - filepath.Ext(filename) -> 扩展名
func SplitPath(path string) (dir, filename, ext string) {
	dir = filepath.Dir(path)
	filename = filepath.Base(path)
	ext = filepath.Ext(filename)
	return
}

func TestSplitPath(t *testing.T) {
	tests := []struct {
		input    string
		wantDir  string
		wantFile string
		wantExt  string
	}{
		{"a/b/c.txt", "a/b", "c.txt", ".txt"},
		{"a/b/c", "a/b", "c", ""},
		{"/tmp/test.log", "/tmp", "test.log", ".log"},
		{"test", ".", "test", ""},
		{"a/b/.env", "a/b", ".env", ""},
		{"", ".", "", ""},
	}

	for _, tt := range tests {
		dir, file, ext := SplitPath(tt.input)
		if dir != tt.wantDir || file != tt.wantFile || ext != tt.wantExt {
			t.Errorf("SplitPath(%q) = (%q, %q, %q); want (%q, %q, %q)",
				tt.input, dir, file, ext, tt.wantDir, tt.wantFile, tt.wantExt)
		}
	}
}
