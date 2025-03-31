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
func New[T any]() *LinkedQueue[T] { return &LinkedQueue[T]{} }

// Push 将一个元素加入队列尾部
func (q *LinkedQueue[T]) Push(value T) {
	q.Lock()
	defer q.Unlock() // 很土, 但好用能用

	node := &linkedQueueNode[T]{value: value}
	if q.tail == nil {
		// 空队列，head 和 tail 都指向新节点
		q.head = node
	} else {
		q.tail.next = node
	}
	q.tail = node

	q.length.Add(1)
}

// Pop 从队列头部取出一个元素
func (q *LinkedQueue[T]) Pop() (value T, ok bool) {
	q.Lock()
	defer q.Unlock()

	if q.head == nil {
		return
	}

	value, ok = q.head.value, true
	q.head = q.head.next
	if q.head == nil {
		// 队列已空，tail 也要清空
		q.tail = nil
	}

	q.length.Add(-1)
	return
}

func (q *LinkedQueue[T]) Peek() (value T, ok bool) {
	q.Lock()
	defer q.Unlock()

	if q.head != nil {
		value, ok = q.head.value, true
	}
	return
}

// Empty 判断队列是否为空 q.Len() == 0

// Len 获取 queue 中待消费数量, 多数用于线上监控, 业务运行状态
func (q *LinkedQueue[T]) Len() int32 {
	return q.length.Load()
}
