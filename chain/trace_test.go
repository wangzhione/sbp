package chain

import (
	"testing"
)

func TestCopyTrace(t *testing.T) {
	if any("X-Request-Id") == any(XRquestID) {
		t.Log("equal") // any("X-Request-Id") == any("X-Request-Id") | type equal , value equal
	} else {
		t.Log("no equal")
	}
}
