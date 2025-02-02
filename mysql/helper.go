package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"runtime/debug"
	"time"

	"sbp/util/chain"

	_ "github.com/go-sql-driver/mysql" // init MySQL 驱动
)

// 设计师有话说
// 对于底层 MySQL 处理, 非常重且重要业务, 实战开发, 只推荐 context 版本
// 这里用的是显示注入 SQL log 行为. 整体而言非常轻量级别 API 设计

// MySQLDriverName 驱动名称, 方便以后有更好的 sql hook 出现, 留一个口子
var MySQLDriverName = "mysql"

// DB 数据库帮助新结构体
type DB sql.DB

// NewDB 创建一个新的 MyMySQL 实例
func NewDB(ctx context.Context, config *MySQLConfig) (s *DB, err error) {
	// 构建 DSN（Data Source Name）
	dsn := config.DataSourceName()
	if chain.EnableLevel == slog.LevelDebug {
		slog.DebugContext(ctx, "dsn and mysql cmd", "mysql", dsn, "command", config.Command())
	}

	// 初始化数据库连接
	db, err := sql.Open(MySQLDriverName, dsn)
	if err != nil {
		slog.ErrorContext(ctx, "failed to connect to MySQL", "dsn", dsn, "reason", err)
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
		slog.ErrorContext(ctx, "failed to ping MySQL", "dsn", dsn, "reason", err)
		return
	}

	slog.InfoContext(ctx, "Connected to MySQL successfully!", "database", config.Database, "username", config.Username)
	return (*DB)(db), nil
}

// Close 关闭数据库连接
func (s *DB) Close() error {
	return (*sql.DB)(s).Close()
}

// Before hook will print the query with it's args and return the context with the timestamp
func Before(ctx context.Context, query string, args ...any) time.Time {
	begin := time.Now()
	slog.InfoContext(ctx, "MySQL before", "begin", begin, "query", query, "args", args)
	return begin
}

// After hook will get the timestamp registered on the Before hook and print the elapsed time
func After(ctx context.Context, begin time.Time) {
	end := time.Now()
	elapsed := end.Sub(begin)
	if elapsed >= time.Second {
		slog.WarnContext(ctx, "MySQL After Warn slow", "elapsed", elapsed, "end", end)
	} else {
		slog.InfoContext(ctx, "MySQL After", "elapsed", elapsed, "end", end)
	}
}

// Exec 执行无返回的 SQL 语句等 例如（INSERT, UPDATE, DELETE）
func (s *DB) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	// 主动注入日志模块
	defer After(ctx, Before(ctx, query, args))

	result, err := (*sql.DB)(s).ExecContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "MySQL Exec error", "query", query, "args", args, "reason", err)
	}
	return result, err
}

// Query 执行查询，返回多行数据
func (s *DB) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	defer After(ctx, Before(ctx, query, args))

	rows, err := (*sql.DB)(s).QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "MySQL Query error", "query", query, "args", args, "reason", err)
	}
	return rows, err
}

// QueryCallBack 执行查询, 内部自行通过闭包来完成参数传递和返回值获取
// callback is for rows.Next() {}
func (s *DB) QueryCallBack(ctx context.Context, callback func(context.Context, *sql.Rows) error, query string, args ...any) error {
	defer After(ctx, Before(ctx, query, args))

	rows, err := (*sql.DB)(s).QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "MySQL QueryCallBack error", "query", query, "args", args, "reason", err)
		return err
	}
	// recover panic and error rollback
	defer func() {
		if cover := recover(); cover != nil {
			slog.ErrorContext(ctx, "MySQL QueryCallBack panic error", "recover", cover, "stack", debug.Stack())
		}

		newerr := rows.Close()
		if newerr != nil {
			slog.ErrorContext(ctx, "MySQL QueryCallBack rows.Close() panic error", "newerr", newerr)
		}
	}()

	err = callback(ctx, rows)
	if err != nil {
		slog.ErrorContext(ctx, "MySQL QueryCallBack callback rows error", "query", query, "args", args, "reason", err)
		return err
	}

	err = rows.Err()
	if err != nil {
		slog.ErrorContext(ctx, "MySQL QueryCallBack rows.Err() error", "query", query, "args", args, "reason", err)
		return err
	}

	return nil
}

// QueryRow FindOne, args is empty 可以是 nil or []any{}
func (s *DB) QueryRow(ctx context.Context, query string, args []any, dest ...any) error {
	defer After(ctx, Before(ctx, query, args))

	err := (*sql.DB)(s).QueryRowContext(ctx, query, args...).Scan(dest...)
	switch err {
	case nil: // success
		return nil
	case sql.ErrNoRows: // 没有记录，返回空值
		slog.InfoContext(ctx, "MySQL QueryRow sql.ErrNoRows", "query", query, "args", args)
		return err // or empty 业务逻辑处理
	default:
		slog.ErrorContext(ctx, "MySQL QueryRow error", "query", query, "args", args, "reason", err)
		return err
	}
}

type Tx sql.Tx

// Transaction 开启事务
func (s *DB) Transaction(ctx context.Context, transaction func(context.Context, *Tx) error) (err error) {
	defer After(ctx, Before(ctx, "Transaction"))

	// opts *sql.TxOptions 用于指定事务的隔离级别和是否为只读事务。可选参数，可以传 nil 使用 mysql 默认配置。
	tx, err := (*sql.DB)(s).BeginTx(ctx, nil)
	if err != nil {
		slog.ErrorContext(ctx, "MySQL Transaction error", "reason", err)
		return err
	}

	// recover panic and error rollback
	defer func() {
		if cover := recover(); cover != nil {
			slog.ErrorContext(ctx, "MySQL Transaction panic error", "recover", cover, "stack", debug.Stack())

			newerr := tx.Rollback()
			if newerr != nil && newerr != sql.ErrTxDone {
				slog.ErrorContext(ctx, "MySQL Transaction Rollback defer panic error", "newerr", newerr)
			}

			err = fmt.Errorf("transaction panic: %v", cover)
			return
		}
	}()

	err = transaction(ctx, (*Tx)(tx))
	if err != nil {
		slog.ErrorContext(ctx, "MySQL Transaction transaction error", "reason", err)

		newerr := tx.Rollback()
		if newerr != nil && newerr != sql.ErrTxDone {
			// 依赖人工每日接入, 追查细节
			slog.ErrorContext(ctx, "MySQL Transaction Rollback panic error", "newerr", newerr)
		}
		return err
	}

	// transaction success commit
	err = tx.Commit()
	if err != nil && err != sql.ErrTxDone {
		slog.ErrorContext(ctx, "MySQL Transaction Commit panic error", "newerr", err)
	}

	return nil
}

// Exec 执行无返回的 SQL 语句等 例如（INSERT, UPDATE, DELETE）
func (t *Tx) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	// 主动注入日志模块
	defer After(ctx, Before(ctx, query, args))

	result, err := (*sql.Tx)(t).ExecContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "MySQL Transaction Exec error", "query", query, "args", args, "reason", err)
	}
	return result, err
}

// Query 执行查询，返回多行数据
func (t *Tx) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	defer After(ctx, Before(ctx, query, args))

	rows, err := (*sql.Tx)(t).QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "MySQL Transaction Query error", "query", query, "args", args, "reason", err)
	}
	return rows, err
}

// QueryCallBack 执行查询, 内部自行通过闭包来完成参数传递和返回值获取
// callback is for rows.Next() {}
func (t *Tx) QueryCallBack(ctx context.Context, callback func(context.Context, *sql.Rows) error, query string, args ...any) error {
	defer After(ctx, Before(ctx, query, args))

	rows, err := (*sql.Tx)(t).QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "MySQL Transaction QueryCallBack error", "query", query, "args", args, "reason", err)
		return err
	}
	// recover panic and error rollback
	defer func() {
		if cover := recover(); cover != nil {
			slog.ErrorContext(ctx, "MySQL Transaction QueryCallBack panic error", "recover", cover, "stack", debug.Stack())
		}

		newerr := rows.Close()
		if newerr != nil {
			slog.ErrorContext(ctx, "MySQL Transaction QueryCallBack rows.Close() panic error", "newerr", newerr)
		}
	}()

	err = callback(ctx, rows)
	if err != nil {
		slog.ErrorContext(ctx, "MySQL Transaction QueryCallBack callback rows error", "query", query, "args", args, "reason", err)
		return err
	}

	err = rows.Err()
	if err != nil {
		slog.ErrorContext(ctx, "MySQL Transaction QueryCallBack rows.Err() error", "query", query, "args", args, "reason", err)
		return err
	}

	return nil
}

// QueryRow FindOne
func (t *Tx) QueryRow(ctx context.Context, query string, args []any, dest ...any) error {
	defer After(ctx, Before(ctx, query, args))

	err := (*sql.Tx)(t).QueryRowContext(ctx, query, args...).Scan(dest...)
	switch err {
	case nil: // success
		return nil
	case sql.ErrNoRows: // 没有记录，返回空值
		slog.InfoContext(ctx, "MySQL Transaction QueryRow sql.ErrNoRows", "query", query, "args", args)
		return err // or empty 业务逻辑处理
	default:
		slog.ErrorContext(ctx, "MySQL Transaction QueryRow error", "query", query, "args", args, "reason", err)
		return err
	}
}
