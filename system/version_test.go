package system

import (
	"testing"
	"uuid"
)

func TestVersion(t *testing.T) {
	t.Log("BuildGoVersion", BuildGoVersion)
	t.Log("GitVersion", GitVersion)
	t.Log("GitLastCommitTime", GitLastCommitTime)
}

func TestUUID(t *testing.T) {
	println(uuid.NewV7().String())
	println(uuid.NewV4().String())
}
