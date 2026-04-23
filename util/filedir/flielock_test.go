package filedir

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTryFileLock(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	path := filepath.Join(dir, "demo.txt")

	lock1, err := TryFileLock(ctx, path)
	if err != nil {
		t.Fatal(err)
	}
	defer lock1.Unlock(ctx)

	_, err = TryFileLock(ctx, path)
	if !errors.Is(err, ErrFileLockBusy) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWithFileLock(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	path := filepath.Join(dir, "demo.txt")

	start := make(chan struct{})
	locked := make(chan struct{})
	done := make(chan struct{})

	go func() {
		err := WithFileLock(ctx, path, func() error {
			close(start)
			time.Sleep(50 * time.Millisecond)
			return nil
		})
		if err != nil {
			t.Error(err)
		}
		close(done)
	}()

	<-start

	go func() {
		err := WithFileLock(ctx, path, func() error {
			close(locked)
			return nil
		})
		if err != nil {
			t.Error(err)
		}
	}()

	select {
	case <-locked:
		t.Fatal("lock should block concurrent caller")
	case <-time.After(20 * time.Millisecond):
	}

	<-done

	select {
	case <-locked:
	case <-time.After(time.Second):
		t.Fatal("lock not released")
	}

	time.Sleep(20 * time.Millisecond)

	_, err := os.Stat(path + ".lock")
	if !os.IsNotExist(err) {
		t.Fatalf("lock file should be removed, err=%v", err)
	}
}
