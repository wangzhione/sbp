package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"runtime/debug"
	"sbp/util/trace"
	"time"

	_ "github.com/go-sql-driver/mysql" // init MySQL 驱动
)

// SQLHelper 数据库帮助结构体
type SQLHelper struct {
	DB *sql.DB
}

// NewMySQLHelper 创建一个新的 MySQLHelper 实例
func NewMySQLHelper(ctx context.Context, config MySQLConfig) (helper *SQLHelper, err error) {
	// 构建 DSN（Data Source Name）
	dsn := config.DataSourceName()
	if trace.EnableLevel == slog.LevelDebug {
		slog.DebugContext(ctx, "dsn and mysql cmd", "mysql", dsn, "command", config.Command())
	}

	// 初始化数据库连接
	db, err := sql.Open("mysql", dsn)
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

	// 测试连接 默认 2s 内如果链接, 不成功, 认为失败
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		slog.ErrorContext(ctx, "failed to ping MySQL", "dsn", dsn, "reason", err)
		return
	}

	slog.InfoContext(ctx, "Connected to MySQL successfully!", "database", config.Database, "username", config.Username)
	helper = &SQLHelper{DB: db}
	return
}

// Close 关闭数据库连接
func (helper *SQLHelper) Close() error {
	return helper.DB.Close()
}

// Exec 执行无返回的 SQL 语句等 例如（INSERT, UPDATE, DELETE）
func (helper *SQLHelper) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	result, err := helper.DB.ExecContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "SQLHelper Exec error", "query", query, "args", args, "reason", err)
	}
	return result, err
}

// Query 执行查询，返回多行数据
func (helper *SQLHelper) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	rows, err := helper.DB.QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "SQLHelper Query error", "query", query, "args", args, "reason", err)
	}
	return rows, err
}

// QueryCallBack 执行查询, 内部自行通过闭包来完成参数传递和返回值获取
// callback is for rows.Next() {}
func (helper *SQLHelper) QueryCallBack(ctx context.Context, callback func(context.Context, *sql.Rows) error, query string, args ...any) error {
	rows, err := helper.DB.QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "SQLHelper QueryCallBack error", "query", query, "args", args, "reason", err)
		return err
	}
	// recover panic and error rollback
	defer func() {
		if cover := recover(); cover != nil {
			slog.ErrorContext(ctx, "SQLHelper QueryCallBack panic error", "recover", cover, "stack", debug.Stack())
		}

		newerr := rows.Close()
		if newerr != nil {
			slog.ErrorContext(ctx, "SQLHelper QueryCallBack rows.Close() panic error", "newerr", newerr)
		}
	}()

	err = callback(ctx, rows)
	if err != nil {
		slog.ErrorContext(ctx, "SQLHelper QueryCallBack callback rows error", "query", query, "args", args, "reason", err)
		return err
	}

	err = rows.Err()
	if err != nil {
		slog.ErrorContext(ctx, "SQLHelper QueryCallBack rows.Err() error", "query", query, "args", args, "reason", err)
		return err
	}

	return nil
}

/*
 QueryRow template

	err := helper.DB.QueryRowContext(ctx, query, args...).Scan(dest...)
	switch err {
	case nil:
		// success
		return nil;
	case sql.ErrNoRows: // 没有记录，返回空值
		slog.InfoContext(ctx, "SQLHelper QueryRow sql.ErrNoRows")
		return err; // or empty 业务逻辑处理
	default:
		slog.ErrorContext(ctx, "SQLHelper QueryRow error", "query", query, "args", args, "reason", err)
		return nil;
	}
*/

// BeginTransaction 开启事务
func (helper *SQLHelper) BeginTransaction(ctx context.Context, transaction func(context.Context, *sql.Tx) error) (err error) {
	// opts *sql.TxOptions 用于指定事务的隔离级别和是否为只读事务。可选参数，可以传 nil 使用 mysql 默认配置。
	tx, err := helper.DB.BeginTx(ctx, nil)
	if err != nil {
		slog.ErrorContext(ctx, "SQLHelper BeginTransaction error", "reason", err)
		return err
	}

	// recover panic and error rollback
	defer func() {
		if cover := recover(); cover != nil {
			slog.ErrorContext(ctx, "SQLHelper BeginTransaction panic error", "recover", cover, "stack", debug.Stack())

			newerr := tx.Rollback()
			if newerr != nil && newerr != sql.ErrTxDone {
				slog.ErrorContext(ctx, "SQLHelper BeginTransaction Rollback defer panic error", "newerr", newerr)
			}

			err = fmt.Errorf("transaction panic: %v", cover)
			return
		}
	}()

	err = transaction(ctx, tx)
	if err != nil {
		slog.ErrorContext(ctx, "SQLHelper BeginTransaction transaction error", "reason", err)

		newerr := tx.Rollback()
		if newerr != nil && newerr != sql.ErrTxDone {
			// 依赖人工每日接入, 追查细节
			slog.ErrorContext(ctx, "SQLHelper BeginTransaction Rollback panic error", "newerr", newerr)
		}
		return err
	}

	// transaction success commit
	err = tx.Commit()
	if err != nil && err != sql.ErrTxDone {
		slog.ErrorContext(ctx, "SQLHelper BeginTransaction Commit panic error", "newerr", err)
	}

	return nil
}
