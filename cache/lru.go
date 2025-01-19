package cache

import (
	"container/list"
	"sync"
)

type pair[K comparable, V any] struct {
	key   K
	value V
}

// TemplateLRU go 1.18 对 LRU cache 包装
type TemplateLRU[K comparable, V any] struct {
	M sync.Mutex // 如果需要并发安全, 所有操作可以通过 M 加锁

	cap  int                 // capacity 容量 <= 0 标识不限制
	list *list.List          // entry container list
	data map[K]*list.Element // 数据池子
}

// NewTemplateLRU new cache LRU
func NewTemplateLRU[K comparable, V any](cap int) *TemplateLRU[K, V] {
	return &TemplateLRU[K, V]{
		cap:  cap,
		list: list.New(),
		data: make(map[K]*list.Element),
	}
}

// Get 获取 value, ok is true 标识存在
func (c *TemplateLRU[K, V]) Get(key K) (value V, ok bool) {
	element, ok := c.data[key]
	if ok {
		value = element.Value.(*pair[K, V]).value
		c.list.MoveToFront(element)
	}
	return
}

// Delete 通过 key 删除操作
func (c *TemplateLRU[K, V]) Delete(key K) {
	element, ok := c.data[key]
	if ok {
		delete(c.data, key)
		c.list.Remove(element)
	}
}

// Put 添加数据
func (c *TemplateLRU[K, V]) Put(key K, value V) {
	element, ok := c.data[key]
	if ok {
		element.Value.(*pair[K, V]).value = value
		c.list.MoveToFront(element)
		return
	}

	// 容量不足, 尝试清理数据
	if c.cap > 0 && c.cap <= c.list.Len() {
		delete(c.data, c.list.Remove(c.list.Back()).(*pair[K, V]).key)
	}

	// 开始填充
	c.data[key] = c.list.PushFront(&pair[K, V]{key: key, value: value})
}
