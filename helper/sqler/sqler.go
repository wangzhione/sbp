package sqler

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"runtime/debug"
	"time"
)

// DB 数据库帮助新结构体, 也可以 (*sql.DB)(s) 调用原生接口
type DB sql.DB

func (s *DB) DB() *sql.DB {
	return (*sql.DB)(s)
}

// Close 关闭数据库连接, 必须主动去执行, 否则无法被回收
func (s *DB) Close(ctx context.Context) (err error) {
	if s != nil {
		err = s.DB().Close()
		// 创建和关闭都是很重的操作需要格外小心
		slog.InfoContext(ctx, "r.DB().Close() info", "error", err)
	}
	return
}

// Before hook will print the query with it's args and return the context with the timestamp
func Before(ctx context.Context, query string, args ...any) time.Time {
	begin := time.Now()
	slog.InfoContext(ctx, "SQLer Before", "begin", begin, "query", query, "args", args)
	return begin
}

// After hook will get the timestamp registered on the Before hook and print the elapsed time
func After(ctx context.Context, begin time.Time) {
	end := time.Now()
	elapsed := end.Sub(begin)
	if elapsed >= time.Second {
		slog.WarnContext(ctx, "SQLer After Warn slow", "elapsed", elapsed, "end", end)
	} else {
		slog.InfoContext(ctx, "SQLer After", "elapsed", elapsed, "end", end)
	}
}

// Exec 执行无返回的 SQL 语句等 例如（INSERT, UPDATE, DELETE）
func (s *DB) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	// 主动注入日志模块
	defer After(ctx, Before(ctx, query, args))

	result, err := s.DB().ExecContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "SQLer Exec error", "query", query, "args", args, "error", err)
	}
	return result, err
}

// QueryCallBack 执行查询, 内部自行通过闭包来完成参数传递和返回值获取
// callback is for rows.Next() { if err := rows.Scan(&, &, &, ...); err != nil { } }
func (s *DB) QueryCallBack(ctx context.Context, callback func(context.Context, *sql.Rows) error, query string, args ...any) error {
	defer After(ctx, Before(ctx, query, args))

	rows, err := s.DB().QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "SQLer QueryCallBack error", "query", query, "args", args, "error", err)
		return err
	}
	// recover panic and error rollback
	defer func() {
		if cover := recover(); cover != nil {
			slog.ErrorContext(ctx, "SQLer QueryCallBack panic error", "recover", cover, "stack", string(debug.Stack()))
		}

		newerr := rows.Close()
		if newerr != nil {
			slog.ErrorContext(ctx, "SQLer QueryCallBack rows.Close() panic error", "newerr", newerr)
		}
	}()

	err = callback(ctx, rows)
	if err != nil {
		slog.ErrorContext(ctx, "SQLer QueryCallBack callback rows error", "query", query, "args", args, "error", err)
		return err
	}

	err = rows.Err()
	if err != nil {
		slog.ErrorContext(ctx, "SQLer QueryCallBack rows.Err() error", "query", query, "args", args, "error", err)
		return err
	}

	return nil
}

// QueryRow FindOne, args is empty 可以是 nil or []any{}
func (s *DB) QueryRow(ctx context.Context, query string, args []any, dest ...any) error {
	defer After(ctx, Before(ctx, query, args))

	err := s.DB().QueryRowContext(ctx, query, args...).Scan(dest...)
	switch err {
	case nil: // success
		return nil
	case sql.ErrNoRows: // 没有记录，返回空值
		slog.InfoContext(ctx, "SQLer QueryRow sql.ErrNoRows", "query", query, "args", args)
		return err // or empty 业务逻辑处理
	default:
		slog.ErrorContext(ctx, "SQLer QueryRow error", "query", query, "args", args, "error", err)
		return err
	}
}

// QueryOne 查询单条记录
func (s *DB) QueryOne(ctx context.Context, query string, args ...any) (result map[string]any, err error) {
	defer After(ctx, Before(ctx, query, args))

	rows, err := s.DB().QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "SQLer QueryOne QueryContext error", "query", query, "args", args, "error", err)
		return
	}
	defer rows.Close()

	// 获取列名
	columns, err := rows.Columns()
	if err != nil {
		slog.ErrorContext(ctx, "SQLer QueryOne rows.Columns() error", "query", query, "args", args, "error", err)
		return
	}

	result = make(map[string]any)
	if len(columns) == 0 {
		return
	}

	// 创建切片用于存储结果
	values := make([]any, len(columns))
	valuePtrs := make([]any, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	if rows.Next() {
		// 读取数据
		if err = rows.Scan(valuePtrs...); err != nil {
			slog.ErrorContext(ctx, "SQLer QueryOne rows.Scan(valuePtrs...) error", "query", query, "args", args, "error", err)
			return
		}
	}

	if err = rows.Err(); err != nil {
		slog.ErrorContext(ctx, "SQLer QueryOne rows.Err() error", "query", query, "args", args, "error", err)
		return
	}

	// 解析数据
	for i, colName := range columns {
		switch val := values[i].(type) {
		case nil:
			result[colName] = nil // 保持 NULL 值
		case []byte:
			result[colName] = string(val) // 转换 []byte 为 string
		default:
			result[colName] = val
		}
	}

	return
}

// QueryAll 查询多条记录
func (s *DB) QueryAll(ctx context.Context, query string, args ...any) (results []map[string]any, err error) {
	defer After(ctx, Before(ctx, query, args))

	rows, err := s.DB().QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "SQLer QueryAll QueryContext error", "query", query, "args", args, "error", err)
		return
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		slog.ErrorContext(ctx, "SQLer QueryAll rows.Columns() error", "query", query, "args", args, "error", err)
		return
	}

	if len(columns) == 0 {
		return
	}

	for rows.Next() {
		// 创建存储每列数据的切片
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i] // 指针绑定
		}

		// 读取数据
		if err = rows.Scan(valuePtrs...); err != nil {
			slog.ErrorContext(ctx, "SQLer QueryAll rows.Scan(valuePtrs...) error", "query", query, "args", args, "error", err)
			return
		}

		// 解析数据，转换 NULL 值
		result := make(map[string]any, len(columns))
		for i, colName := range columns {
			switch val := values[i].(type) {
			case nil:
				result[colName] = nil // 保持 NULL 值
			case []byte:
				result[colName] = string(val) // 转换 []byte 为 string
			default:
				result[colName] = val
			}
		}

		results = append(results, result)
	}

	// 检查迭代过程中是否出错
	if err = rows.Err(); err != nil {
		slog.ErrorContext(ctx, "SQLer QueryAll rows.Err() error", "query", query, "args", args, "error", err)
		return
	}

	return
}

// Transaction 开启事务
func (s *DB) Transaction(ctx context.Context, transaction func(context.Context, *sql.Tx) error) (err error) {
	defer After(ctx, Before(ctx, "Transaction"))

	// opts *sql.TxOptions 用于指定事务的隔离级别和是否为只读事务。可选参数，可以传 nil 使用 mysql 默认配置。
	tx, err := s.DB().BeginTx(ctx, nil)
	if err != nil {
		slog.ErrorContext(ctx, "SQLer Transaction error", "error", err)
		return err
	}

	// recover panic and error rollback
	defer func() {
		if cover := recover(); cover != nil {
			slog.ErrorContext(ctx, "SQLer Transaction panic error", "recover", cover, "stack", string(debug.Stack()))

			newerr := tx.Rollback()
			if newerr != nil && newerr != sql.ErrTxDone {
				slog.ErrorContext(ctx, "SQLer Transaction Rollback defer panic error", "newerr", newerr)
			}

			err = fmt.Errorf("transaction panic: %v", cover)
			return
		}
	}()

	err = transaction(ctx, tx)
	if err != nil {
		slog.ErrorContext(ctx, "SQLer Transaction transaction error", "error", err)

		newerr := tx.Rollback()
		if newerr != nil && newerr != sql.ErrTxDone {
			// 依赖人工每日接入, 追查细节
			slog.ErrorContext(ctx, "SQLer Transaction Rollback panic error", "newerr", newerr)
		}
		return err
	}

	// transaction success commit
	err = tx.Commit()
	if err != nil && err != sql.ErrTxDone {
		slog.ErrorContext(ctx, "SQLer Transaction Commit panic error", "newerr", err)
	}

	return nil
}
