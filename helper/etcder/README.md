# etcd

go 操作 etcd helper 帮助类, 主要围绕常用的, 服务发现, 远端配置, 分布式锁

# etcd 环境搭建

```Docker

# 使用指定的 compose 文件在后台启动服务, 没有容器自动创建
docker compose -p etcder -f docker-compose.yml up -d

# etcd 健康状况查看
docker exec etcd1 etcdctl --endpoints=http://etcd1:2379 endpoint health
docker exec etcd1 etcdctl --endpoints=http://etcd1:2379,etcd2:2379,etcd3:2379 member list

# 暂停相关容器
docker compose -f docker-compose.yml pause

# 关闭并移除容器、网络等资源
docker compose -f docker-compose.yml down
```
