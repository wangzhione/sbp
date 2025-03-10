package sqlite

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	_ "github.com/mattn/go-sqlite3" // init SQLite 驱动
	"github.com/wangzhione/sbp/sqler"
)

// NewDB 创建 sqlite3 实例
func NewDB(ctx context.Context, command string) (s *sqler.DB, err error) {
	db, err := sql.Open("sqlite3", command)
	if err != nil {
		slog.ErrorContext(ctx, "failed to connect to SQLite error", "command", command, "error", err)
		return
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// Ping 运行一个简单的 SQL 查询，确保数据库可用
	var test int
	err = db.QueryRowContext(ctx, "SELECT 1").Scan(&test)
	if err != nil {
		slog.ErrorContext(ctx, "failed to ping SQLite error", "command", command, "error", err)
		return
	}

	slog.InfoContext(ctx, "Connected to SQLite successfully!", "command", command)
	s = (*sqler.DB)(db)
	return
}
