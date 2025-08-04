package chain

import (
	"os"
	"testing"
)

func TestGetExeName(t *testing.T) {
	t.Log(ExeName)

	exePath, err := os.Executable()
	t.Log(exePath, err)

	t.Log(ExeDir)
}

func TestUUID(t *testing.T) {
	id := UUID()

	t.Logf("id = %s", id)
}
