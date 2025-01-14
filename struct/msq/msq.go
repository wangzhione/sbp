package msq

import (
	"sync/atomic"
	"unsafe"
)

type MSQueue struct {
	head unsafe.Pointer // *msqv1node
	tail unsafe.Pointer // *msqv1node
}

type msqnode struct {
	next  unsafe.Pointer // *msqv1node
	value uint64
}

func New() *MSQueue {
	node := unsafe.Pointer(new(msqnode))
	return &MSQueue{head: node, tail: node}
}

func (q *MSQueue) Enqueue(value uint64) {
	node := unsafe.Pointer(&msqnode{value: value})
	for {
		tail := atomic.LoadPointer(&q.tail)
		tailnode := (*msqnode)(tail)
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

func (q *MSQueue) Dequeue() (uint64, bool) {
	for {
		head := atomic.LoadPointer(&q.head)
		tail := atomic.LoadPointer(&q.tail)
		headnode := (*msqnode)(head)
		next := atomic.LoadPointer(&headnode.next)
		if head == atomic.LoadPointer(&q.head) {
			if head == tail {
				if next == nil {
					return 0, false
				}
				atomic.CompareAndSwapPointer(&q.tail, tail, next)
			} else {
				value := ((*msqnode)(next)).value
				if atomic.CompareAndSwapPointer(&q.head, head, next) {
					return value, true
				}
			}
		}
	}
}
