// Package structs.linkedq provides a thread-safe generic linked queue implementation.
package structs

import (
	"sync"
	"sync/atomic"
)

// LinkedQueue 是线程安全的泛型链表队列
type LinkedQueue[T any] struct {
	sync.Mutex
	head   *linkedQueueNode[T]
	tail   *linkedQueueNode[T]
	length atomic.Int32 // 单纯统计监控维度; 大多数场景不会也不应该超过 21 亿个元素
}

type linkedQueueNode[T any] struct {
	next  *linkedQueueNode[T]
	value T
}

// NewLinkedQueue 创建一个新的泛型队列
func NewLinkedQueue[T any]() *LinkedQueue[T] { return &LinkedQueue[T]{} }

// Push 将一个元素加入队列尾部
//
// 字符画主视图：
// 空队列 Push(1):
//
//	head: nil, tail: nil
//	Push(1) → [1] ← head/tail
//
// 非空队列 Push(2):
//
//	head: [1] → nil, tail: [1]
//	Push(2) → [1] → [2] ← tail
//	          ↑
//	        head
//
// 继续 Push(3):
//
//	head: [1] → [2] → nil, tail: [2]
//	Push(3) → [1] → [2] → [3] ← tail
//	          ↑
//	        head
func (q *LinkedQueue[T]) Push(value T) {
	q.Lock()
	node := &linkedQueueNode[T]{value: value}
	if q.tail == nil {
		// 空队列，head 和 tail 都指向新节点
		q.head = node
	} else {
		q.tail.next = node
	}
	q.tail = node
	q.Unlock()

	q.length.Add(1)
}

// Pop 从队列头部取出一个元素
//
// 字符画主视图：
// 队列状态 [1] → [2] → [3] → nil
//
//	  ↑              ↑
//	head            tail
//
// Pop() 操作：
//  1. 取出 head.value = 1
//  2. head 移动到下一个节点
//  3. 结果: [2] → [3] → nil
//     ↑        ↑
//     head      tail
//
// 最后一个元素 Pop():
// 队列: [3] → nil
//
//	    ↑
//	head/tail
//
// Pop() 后: nil (空队列)
//
//	head/tail 都指向 nil
func (q *LinkedQueue[T]) Pop() (value T, ok bool) {
	q.Lock()
	if q.head == nil {
		q.Unlock()
		return
	}

	value, ok = q.head.value, true
	q.head = q.head.next
	if q.head == nil {
		// 队列已空，tail 也要清空
		q.tail = nil
	}
	q.Unlock()

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
