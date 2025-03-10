package sqlite

import (
	"testing"

	"github.com/wangzhione/sbp/util/chain"
)

func TestNewDB(t *testing.T) {
	ctx := chain.Context()

	command := "./test.db"

	s, err := NewDB(ctx, command)
	if err != nil {
		t.Fatal("NewDB fatal", err)
	}
	_ = s
}
