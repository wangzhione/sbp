package jsou

import (
	"testing"
)

func TestDebug(t *testing.T) {
	s := []string{"123", "456", `789
	`, "8a\n\n\"bc"}

	Debug(s)
	Debug(nil)

	type XX struct {
		A int
		B string
	}

	Debug(XX{A: 2, B: "XX"})
}

func TestDebugPrefix(t *testing.T) {
	type XX struct {
		A int
		B string
		C func()
	}

	Debug(nil, XX{A: 2, B: "XX"}, 2, 3)
}
