// Package etcder provides etcd-based config watcher utilities
package etcder

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/wangzhione/sbp/helper/safego"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// ConfigWatcher 监听单个配置 key（如 /configs/app.json）并自动热更新
type ConfigWatcher struct {
	ctx context.Context
	cli *clientv3.Client

	// Close cw.Close() 停止监听, 理论上不用执行, 因为往往 config watcher 应该伴随服务 or 程序长期持有
	Close context.CancelFunc

	key      string
	onUpdate func(ctx context.Context, data []byte) // 更新回调（可选）
}

// NewConfigWatcher 创建配置监听器
// key：例如 /configs/app.json
// onUpdate：
/*
    var currentConfig atomic.Pointer[AppConfig]

	func applyConfig(ctx context.Context, data []byte) {
		if data == nil {
			slog.WarnContext(ctx, "配置被删除，保持旧配置")
			return
		}

		var cfg *AppConfig
		if err := json.Unmarshal(data, &cfg); err != nil {
			slog.ErrorContext(ctx, "配置解析失败", slog.Any("error", err), slog.String("raw", string(data)))
			return
		}

		currentConfig.Store(&cfg)
		slog.InfoContext(ctx, "✅ 配置已更新", slog.Any("config", cfg))
	}
*/
func NewConfigWatcher(ctx context.Context, cli *clientv3.Client, key string, onUpdate func(context.Context, []byte)) (*ConfigWatcher, error) {
	if onUpdate == nil {
		return nil, fmt.Errorf("onUpdate must not be nil")
	}
	ctx, cancel := context.WithCancel(ctx)

	cw := &ConfigWatcher{
		ctx:      ctx,
		Close:    cancel,
		cli:      cli,
		key:      key,
		onUpdate: onUpdate,
	}

	// 初次拉取（带短超时），并得到当前 revision
	err := cw.initial()
	if err != nil {
		cancel()
		return nil, err
	}

	// 启动监听
	safego.Go(ctx, func() { cw.watch() })

	return cw, nil
}

// initial 加载初始配置值
func (cw *ConfigWatcher) initial() error {
	resp, err := cw.cli.Get(cw.ctx, cw.key)
	if err != nil {
		slog.ErrorContext(cw.ctx, "failed to load initial config", slog.Any("error", err), slog.String("key", cw.key))
		return err
	}
	if len(resp.Kvs) == 0 {
		slog.WarnContext(cw.ctx, "config key not found", slog.String("key", cw.key))
		return nil
	}

	safego.So(cw.ctx, func() { cw.onUpdate(cw.ctx, resp.Kvs[0].Value) })

	slog.InfoContext(cw.ctx, "initial config loaded", slog.String("key", cw.key))
	return nil
}

// watch 持续监听配置变化（支持 ctx.Done()，处理 watchChan 关闭）
func (cw *ConfigWatcher) watch() {
	slog.InfoContext(cw.ctx, "watching config key", slog.String("key", cw.key))

	watchChan := cw.cli.Watch(cw.ctx, cw.key)

	for {
		select {
		case <-cw.ctx.Done():
			slog.WarnContext(cw.ctx, "config watcher stopped by context", slog.String("key", cw.key), slog.Any("err", cw.ctx.Err()))
			return

		case resp, ok := <-watchChan:
			if !ok {
				slog.WarnContext(cw.ctx, "config watch channel closed", slog.String("key", cw.key))
				return
			}

			if err := resp.Err(); err != nil {
				slog.ErrorContext(cw.ctx, "watch error", slog.Any("error", err))
				continue
			}

			for _, event := range resp.Events {
				switch event.Type {
				case clientv3.EventTypePut:
					slog.InfoContext(cw.ctx, "config updated",
						slog.String("key", cw.key),
						slog.String("value", string(event.Kv.Value)))
					safego.So(cw.ctx, func() { cw.onUpdate(cw.ctx, event.Kv.Value) })

				case clientv3.EventTypeDelete:
					slog.ErrorContext(cw.ctx, "config deleted panic error", slog.String("key", cw.key))
					safego.So(cw.ctx, func() { cw.onUpdate(cw.ctx, nil) })
				}
			}
		}
	}
}

// Get 主动从 etcd 获取最新配置内容（同步拉取一次）; 更为直接可以使用 client.go::Get 方法
// 返回值 data 为 nil 表示 key 不存在
func (cw *ConfigWatcher) Get() (data []byte, err error) {
	resp, err := cw.cli.Get(cw.ctx, cw.key)
	if err != nil {
		slog.ErrorContext(cw.ctx, "failed to get config", slog.Any("error", err), slog.String("key", cw.key))
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		slog.WarnContext(cw.ctx, "config key not found", slog.String("key", cw.key))
		return nil, nil
	}
	return resp.Kvs[0].Value, nil
}
