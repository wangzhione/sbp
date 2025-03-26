package filedir

import (
	"runtime"
	"testing"
	"time"

	"github.com/wangzhione/sbp/chain"
	"github.com/wangzhione/sbp/util/jsou"
)

var ctx = chain.Context()

func TestOpenFile(t *testing.T) {
	file, err := OpenFile(ctx, "filedir_test.go")
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

	files, err := FileList(ctx, dirname)
	if err != nil {
		t.Fatal(err)
	}

	jsou.DEBUG(files)
}

func TestCreateDir(t *testing.T) {
	path := `E:\github.com\wangzhione\sbp\util\filedir\aa\bb\cc\filedir_test.go`

	err := CreateDir(ctx, path)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Succes")
}
