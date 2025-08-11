package etcder

import (
	"context"
	"errors"
	"log/slog"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

/*
	// 2️⃣ 创建分布式锁
	lock, err := etcder.NewETCDLock(ctx, ETCDClient, "/locks/demo-task", 0)
	if err != nil {
		return
	}
	defer lock.Close() // 程序结束前释放 session

	// 3️⃣ 尝试在 3 秒内获取锁
	if err := lock.TimeoutLock(ctx, 3*time.Second); err != nil {
		return
	}
	defer lock.Unlock(ctx)

	... 执行需要锁保护的业务操作 ...
*/

type ETCDLock struct {
	client  *clientv3.Client
	session *concurrency.Session
	mutex   *concurrency.Mutex
	key     string
}

// NewETCDLock 创建一个 etcd 分布式锁
// If TTL is <= 0, the default 60 seconds TTL will be used. 不知道怎么设置时候, 默认填入 0
// 租约时间（TTL） → 由 Session 控制，影响锁持有多久后会被 etcd 自动释放（如果不 KeepAlive）
func NewETCDLock(ctx context.Context, cli *clientv3.Client, key string, secondTTL int) (lock *ETCDLock, err error) {
	session, err := concurrency.NewSession(cli, concurrency.WithTTL(secondTTL))
	if err != nil {
		slog.ErrorContext(ctx, "failed to create etcd session",
			slog.Any("error", err),
			slog.String("key", key),
			slog.Int("secondTTL", secondTTL),
		)
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

// Close 释放 session（锁和租约自动释放）
func (l *ETCDLock) Close() (err error) {
	err = l.session.Close()
	if err != nil {
		slog.ErrorContext(context.Background(), "failed to close etcd session",
			slog.Any("error", err),
			slog.String("key", l.key),
		)
		return
	}

	slog.InfoContext(context.Background(), "etcd session closed", slog.String("key", l.key))
	return
}

// TimeoutLock 尝试在 timeout 内获取锁（阻塞）
func (l *ETCDLock) TimeoutLock(ctx context.Context, timeout time.Duration) (err error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err = l.mutex.Lock(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.WarnContext(ctx, "lock attempt timed out", slog.String("key", l.key))
		} else {
			slog.ErrorContext(ctx, "failed to lock key",
				slog.String("key", l.key),
				slog.Any("error", err),
			)
		}
		return
	}

	slog.InfoContext(ctx, "lock acquired successfully", slog.String("key", l.key))
	return
}

// Unlock 释放锁
func (l *ETCDLock) Unlock(ctx context.Context) (err error) {
	err = l.mutex.Unlock(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to unlock key",
			slog.String("key", l.key),
			slog.Any("error", err),
		)
		return
	}

	slog.InfoContext(ctx, "lock released successfully", slog.String("key", l.key))
	return
}
