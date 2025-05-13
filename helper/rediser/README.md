redis

详情 

@see https://github.com/redis/go-redis

@see https://redis.uptrace.dev/zh/guide/

# 帮助手册

redis 相关业务实战项目需要的, 会纪录在这里

## 1. set

```
SET mykey "value"          # 直接设置
SET mykey "value" EX 60    # 设置并 60 秒后过期
SET mykey "value" PX 500   # 500 毫秒后过期
SET mykey "value" NX       # 仅在 mykey 不存在时设置
SET mykey "value" XX       # 仅在 mykey 存在时更新
```

## redis stream 用作任务队列

@see queue.go 部分. 非常适合小公司组内, 少范围 B 端业务. 怎么方便, 怎么省钱怎么来.

对于消息队列部分, 比较复杂, 善于使用 ChatGPT 人脑存储