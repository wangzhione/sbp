package rediser

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"strings"

	"github.com/redis/go-redis/v9"
)

// NewClusterClient redis cluster 集群模式
//
//	redis.NewClusterClient(&redis.ClusterOptions{
//	    Addrs: []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"},
//	})
func NewClusterClient(ctx context.Context, options *redis.ClusterOptions) (r *Client, err error) {
	cluster := redis.NewClusterClient(options)

	err = cluster.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
		return shard.Ping(ctx).Err()
	})
	if err != nil {
		slog.ErrorContext(ctx, "rdb.ForEachShard error", "error", err, "Addr", options.Addrs)
		return nil, err
	}

	r = &Client{cluster}
	return
}

func NewRedis(ctx context.Context, options *redis.Options) (r *Client, err error) {
	rdb := redis.NewClient(options)

	// 测试连接
	result, err := rdb.Ping(ctx).Result()
	if err != nil {
		slog.ErrorContext(ctx, "rdb.Ping(ctx).Result() error", "error", err, "Addr", options.Addr)
		return
	}
	slog.InfoContext(ctx, "Redis Success "+options.Addr, "result", result)

	r = &Client{rdb}
	return
}

// NewDefaultRedis 构建默认的 redis client
func NewDefaultRedis(ctx context.Context, command string) (r *Client, err error) {
	options, err := ParseRedisCommand(command)
	if err != nil {
		slog.ErrorContext(ctx, "ParseRedisCommand is error", "error", err, "command", command)
		return
	}

	return NewRedis(ctx, options)
}

// ParseRedisCommand 解析 redis-cli 命令并返回 redis.Options
func ParseRedisCommand(command string) (*redis.Options, error) {
	// 分割命令行参数
	args := strings.Fields(command)

	// 检查是否以 "redis-cli" 开头
	if len(args) == 0 || !strings.EqualFold(args[0], "redis-cli") {
		return nil, errors.New("not redis-cli head")
	}

	// 使用 flag 库解析命令行参数
	lag := flag.NewFlagSet("redis-cli", flag.ContinueOnError)
	host := lag.String("h", "localhost", "Redis host")
	port := lag.String("p", "6379", "Redis port")
	password := lag.String("a", "", "Redis passwd")
	username := lag.String("u", "", "Redis Username")
	db := lag.Int("n", 0, "Redis database num") // 默认数据库编号就是 0, 不应该设置

	// 解析参数（跳过 "redis-cli"）
	err := lag.Parse(args[1:])
	if err != nil {
		return nil, err
	}

	// 设置解析后的值
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", *host, *port), // 默认地址
		Password: *password,                          // 默认无密码
		DB:       *db,                                // 默认 DB 0
		Username: *username,                          // 默认 "" 一般不配置
	}

	return options, nil
}
