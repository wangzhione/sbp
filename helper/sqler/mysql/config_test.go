package mysql

import (
	"fmt"
	"strings"
	"testing"

	mysqlDriver "github.com/go-sql-driver/mysql"
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

func TestDataSourceNameEscape(t *testing.T) {
	config := &MySQLConfig{
		Username: "root",
		Password: "123456",
		Host:     "127.0.0.1",
		Port:     3306,
		Database: "a/b",
		Location: "Asia/Shanghai",
	}

	dsn := config.DataSourceName()
	parsed, err := mysqlDriver.ParseDSN(dsn)
	if err != nil {
		t.Fatalf("DataSourceName should build valid DSN, dsn=%s, err=%v", dsn, err)
	}
	if parsed.DBName != config.Database {
		t.Fatalf("database = %q, want %q, dsn=%s", parsed.DBName, config.Database, dsn)
	}
	if parsed.Loc == nil || parsed.Loc.String() != config.Location {
		t.Fatalf("loc = %v, want %q, dsn=%s", parsed.Loc, config.Location, dsn)
	}
}

func TestConvertDSNToCommandIPv6(t *testing.T) {
	command, err := ConvertDSNToCommand("root:123456@tcp([de:ad:be:ef::ca:fe]:80)/test_db?charset=utf8mb4")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(command, "-h de:ad:be:ef::ca:fe -P 80") {
		t.Fatalf("command should keep IPv6 host and port, got %q", command)
	}
}
