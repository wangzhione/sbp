package mysql

import (
	"fmt"
	"sbp/util/cast"
	"strings"
)

// MySQLConfig 用于配置 MySQL 数据库连接
type MySQLConfig struct {
	Username string
	Password string
	Host     string
	Port     uint16
	Database string
	Location string

	// 默认不用管这些低级配置, 如果需要设置要么性能出现问题, 还需要压测配合, 或者 DB 方需要我们配合, 配合获取一个好的经验值

	// MaxOpenConns 控制最大连接数，避免数据库因过多连接而过载。默认 0 无限
	// 如果你知道数据库的最大连接数上限，建议设置为略低于此值，例如 70%-80%
	// 如果应用存在多个实例（例如在负载均衡下），将上限除以实例数
	// 每个实例的最大连接数 = 数据库最大连接数 / 实例数量
	// 如果设置太小：并发请求多时，连接池耗尽，导致请求排队。
	// 如果设置太大：数据库可能被连接耗尽，出现性能瓶颈。
	MaxOpenConns *int
	MaxIdleConns *int // 控制空闲连接数，优化连接复用。默认 2 最多 2 个空闲连接

	// 对于连接相关最大生命周期, 默认是 0 表示永久, 随自行 close 或 mysql 主动关闭
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

func (config *MySQLConfig) Command() string {
	// 生成 MySQL 命令行连接字符串
	// mysql -u <用户名> -p<密码> -h <主机名> -P <端口号> --default-character-set=utf8mb4 <数据库名>
	return fmt.Sprintf("mysql -u %s -p%s -h %s -P %d %s --default-character-set=utf8mb4",
		config.Username, config.Password, config.Host, config.Port, config.Database)
}

// ParseCommand parses a MySQL command string into a MySQLConfig object
func ParseCommand(command string) (*MySQLConfig, error) {
	// Split the command into arguments
	args := strings.Fields(command)

	if len(args) < 2 || args[0] != "mysql" {
		return nil, fmt.Errorf("invalid MySQL command format")
	}

	config := &MySQLConfig{}

	// Parse each argument
	for i := 1; i < len(args); i++ {
		arg := args[i]

		if arg == "-u" {
			// Handle `-u` with a value in the next argument
			if i+1 < len(args) {
				config.Username = args[i+1]
				i++ // Skip the next argument since it's already used
			}
		} else if arg == "-p" {
			// Handle `-p` with a value in the next argument
			if i+1 < len(args) {
				config.Password = args[i+1]
				i++ // Skip the next argument since it's already used
			}
		} else if strings.HasPrefix(arg, "-p") {
			// Handle `-pPassword` format
			config.Password = strings.TrimPrefix(arg, "-p")
		} else if arg == "-h" {
			// Handle `-h` with a value in the next argument
			if i+1 < len(args) {
				config.Host = args[i+1]
				i++ // Skip the next argument since it's already used
			}
		} else if arg == "-P" {
			// Handle `-P` with a value in the next argument
			if i+1 < len(args) {
				port, err := cast.StringToIntE[uint16](args[i+1])
				if err != nil {
					return nil, fmt.Errorf("invalid port format: %s %v", args[i+1], err)
				}
				config.Port = port
				i++ // Skip the next argument since it's already used
			}
		} else if strings.HasPrefix(arg, "--default-character-set=") {
			// Handle `--default-character-set` format
			continue
		} else {
			// Assume the last argument is the database name
			config.Database = arg
		}
	}

	// Validate required fields
	if config.Username == "" || config.Password == "" || config.Host == "" || config.Port == 0 || config.Database == "" {
		return nil, fmt.Errorf("missing required fields in MySQL command")
	}

	return config, nil
}

// ConvertDSNToCommand 将 DataSourceName 转换为 mysql 命令行格式
func ConvertDSNToCommand(dsn string) (string, error) {
	// 分割 DSN 为主连接部分和查询参数部分
	parts := strings.Split(dsn, "?")
	connPart := parts[0]

	// 校验格式是否包含 "@"
	userInfoAndAddr := strings.SplitN(connPart, "@", 2)
	if len(userInfoAndAddr) < 2 {
		return "", fmt.Errorf("invalid DSN: missing '@' separator")
	}

	// 提取用户信息
	userInfo := userInfoAndAddr[0]
	protocolAndAddr := userInfoAndAddr[1]

	// 分割用户名和密码
	userAndPass := strings.SplitN(userInfo, ":", 2)
	username := userAndPass[0]
	password := ""
	if len(userAndPass) > 1 {
		password = userAndPass[1]
	}

	// 校验地址是否包含 "("
	protocolAndAddress := strings.SplitN(protocolAndAddr, "(", 2)
	if len(protocolAndAddress) < 2 {
		return "", fmt.Errorf("invalid DSN: missing '(' separator for address")
	}

	// 提取地址和数据库名
	addressAndDb := strings.SplitN(protocolAndAddress[1], ")/", 2)
	if len(addressAndDb) < 2 {
		return "", fmt.Errorf("invalid DSN: missing database name")
	}
	address := strings.TrimRight(addressAndDb[0], ")")
	dbName := addressAndDb[1]

	// 分割地址为主机名和端口
	hostAndPort := strings.SplitN(address, ":", 2)
	host := hostAndPort[0]
	port := "3306" // 默认端口
	if len(hostAndPort) > 1 {
		port = hostAndPort[1]
	}

	// 构造 mysql 命令行格式
	command := fmt.Sprintf(
		"mysql -u %s -p%s -h %s -P %s %s --default-character-set=utf8mb4",
		username, password, host, port, dbName,
	)

	return command, nil
}
