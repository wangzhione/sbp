// Package etcder provides etcd-based service discovery helpers
package etcder

import (
	"context"
	"log/slog"
	"time"

	"github.com/wangzhione/sbp/structs"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// DefaultDialTimeout is the default timeout for etcd client connections
var DefaultDialTimeout = 6 * time.Second

// New 创建一个 etcd client（推荐复用）
func New(ctx context.Context, endpoints []string, timeout ...time.Duration) (cli *clientv3.Client, err error) {
	// 如果没有传入超时时间，则使用默认值
	dialTimeout := structs.Ternary(len(timeout) == 0, DefaultDialTimeout, timeout[0])

	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to create etcd client",
			slog.Any("error", err), slog.Any("endpoints", endpoints), slog.Duration("timeout", dialTimeout))
		return
	}

	slog.InfoContext(ctx, "etcd client connected", slog.Any("endpoints", endpoints))
	return
}

// Close 关闭 etcd 客户端连接
func Close(ctx context.Context, cli *clientv3.Client) {
	endpoints := cli.Endpoints()

	// 理论上, 很难执行这个方法, 实战场景可以不主动尝试 Clone
	if err := cli.Close(); err != nil {
		slog.ErrorContext(ctx, "failed to close etcd client", slog.Any("error", err), slog.Any("endpoints", endpoints))
	} else {
		slog.InfoContext(ctx, "etcd client closed", slog.Any("endpoints", endpoints))
	}
}

// Get 获取指定 key 的值
func Get(ctx context.Context, cli *clientv3.Client, key string) (data []byte, err error) {
	resp, err := cli.Get(ctx, key)
	if err != nil {
		slog.ErrorContext(ctx, "etcd get failed", slog.Any("error", err), slog.String("key", key))
		return
	}
	if len(resp.Kvs) == 0 {
		slog.WarnContext(ctx, "key not found", slog.String("key", key))
		return nil, nil
	}
	data = resp.Kvs[0].Value
	return
}

// GetPrefix 获取指定前缀下的所有键值对
func GetPrefix(ctx context.Context, cli *clientv3.Client, prefix string) (services map[string]string, err error) {
	// 使用 WithPrefix 获取所有以 prefix 开头的键值对
	resp, err := cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		slog.ErrorContext(ctx, "etcd get failed", slog.Any("error", err), slog.String("prefix", prefix))
		return
	}

	services = make(map[string]string)
	for _, kv := range resp.Kvs {
		services[string(kv.Key)] = string(kv.Value)
	}
	return
}

// GetPrefixKeys 获取指定前缀下的所有 keys
func GetPrefixKeys(ctx context.Context, cli *clientv3.Client, prefix string) (keys []string, err error) {
	resp, err := cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		slog.ErrorContext(ctx, "etcd get with prefix failed", slog.Any("error", err), slog.String("prefix", prefix))
		return
	}

	keys = make([]string, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		keys = append(keys, string(kv.Key))
	}
	return
}

// GetPrefixValues 获取指定前缀下的所有 values
func GetPrefixValues(ctx context.Context, cli *clientv3.Client, prefix string) (values []string, err error) {
	resp, err := cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		slog.ErrorContext(ctx, "etcd get with prefix failed", slog.Any("error", err), slog.String("prefix", prefix))
		return
	}

	values = make([]string, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		values = append(values, string(kv.Value))
	}
	return
}

// WatchPrefix 监听指定前缀下键值对变化（新增/更新/删除）
// 每次变更都会调用 onChange 回调：onChange func(ctx context.Context, isDelete bool, key, value string)
func WatchPrefix(ctx context.Context, cli *clientv3.Client, prefix string, onChange func(ctx context.Context, isDelete bool, key, value string)) {
	watchChan := cli.Watch(ctx, prefix, clientv3.WithPrefix())
	slog.InfoContext(ctx, "watching key-values", slog.String("prefix", prefix))

	for resp := range watchChan {
		if err := resp.Err(); err != nil {
			slog.ErrorContext(ctx, "etcd watch error", slog.Any("error", err))
			continue
		}

		for _, ev := range resp.Events {
			key, val := string(ev.Kv.Key), string(ev.Kv.Value)
			switch ev.Type {
			case clientv3.EventTypePut: // 服务注册 or 变更(如 ip 变了)
				slog.InfoContext(ctx, "event added or updated", slog.String("key", key))
				onChange(ctx, false, key, val)
			case clientv3.EventTypeDelete: // 服务下线, 异常失去联系等
				slog.InfoContext(ctx, "event deleted", slog.String("key", key))
				onChange(ctx, true, key, val)
			}
		}
	}
}

// Grant 注册一个服务，写入指定 key 和 value，并绑定 TTL 租约
func Grant(ctx context.Context, cli *clientv3.Client, key, value string, ttlSeconds int64) (leaseID clientv3.LeaseID, err error) {
	// 创建租约
	leaseResp, err := cli.Grant(ctx, ttlSeconds)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create lease", slog.Any("error", err))
		return
	}

	leaseID = leaseResp.ID
	// 写入带租约的 key-value
	_, err = cli.Put(ctx, key, value, clientv3.WithLease(leaseID))
	if err != nil {
		slog.ErrorContext(ctx, "failed to register service", slog.Any("error", err),
			slog.String("key", key), slog.String("value", value),
			slog.Int64("leaseID", int64(leaseID)))
		return
	}

	slog.InfoContext(ctx, "grant service registered", slog.String("key", key), slog.Int64("ttl", ttlSeconds), slog.Int64("leaseID", int64(leaseID)))
	return
}

// KeepAlive 启动租约续租 goroutine，确保服务不会因超时自动下线
func KeepAlive(ctx context.Context, cli *clientv3.Client, leaseID clientv3.LeaseID) {
	ch, err := cli.KeepAlive(ctx, leaseID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to start lease keepalive", slog.Any("error", err), slog.Int64("leaseID", int64(leaseID)))
		return
	}

	slog.InfoContext(ctx, "started lease keepalive", slog.Int64("leaseID", int64(leaseID)))
	for {
		select {
		case <-ctx.Done():
			slog.WarnContext(ctx, "lease keepalive stopped by context", slog.Int64("leaseID", int64(leaseID)))
			return
		// 收到响应, 续租情况
		case resp, ok := <-ch:
			if !ok {
				slog.WarnContext(ctx, "lease keepalive channel closed", slog.Int64("leaseID", int64(leaseID)))
				return
			}
			slog.DebugContext(ctx, "lease keepalive response", slog.Int64("leaseID", int64(resp.ID)), slog.Int64("ttl", int64(resp.TTL)))
		}
	}
}

// Revoke 主动撤销租约（可用于服务下线）
func Revoke(ctx context.Context, cli *clientv3.Client, leaseID clientv3.LeaseID) error {
	resp, err := cli.Revoke(ctx, leaseID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to revoke lease",
			slog.Any("error", err),
			slog.Int64("leaseID", int64(leaseID)),
		)
		return err
	}

	slog.InfoContext(ctx, "lease revoked",
		slog.Any("header", resp.Header), // 记录集群元信息
		slog.Int64("leaseID", int64(leaseID)),
	)
	return nil
}

// Put 设置指定 key 的值
func Put(ctx context.Context, cli *clientv3.Client, key, value string) error {
	resp, err := cli.Put(ctx, key, value, clientv3.WithPrevKV()) // 加 WithPrevKV 可取旧值
	if err != nil {
		slog.ErrorContext(ctx, "etcd put failed",
			slog.Any("error", err),
			slog.String("key", key),
			slog.String("value", value),
		)
		return err
	}

	if resp.PrevKv != nil {
		slog.InfoContext(ctx, "key updated",
			slog.Any("header", resp.Header),
			slog.String("key", key),
			slog.String("old_value", string(resp.PrevKv.Value)),
			slog.String("new_value", value),
		)
	} else {
		slog.InfoContext(ctx, "key created",
			slog.Any("header", resp.Header),
			slog.String("key", key),
			slog.String("value", value),
		)
	}
	return nil
}

// Delete 删除指定 key
func Delete(ctx context.Context, cli *clientv3.Client, key string) error {
	resp, err := cli.Delete(ctx, key, clientv3.WithPrevKV())
	if err != nil {
		slog.ErrorContext(ctx, "etcd delete failed",
			slog.Any("error", err),
			slog.String("key", key),
		)
		return err
	}

	slog.InfoContext(ctx, "delete completed",
		slog.String("key", key),
		slog.Int64("deleted_count", resp.Deleted),
		slog.Any("header", resp.Header),
	)

	if len(resp.PrevKvs) > 0 {
		for _, kv := range resp.PrevKvs {
			slog.InfoContext(ctx, "key deleted",
				slog.String("key", string(kv.Key)),
				slog.String("old_value", string(kv.Value)),
			)
		}
	}
	return nil
}
