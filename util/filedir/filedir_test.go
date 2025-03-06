package filedir

import (
	"runtime"
	"testing"
	"time"
)

func TestOpenFile(t *testing.T) {
	file, err := OpenFile("filedir_test.go")
	if err != nil {
		t.Fatal("OpenFile fatal", err)
	}

	// 模拟使用 fileResource
	_ = file
	t.Log("File opened successfully")

	// 当 file 不再被引用时，AddCleanup 会自动关闭文件
	file = nil

	// 强制触发 GC，以便清理 fileResource
	runtime.GC()

	time.Sleep(time.Second * 6)
}

func TestFileList(t *testing.T) {
	dirname := "../"

	files, err := FileList(dirname)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(files)
}
