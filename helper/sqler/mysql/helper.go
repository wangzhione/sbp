package mysql

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	_ "github.com/go-sql-driver/mysql" // init MySQL 驱动
	"github.com/wangzhione/sbp/chain"
	"github.com/wangzhione/sbp/helper/sqler"
)

// 设计师有话说
// 对于底层 MySQL 处理, 非常重且重要业务, 实战开发, 只推荐 context 版本
// 这里用的是显示注入 SQL log 行为. 整体而言非常轻量级别 API 设计

// MySQLDriverName 驱动名称, 方便以后有更好的 sql hook 出现, 留一个口子
var MySQLDriverName = "mysql"

// NewDBWithConfig 创建一个新的 MyMySQL 实例, 需要自行 Close 释放资源
func NewDBWithConfig(ctx context.Context, config *MySQLConfig) (s *sqler.DB, err error) {
	// 构建 DSN（Data Source Name）
	dsn := config.DataSourceName()
	if chain.EnableLevel == slog.LevelDebug {
		slog.DebugContext(ctx, "dsn and mysql cmd", "mysql", dsn, "command", config.Command())
	}

	// 初始化数据库连接
	db, err := sql.Open(MySQLDriverName, dsn)
	if err != nil {
		slog.ErrorContext(ctx, "failed to connect to MySQL error", "dsn", dsn, "error", err, "cmd", config.Command())
		return
	}

	// 配置连接池
	if config.MaxOpenConns != nil {
		db.SetMaxOpenConns(*config.MaxOpenConns)
	} else {
		// 单个连接消耗系统资源接近 4MB, 256 个连接差不多 1G, 而且这只是单台机器. 高并发请求 MySQL 本身存在瓶颈
		// 当然如果真的跑满了 128 , 哪怕此刻设置为 0, 在横向服务器组中, 对 MySQL 压力巨大的, 可能会让其拒绝服务(处理不过来)
		db.SetMaxOpenConns(128)
	}
	if config.MaxIdleConns != nil {
		db.SetMaxIdleConns(*config.MaxIdleConns)
	}

	// 测试连接 默认 2s 内如果链接, 不成功, 认为失败
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		slog.ErrorContext(ctx, "failed to ping MySQL panic error", "dsn", dsn, "error", err, "cmd", config.Command())
		return
	}

	slog.InfoContext(ctx, "Connected to MySQL successfully", "database", config.Database, "username", config.Username)
	s = (*sqler.DB)(db)
	return
}

func NewDB(ctx context.Context, command string) (s *sqler.DB, err error) {
	config, err := ParseCommand(command)
	if err != nil {
		slog.ErrorContext(ctx, "ParseCommand error", "error", err, "command", command)
		return
	}

	return NewDBWithConfig(ctx, config)
}
