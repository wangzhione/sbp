package mq

import (
	"sync/atomic"
	"unsafe"
)

type MQueue struct {
	head unsafe.Pointer // *msqv1node
	tail unsafe.Pointer // *msqv1node
}

type msnode struct {
	next  unsafe.Pointer // *msqv1node
	value uint64
}

func New() *MQueue {
	node := unsafe.Pointer(new(msnode))
	return &MQueue{head: node, tail: node}
}

func (q *MQueue) Enqueue(value uint64) {
	node := unsafe.Pointer(&msnode{value: value})
	for {
		tail := atomic.LoadPointer(&q.tail)
		tailnode := (*msnode)(tail)
		next := atomic.LoadPointer(&tailnode.next)
		if tail == atomic.LoadPointer(&q.tail) {
			if next == nil {
				// tail.next is empty, insert new node
				if atomic.CompareAndSwapPointer(&tailnode.next, next, node) {
					atomic.CompareAndSwapPointer(&q.tail, tail, node)
					break
				}
			} else {
				atomic.CompareAndSwapPointer(&q.tail, tail, next)
			}
		}
	}
}

func (q *MQueue) Dequeue() (uint64, bool) {
	for {
		head := atomic.LoadPointer(&q.head)
		tail := atomic.LoadPointer(&q.tail)
		headnode := (*msnode)(head)
		next := atomic.LoadPointer(&headnode.next)
		if head == atomic.LoadPointer(&q.head) {
			if head == tail {
				if next == nil {
					return 0, false
				}
				atomic.CompareAndSwapPointer(&q.tail, tail, next)
			} else {
				value := ((*msnode)(next)).value
				if atomic.CompareAndSwapPointer(&q.head, head, next) {
					return value, true
				}
			}
		}
	}
}
