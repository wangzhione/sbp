package mysql

import (
	"context"
	"database/sql"
	"log/slog"
	"sbp/util/trace"

	_ "github.com/go-sql-driver/mysql" // 导入 MySQL 驱动
)

// NewMySQLHelper 创建一个新的 MySQLHelper 实例
func NewMySQLHelper(ctx context.Context, config MySQLConfig) (db *sql.DB, err error) {
	// 构建 DSN（Data Source Name）
	dsn := config.DataSourceName()
	if trace.EnableLevel == slog.LevelDebug {
		slog.DebugContext(ctx, "dsn and mysql cmd", "mysql", dsn, "command", config.Command())
	}

	// 初始化数据库连接
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		slog.ErrorContext(ctx, "failed to connect to MySQL", "dsn", dsn, "reason", err)
		return
	}

	// 配置连接池
	if config.MaxIdleConns != nil {
		db.SetMaxOpenConns(*config.MaxOpenConns)
	} else {
		// 默认有 3 个空闲连接, 库本身默认 2 个, 这边多一个, 尝试用于 goroutine chan Exec
		db.SetMaxOpenConns(3)
	}
	if config.MaxIdleConns != nil {
		db.SetMaxIdleConns(*config.MaxIdleConns)
	} else {
		// 单个连接消耗系统资源接近 4MB, 256 个连接差不多 1G, 而且这只是单台机器. 高并发请求 MySQL 本身存在瓶颈
		db.SetMaxIdleConns(128)
	}

	// 测试连接
	if err = db.Ping(); err != nil {
		slog.ErrorContext(ctx, "failed to ping MySQL", "dsn", dsn, "reason", err)
		return
	}

	slog.InfoContext(ctx, "Connected to MySQL successfully!", "database", config.Database, "username", config.Username)
	return
}
