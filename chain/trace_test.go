package chain

import (
	"testing"
)

func TestCopyTrace(t *testing.T) {
	newxRequestID := any(XRquestID)

	if newxRequestID == xRquestID {
		t.Log("equal") // any("X-Request-Id") == any("X-Request-Id") | type equal , value equal
	} else {
		t.Log("no equal")
	}
}
