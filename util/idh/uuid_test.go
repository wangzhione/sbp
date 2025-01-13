package idh

import "testing"

func TestUUID(t *testing.T) {
	id := UUID()

	t.Logf("id = %s", id)
}
