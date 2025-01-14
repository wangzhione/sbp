package linkedq

import "sync"

type LinkedQueue struct {
	head *linkedQueueNode
	tail *linkedQueueNode
	mu   sync.Mutex
}

type linkedQueueNode struct {
	next  *linkedQueueNode
	value uint64
}

func New() *LinkedQueue {
	node := new(linkedQueueNode)
	return &LinkedQueue{head: node, tail: node}
}

func (q *LinkedQueue) Enqueue(value uint64) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.tail.next = &linkedQueueNode{value: value}
	q.tail = q.tail.next
}

func (q *LinkedQueue) Dequeue() (value uint64, ok bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.head.next != nil {
		value, ok = q.head.next.value, true
		q.head = q.head.next
	}

	return
}

func (q *LinkedQueue) Empty() bool {
	return q.head == q.tail
}
