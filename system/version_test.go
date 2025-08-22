package system

import "testing"

func TestVersion(t *testing.T) {
	t.Log("BuildVersion", BuildVersion)
	t.Log("GitVersion", GitVersion)
	t.Log("GitCommitTime", GitCommitTime)
}
