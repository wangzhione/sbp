package filedir

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestFSyncWriteReader(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "atomic.txt")
	data := []byte("atomic write by reader")

	err := FSyncWriteReader(path, bytes.NewReader(data), 0o664)
	if err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != string(data) {
		t.Fatalf("unexpected file content: %q", string(got))
	}
}

func TestFSyncWriteReaderNil(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "atomic.txt")

	err := FSyncWriteReader(path, nil, 0o664)
	if !errors.Is(err, os.ErrInvalid) {
		t.Fatalf("unexpected error: %v", err)
	}
}
