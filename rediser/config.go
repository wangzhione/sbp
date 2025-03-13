package rediser

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
)

func NewCommand(ctx context.Context, command string) (options *redis.Options, err error) {
	return
}

// ParseRedisCLI 解析 redis-cli 命令并返回 redis.Options
func ParseRedisCLI(command string) (*redis.Options, error) {
	// 分割命令行参数
	args := strings.Fields(command)

	// 检查是否以 "redis-cli" 开头
	if len(args) == 0 || !strings.EqualFold(args[0], "redis-cli") {
		return nil, errors.New("not redis-cli head")
	}

	// 定义默认值
	options := &redis.Options{
		Addr:     "localhost:6379", // 默认地址
		Password: "",               // 默认无密码
		DB:       0,                // 默认 DB 0
	}

	// 使用 flag 库解析命令行参数
	lag := flag.NewFlagSet("redis-cli", flag.ContinueOnError)
	host := lag.String("h", "localhost", "Redis host")
	port := lag.String("p", "6379", "Redis port")
	password := lag.String("a", "", "Redis passwd")

	// 默认数据库编号就是 0, 不应该设置
	db := lag.Int("n", 0, "Redis database num")

	// 解析参数（跳过 "redis-cli"）
	err := lag.Parse(args[1:])
	if err != nil {
		return nil, err
	}

	// 设置解析后的值
	options.Addr = fmt.Sprintf("%s:%s", *host, *port)
	options.Password = *password
	options.DB = *db

	return options, nil
}
