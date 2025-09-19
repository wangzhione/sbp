package chain

import (
	"os"
	"path/filepath"
	"sync"
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

func TestUUID(t *testing.T) {
	id := UUID()
	t.Logf("id = %s", id) // id = 22ba3cffc8de4a2d9dc8a95d09ed03e1

	// import "github.com/google/uuid"
	// go mod tidy
	//
	// id := uuid.New().String()
	// t.Logf("id = %s", id) // id = 22ba3cff-c8de-4a2d-9dc8-a95d09ed03e1
}

func BenchmarkUUID(b *testing.B) {
	for b.Loop() {
		_ = UUID()
	}
}

// 测试高并发下的性能表现
func BenchmarkUUIDHighConcurrency(b *testing.B) {
	b.Run("HighConcurrency_1000", func(b *testing.B) {
		var wg sync.WaitGroup
		concurrency := 1000

		b.ResetTimer()

		for range concurrency {
			wg.Go(func() {
				iterations := b.N / concurrency
				if iterations == 0 {
					iterations = 1
				}
				for j := 0; j < iterations; j++ {
					_ = UUID()
				}
			})
		}

		wg.Wait()
	})
}
