package mysql

import (
	"fmt"
	"testing"
)

func TestParseCommand(t *testing.T) {
	command := "mysql -u root -p123456 -h 127.0.0.1 -P 3306 test_db"

	config, err := ParseCommand(command)
	if err != nil {
		t.Error("Error:", err)
		return
	}

	t.Logf("Parsed Config: %+v\n", config)
	t.Log(config.DataSourceName())
}

func TestParseCommand2(t *testing.T) {
	command := "mysql -uroot -p123456 -hlocalhost -P3306 resource_ai_drama"

	config, err := ParseCommand(command)
	if err != nil {
		t.Error("Error:", err)
		return
	}

	t.Logf("Parsed Config: %+v\n", config)
	t.Log(config.DataSourceName())
}

func TestConvertDSNToCommand(t *testing.T) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/test_db?charset=utf8mb4&parseTime=true&loc=UTC"

	command, err := ConvertDSNToCommand(dsn)
	if err != nil {
		t.Error("Error: ", err)
	}

	fmt.Printf("\n%s\n\n", command)
}
