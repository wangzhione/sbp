package jsou

import (
	"testing"
)

func TestDEBUG(t *testing.T) {
	s := []string{"123", "456", `789
	`, "8a\n\n\"bc"}

	DEBUG(s)
	DEBUG(nil)

	type XX struct {
		A int
		B string
	}

	DEBUG(XX{A: 2, B: "XX"})
}

func TestDEBUGPrefix(t *testing.T) {
	type XX struct {
		A int
		B string
		C func()
	}

	DEBUG(nil, XX{A: 2, B: "XX"}, 2, 3)
}
