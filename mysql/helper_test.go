package mysql

import (
	"testing"

	"sbp/util/chain"

	_ "github.com/go-sql-driver/mysql"
)

var connects = "mysql -u root -p123456 -h 127.0.0.1 -P 3306 demo"

func TestNewDB(t *testing.T) {
	connects = "mysql -u root -p123456 demo"

	config, err := ParseCommand(connects)
	if err != nil {
		t.Fatal("ParseCommand error", connects, err)
	}
	t.Log(config.DataSourceName(), config.Command())

	s, err := NewDB(chain.Background, config)
	if err != nil {
		t.Fatal("NewDB fatal", err)
	}

	if s != nil {
		t.Log("Success")
	}
}

func TestQueryRow(t *testing.T) {
	config, err := ParseCommand(connects)
	if err != nil {
		t.Fatal("ParseCommand error", connects, err)
	}
	s, err := NewDB(chain.Background, config)
	if err != nil {
		t.Fatal("NewDB fatal", err)
	}

	var count int
	err = s.QueryRow(chain.Background, "select count(*) from t_user", nil, &count)
	if err != nil {
		t.Fatal("s.QueryRow fatal", err)
	}
	t.Log("count = ", count)

	err = s.QueryRow(chain.Background, "select count(*) from t_user where id > ?", []any{6}, &count)
	if err != nil {
		t.Fatal("s.QueryRow fatal", err)
	}
	t.Log("count = ", count)
}
