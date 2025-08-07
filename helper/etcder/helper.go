// Package etcdhelper provides etcd-based service discovery helpers
package etcdhelper

import (
	"context"
	"log/slog"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// DefaultDialTimeout is the default timeout for etcd client connections
var DefaultDialTimeout = 6 * time.Second

// NewClientV3 创建一个 etcd client（推荐复用）
func NewClientV3(ctx context.Context, endpoints []string) (cli *clientv3.Client, err error) {
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: DefaultDialTimeout,
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to create etcd client", slog.Any("error", err), slog.Any("endpoints", endpoints))
		return
	}

	slog.InfoContext(ctx, "etcd client connected", slog.Any("endpoints", endpoints))
	return
}

// CloseClient 关闭 etcd 客户端连接
func CloseClient(ctx context.Context, cli *clientv3.Client) {
	if cli != nil {
		if err := cli.Close(); err != nil {
			slog.ErrorContext(ctx, "failed to close etcd client", slog.Any("error", err))
		} else {
			slog.InfoContext(ctx, "etcd client closed")
		}
	}
}

// GetKeyValues 获取指定前缀下的所有键值对
func GetKeyValues(ctx context.Context, cli *clientv3.Client, prefix string) (services map[string]string, err error) {
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

// WatchKeyValues 监听指定前缀下键值对变化（新增/更新/删除）
// 每次变更都会调用 onChange 回调：onChange func(ctx context.Context, isDelete bool, key, value string)
func WatchKeyValues(ctx context.Context, cli *clientv3.Client, prefix string, onChange func(ctx context.Context, isDelete bool, key, value string)) {
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
				slog.InfoContext(ctx, "event added or updated", slog.String("key", key), slog.String("value", val))
				onChange(ctx, false, key, val)
			case clientv3.EventTypeDelete: // 服务下线, 异常失去联系等
				slog.InfoContext(ctx, "event deleted", slog.String("key", key), slog.String("value", val))
				onChange(ctx, true, key, val)
			default:
				slog.WarnContext(ctx, "unknown event type", slog.String("type", ev.Type.String()), slog.String("key", key), slog.String("value", val))
			}
		}
	}
}

// RegisterService 注册一个服务，写入指定 key 和 value，并绑定 TTL 租约
func RegisterService(ctx context.Context, cli *clientv3.Client, key, value string, ttlSeconds int64) (leaseID clientv3.LeaseID, err error) {
	// 创建租约
	leaseResp, err := cli.Grant(ctx, ttlSeconds)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create lease", slog.Any("error", err))
		return
	}

	// 写入带租约的 key-value
	_, err = cli.Put(ctx, key, value, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		slog.ErrorContext(ctx, "failed to register service", slog.Any("error", err), slog.String("key", key), slog.String("value", value))
		return
	}

	slog.InfoContext(ctx, "service registered", slog.String("key", key), slog.String("value", value), slog.Int64("ttl", ttlSeconds))
	leaseID = leaseResp.ID
	return
}

// KeepAliveLease 启动租约续租 goroutine，确保服务不会因超时自动下线
func KeepAliveLease(ctx context.Context, cli *clientv3.Client, leaseID clientv3.LeaseID) {
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
		case resp, ok := <-ch:
			if !ok {
				slog.WarnContext(ctx, "lease keepalive channel closed", slog.Int64("leaseID", int64(leaseID)))
				return
			}
			slog.DebugContext(ctx, "lease keepalive response", slog.Int64("leaseID", int64(resp.ID)), slog.Int64("ttl", int64(resp.TTL)))
		}
	}
}

// RevokeLease 主动撤销租约（可用于服务下线）
func RevokeLease(ctx context.Context, cli *clientv3.Client, leaseID clientv3.LeaseID) error {
	_, err := cli.Revoke(ctx, leaseID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to revoke lease", slog.Any("error", err), slog.Int64("leaseID", int64(leaseID)))
		return err
	}
	slog.InfoContext(ctx, "lease revoked", slog.Int64("leaseID", int64(leaseID)))
	return nil
}

// ExistsKey 检查指定 key 是否存在
func ExistsKey(ctx context.Context, cli *clientv3.Client, key string) (exists bool, err error) {
	resp, err := cli.Get(ctx, key)
	if err != nil {
		slog.ErrorContext(ctx, "etcd get failed", slog.Any("error", err), slog.String("key", key))
		return
	}
	exists = len(resp.Kvs) > 0
	return
}

// PutKey 设置指定 key 的值
func PutKey(ctx context.Context, cli *clientv3.Client, key, value string) (err error) {
	_, err = cli.Put(ctx, key, value)
	if err != nil {
		slog.ErrorContext(ctx, "etcd put failed", slog.Any("error", err), slog.String("key", key), slog.String("value", value))
		return
	}
	slog.InfoContext(ctx, "key set successfully", slog.String("key", key), slog.String("value", value))
	return
}

// DeleteKey 删除指定 key
func DeleteKey(ctx context.Context, cli *clientv3.Client, key string) (err error) {
	_, err = cli.Delete(ctx, key)
	if err != nil {
		slog.ErrorContext(ctx, "etcd delete failed", slog.Any("error", err), slog.String("key", key))
		return
	}
	slog.InfoContext(ctx, "key deleted successfully", slog.String("key", key))
	return
}

// GetKey 获取指定 key 的值
func GetKey(ctx context.Context, cli *clientv3.Client, key string) (value string, err error) {
	resp, err := cli.Get(ctx, key)
	if err != nil {
		slog.ErrorContext(ctx, "etcd get failed", slog.Any("error", err), slog.String("key", key))
		return
	}
	if len(resp.Kvs) == 0 {
		slog.WarnContext(ctx, "key not found", slog.String("key", key))
		return "", nil
	}
	value = string(resp.Kvs[0].Value)
	slog.InfoContext(ctx, "key retrieved successfully", slog.String("key", key), slog.String("value", value))
	return
}

// GetKeys 获取指定前缀下的所有 keys
func GetKeys(ctx context.Context, cli *clientv3.Client, prefix string) (keys []string, err error) {
	resp, err := cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		slog.ErrorContext(ctx, "etcd get with prefix failed", slog.Any("error", err), slog.String("prefix", prefix))
		return
	}

	keys = make([]string, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		keys = append(keys, string(kv.Key))
	}
	slog.InfoContext(ctx, "keys retrieved with prefix", slog.String("prefix", prefix), slog.Int("count", len(keys)))
	return
}

// GetValues 获取指定前缀下的所有 values
func GetValues(ctx context.Context, cli *clientv3.Client, prefix string) (values []string, err error) {
	resp, err := cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		slog.ErrorContext(ctx, "etcd get with prefix failed", slog.Any("error", err), slog.String("prefix", prefix))
		return
	}

	values = make([]string, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		values = append(values, string(kv.Value))
	}
	slog.InfoContext(ctx, "values retrieved with prefix", slog.String("prefix", prefix), slog.Int("count", len(values)))
	return
}
