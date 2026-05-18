package system

import "testing"

func TestVersion(t *testing.T) {
	t.Log("BuildGoVersion", BuildGoVersion)
	t.Log("GitVersion", GitVersion)
	t.Log("GitLastCommitTime", GitLastCommitTime)
}
