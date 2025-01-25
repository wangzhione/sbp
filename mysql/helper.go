package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/url"
	"sbp/util/trace"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // 导入 MySQL 驱动
)

// MySQLConfig 用于配置 MySQL 数据库连接
type MySQLConfig struct {
	Username string
	Password string
	Host     string
	Port     uint16
	Database string
	Location string

	// 默认不用管这些低级配置, 如果需要设置需要走压测, 配合获取一个好的经验值

	// MaxOpenConns 控制最大连接数，避免数据库因过多连接而过载。默认 0 无限
	// 如果你知道数据库的最大连接数上限，建议设置为略低于此值，例如 70%-80%
	// 如果应用存在多个实例（例如在负载均衡下），将上限除以实例数
	// 每个实例的最大连接数 = 数据库最大连接数 / 实例数量
	// 如果设置太小：并发请求多时，连接池耗尽，导致请求排队。
	// 如果设置太大：数据库可能被连接耗尽，出现性能瓶颈。
	MaxOpenConns    int
	MaxIdleConns    int           // 控制空闲连接数，优化连接复用。默认 2 最多 2 个空闲连接
	ConnMaxIdleTime time.Duration // 控制每个空闲连接的最大生命周期
	ConnLifetime    time.Duration // 控制每个连接的最大生命周期，确保连接池中的连接健康。默认 0, 直到 MySQL 主动关闭
}

func (config *MySQLConfig) DataSourceName() string {
	// ?charset=utf8mb4：
	//	指定字符集。utf8mb4 是 MySQL 中一种支持更广泛字符集（包括表情符号）的字符集。它比 utf8 更加完整。
	// &parseTime=true:
	//	告诉 Go 驱动程序将数据库中的时间类型（如 DATETIME、TIMESTAMP）转换为 Go 中的 time.Time 类型。
	//	默认情况下，Go 驱动程序可能将这些类型解析为字符串，而设置 parseTime=true 会使它们正确地解析为 Go 的时间类型。

	// &loc=Local:
	//	指定数据库连接时使用的时区。Local 表示使用本地时区。默认使用 UTC 时间, UTC 时区 对于分布式系统很重要。
	//  这个设置只表示解析为 time.Time 类型时，使用的配置。并不改变 MySQL 的 time zone 时区信息 time_zone setting。
	if config.Location == "" {
		config.Location = "UTC"
	}

	// 构建 DSN（Data Source Name）
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=%s",
		config.Username, config.Password,
		config.Host, config.Port,
		config.Database,
		config.Location,
	)
}

func DSNToMySQLCommand(dsn string) string {
	// 解析 DSN 字符串
	u, err := url.Parse(dsn)
	if err != nil {
		slog.Error("dsn to mysql command error", "dsn", dsn, "reason", err)
		return ""
	}

	// 提取用户名和密码
	user := u.User.Username()
	password, _ := u.User.Password()

	// 提取主机地址和端口
	host := u.Hostname()
	port := u.Port()

	// 提取数据库名称
	dbname := strings.TrimPrefix(u.Path, "/")

	// 提取字符集
	q := u.Query()
	charset := q.Get("charset")

	// 生成 MySQL 命令行连接字符串
	return fmt.Sprintf("mysql -u %s -p%s -h %s -P %s --default-character-set=%s %s",
		user, password, host, port, charset, dbname)
}

// MySQLHelper 是操作 MySQL 的帮助类
type MySQLHelper struct {
	DB *sql.DB
}

// NewMySQLHelper 创建一个新的 MySQLHelper 实例
func NewMySQLHelper(ctx context.Context, config MySQLConfig) (*MySQLHelper, error) {
	// 构建 DSN（Data Source Name）
	dsn := config.DataSourceName()
	if trace.EnableDebug() {
		slog.DebugContext(ctx, "dsn and mysql cmd", "mysql", dsn, "command", DSNToMySQLCommand(dsn))
	}

	// 初始化数据库连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		slog.ErrorContext(ctx, "failed to connect to MySQL", "dsn", dsn, "reason", err)
		return nil, err
	}

	// 配置连接池
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)
	db.SetConnMaxLifetime(config.ConnLifetime)

	// 测试连接
	if err := db.Ping(); err != nil {
		slog.ErrorContext(ctx, "failed to ping MySQL", "dsn", dsn, "reason", err)
		return nil, err
	}

	slog.InfoContext(ctx, "Connected to MySQL successfully!", "database", config.Database, "username", config.Username)
	return &MySQLHelper{DB: db}, nil
}
