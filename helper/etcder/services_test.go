package etcdhelper

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func Test_ServiceRegistry(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// ✅ 连接到 etcd 集群
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("failed to connect to etcd: %v", err)
	}
	defer cli.Close()

	// ✅ 创建注册对象
	key := "/services/test-service/instance1"
	value := "http://127.0.0.1:8080"
	reg := NewServiceRegistry(ctx, cli, key, value)

	// ✅ 注册并自动续约（推荐 TTL 10s）
	if err := reg.Register(10); err != nil {
		log.Fatalf("register failed: %v", err)
	}

	// ✅ 启动监听（只演示一次）
	reg.WatchServices("/services/test-service/", func(ctx context.Context, isDelete bool, key, value string) {
		if isDelete {
			fmt.Printf("🔴 服务下线: %s → %s\n", key, value)
		} else {
			fmt.Printf("🟢 服务变更: %s → %s\n", key, value)
		}
	})

	// ✅ 等待 2 秒，执行修改
	time.Sleep(2 * time.Second)
	_, err = cli.Put(ctx, key, "http://127.0.0.1:9090")
	if err != nil {
		log.Fatalf("put failed: %v", err)
	}

	// ✅ 再等 2 秒，执行删除
	time.Sleep(2 * time.Second)
	_, err = cli.Delete(ctx, key)
	if err != nil {
		log.Fatalf("delete failed: %v", err)
	}

	// ✅ 测试运行 10 秒
	time.Sleep(10 * time.Second)

	// ✅ 主动下线
	reg.Stop()
}
