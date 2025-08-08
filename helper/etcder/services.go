package etcder

import (
	"context"
	"log/slog"
	"time"

	"github.com/wangzhione/sbp/helper/safego"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// ServiceRegistry 封装服务注册与监听逻辑（基于 etcd 实现）
type ServiceRegistry struct {
	ctx        context.Context
	cancelFunc context.CancelFunc

	client  *clientv3.Client
	leaseID clientv3.LeaseID

	key   string
	value string
}

// NewServiceRegistry 创建一个服务注册对象（推荐 client 复用）
func NewServiceRegistry(ctx context.Context, client *clientv3.Client, key, value string) *ServiceRegistry {
	ctx, cancel := context.WithCancel(ctx)
	return &ServiceRegistry{
		client:     client,
		ctx:        ctx,
		cancelFunc: cancel,
		key:        key,
		value:      value,
	}
}

// Register 注册服务并自动续约
// secondTTL: 服务注册的租约时间（单位：秒）普通 Web 服务推荐 10s; 后台服务推荐 30s; 高敏感服务推荐 3-5s;
func (s *ServiceRegistry) Register(secondTTL int64) error {
	leaseResp, err := s.client.Grant(s.ctx, secondTTL)
	if err != nil {
		slog.ErrorContext(s.ctx, "failed to create lease", slog.Any("error", err), slog.Int64("secondTTL", secondTTL))
		return err
	}

	_, err = s.client.Put(s.ctx, s.key, s.value, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		slog.ErrorContext(s.ctx, "failed to register service", slog.Any("error", err),
			slog.String("key", s.key), slog.String("value", s.value),
			slog.Int64("secondTTL", secondTTL), slog.Int64("leaseID", int64(leaseResp.ID)))
		return err
	}

	s.leaseID = leaseResp.ID
	slog.InfoContext(s.ctx, "service registered", slog.String("key", s.key), slog.String("value", s.value),
		slog.Int64("leaseID", int64(leaseResp.ID)))

	safego.Go(s.ctx, func(context.Context) { s.keepAlive() })
	return nil
}

func (s *ServiceRegistry) keepAlive() {
	const maxRetries = 3

	retryCount := 0
	for {
		if s.ctx.Err() != nil {
			slog.InfoContext(s.ctx, "keepalive stopped (context canceled)", slog.Int64("leaseID", int64(s.leaseID)))
			return
		}

		ch, err := s.client.KeepAlive(s.ctx, s.leaseID)
		if err != nil {
			retryCount++
			slog.ErrorContext(s.ctx, "failed to start keepalive", slog.Any("error", err), slog.Int("retryCount", retryCount))
			if retryCount >= maxRetries {
				slog.ErrorContext(s.ctx, "keepalive retry limit reached, exiting")
				return
			}
			select {
			case <-s.ctx.Done():
				return
			case <-time.After(time.Second):
			}
			continue
		}

		// 启动成功，重置计数
		retryCount = 0

		for {
			select {
			case <-s.ctx.Done():
				return
			case ka, ok := <-ch:
				if !ok {
					retryCount++
					slog.WarnContext(s.ctx, "keepalive channel closed", slog.Int("retryCount", retryCount))
					if retryCount >= maxRetries {
						slog.ErrorContext(s.ctx, "keepalive stream retry limit reached, exiting")
						return
					}
					time.Sleep(time.Second)
					break // 跳出内层，重连
				}
				slog.DebugContext(s.ctx, "lease keepalive", slog.Int64("secondTTL", int64(ka.TTL)))
			}
		}
	}
}

// WatchServices 异步监听服务变化（自动恢复、支持 ctx.Done()）
func (s *ServiceRegistry) WatchServices(prefix string, onChange func(ctx context.Context, isDelete bool, key, value string)) {
	safego.Go(s.ctx, func(ctx context.Context) {
		slog.InfoContext(ctx, "starting watch", slog.String("prefix", prefix))

		watchChan := s.client.Watch(ctx, prefix, clientv3.WithPrefix())

		for {
			select {
			case <-ctx.Done():
				slog.WarnContext(ctx, "watch cancelled by context", slog.String("prefix", prefix), slog.Any("err", ctx.Err()))
				return

			case resp, ok := <-watchChan:
				if !ok {
					slog.WarnContext(ctx, "watch channel closed", slog.String("prefix", prefix))
					return
				}

				if err := resp.Err(); err != nil {
					slog.ErrorContext(ctx, "watch error", slog.Any("error", err), slog.String("prefix", prefix))
					continue // 自动恢复
				}

				for _, event := range resp.Events {
					key, val := string(event.Kv.Key), string(event.Kv.Value)
					switch event.Type {
					case clientv3.EventTypePut:
						slog.InfoContext(ctx, "event added or updated", slog.String("key", key), slog.String("value", val))
						onChange(ctx, false, key, val)

					case clientv3.EventTypeDelete:
						slog.InfoContext(ctx, "event deleted", slog.String("key", key), slog.String("value", val))
						onChange(ctx, true, key, val)
					}
				}
			}
		}
	})
}

// Stop 停止服务注册、撤销租约
func (s *ServiceRegistry) Stop() {
	s.cancelFunc()
	if s.leaseID != 0 {
		_ = Revoke(s.ctx, s.client, s.leaseID)
	}
	slog.InfoContext(s.ctx, "service stopped", slog.String("key", s.key), slog.String("value", s.value), slog.Int64("leaseID", int64(s.leaseID)))
}
