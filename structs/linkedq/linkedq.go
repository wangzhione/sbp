package linkedq

import (
	"sync"
	"sync/atomic"
)

// LinkedQueue 是线程安全的泛型链表队列
type LinkedQueue[T any] struct {
	sync.Mutex
	head   *linkedQueueNode[T]
	tail   *linkedQueueNode[T]
	length atomic.Int32
}

type linkedQueueNode[T any] struct {
	next  *linkedQueueNode[T]
	value T
}

// New 创建一个新的泛型队列
func New[T any]() *LinkedQueue[T] {
	node := new(linkedQueueNode[T]) // 哨兵节点（dummy node）
	return &LinkedQueue[T]{head: node, tail: node}
}

// Enqueue 将一个元素加入队列尾部
func (q *LinkedQueue[T]) Enqueue(value T) {
	q.Lock()
	q.tail.next = &linkedQueueNode[T]{value: value}
	q.tail = q.tail.next
	q.Unlock()

	// 业务在 Enqueue 后立即读取 Len()，是可能会提前看到不一致的数量
	// 依赖业务最终一致性, 属于 一致性 vs 性能权衡
	q.length.Add(1)
}

// Dequeue 从队列头部取出一个元素
func (q *LinkedQueue[T]) Dequeue() (value T, ok bool) {
	if q.length.Load() == 0 {
		return
	}

	q.Lock()
	if q.head.next != nil {
		value, ok = q.head.next.value, true
		q.head = q.head.next

		// ⚠️ 如果弹出的是最后一个节点，重置 tail 指回 dummy 节点
		if q.head.next == nil {
			q.tail = q.head
		}

		q.Unlock()

		q.length.Add(-1)
		return
	}
	q.Unlock()
	return
}

func (q *LinkedQueue[T]) Peek() (value T, ok bool) {
	if q.length.Load() == 0 {
		return
	}

	q.Lock()
	if q.head.next != nil {
		value, ok = q.head.next.value, true
	}
	q.Unlock()
	return
}

// Empty 判断队列是否为空 q.Len() == 0

// Len 获取 queue 中待消费数量, 多数用于线上监控, 业务运行状态
func (q *LinkedQueue[T]) Len() int32 {
	return q.length.Load()
}
