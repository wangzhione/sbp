package etcder

import (
	"context"
	"log/slog"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type ETCDLock struct {
	client  *clientv3.Client
	session *concurrency.Session
	mutex   *concurrency.Mutex
	key     string
}

// NewETCDLock 创建一个 etcd 分布式锁
func NewETCDLock(ctx context.Context, cli *clientv3.Client, key string, secondTTL int) (lock *ETCDLock, err error) {
	session, err := concurrency.NewSession(cli, concurrency.WithTTL(secondTTL))
	if err != nil {
		slog.ErrorContext(ctx, "failed to create etcd session", slog.Any("error", err), slog.String("key", key), slog.Int("secondTTL", secondTTL))
		return nil, err
	}

	mutex := concurrency.NewMutex(session, key)

	lock = &ETCDLock{
		client:  cli,
		session: session,
		mutex:   mutex,
		key:     key,
	}

	return
}

// TimeoutLock 尝试在 timeout 内获取锁（阻塞）
func (l *ETCDLock) TimeoutLock(ctx context.Context, timeout time.Duration) (err error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err = l.mutex.Lock(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to lock key", slog.String("key", l.key), slog.Any("error", err))
		return
	}

	return
}

// Unlock 释放锁
func (l *ETCDLock) Unlock(ctx context.Context) (err error) {
	err = l.mutex.Unlock(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to unlock key %s", slog.String("key", l.key), slog.Any("error", err))
	}
	return
}

// Close 释放 session（锁和租约自动释放）
func (l *ETCDLock) Close() (err error) {
	err = l.session.Close()
	if err != nil {
		slog.ErrorContext(context.Background(), "failed to close etcd session", slog.Any("error", err), slog.String("key", l.key))
	}
	return
}
