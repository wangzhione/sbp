// Package structs list provides an implementation of a doubly-linked list with a front
// and back. The individual nodes of the list are publicly exposed so that the
// user can have fine-grained control over the list.
package structs

// Lode 是 ListNode 简写, 表示链表节点.
type Lode[T any] struct {
	// Prev = Previous, 表示前一个节点; Next 表示后一个节点.
	Prev, Next *Lode[T]
	V          T
}

func NewLode[T any](value T) *Lode[T] { return &Lode[T]{V: value} }

// Value returns the value of the node. If the node is nil, it returns the zero value of T.
func (node *Lode[T]) Value() (value T) {
	if node != nil {
		return node.V
	}
	return
}

// List implements a doubly-linked list.
type List[T any] struct {
	// Head 表示头节点; Tail 表示尾节点
	Head, Tail *Lode[T]
}

// List[T]{} or &List[T]{} 都可以, 推荐前者 struct 方式声明

func NewList[T any](values ...T) *List[T] {
	if len(values) == 0 {
		return &List[T]{}
	}

	// init head and tail
	list := &List[T]{Head: &Lode[T]{V: values[0]}}
	list.Tail = list.Head

	for i := 1; i < len(values); i++ {
		list.Tail.Next = &Lode[T]{Prev: list.Tail, V: values[i]}
		list.Tail = list.Tail.Next
	}

	return list
}

// PushBackPartial adds the node 'node' to the back of the list.
func (list *List[T]) PushBackPartial(node *Lode[T]) {
	if list.Tail == nil {
		list.Head = node
		list.Tail = node
		return
	}

	node.Prev = list.Tail
	list.Tail.Next = node
	list.Tail = node
}

// PushBackNode adds the node 'node' to the back of the list.
func (list *List[T]) PushBackNode(node *Lode[T]) {
	if node == nil {
		return
	}

	node.Next = nil
	list.PushBackPartial(node)
}

// PushBack adds 'value' to the end of the list.
func (list *List[T]) PushBack(value T) {
	list.PushBackPartial(NewLode(value))
}

// PushFrontPartial adds the node 'node' to the front of the list.
func (list *List[T]) PushFrontPartial(node *Lode[T]) {
	if list.Head == nil {
		list.Head = node
		list.Tail = node
		return
	}

	node.Next = list.Head
	list.Head.Prev = node
	list.Head = node
}

// PushFrontNode adds the node 'node' to the front of the list.
func (list *List[T]) PushFrontNode(node *Lode[T]) {
	if node == nil {
		return
	}

	node.Prev = nil
	list.PushFrontPartial(node)
}

// PushFront adds 'value' to the beginning of the list.
func (list *List[T]) PushFront(value T) {
	list.PushFrontPartial(NewLode(value))
}

// Detach 从链表中摘除节点 'node'。摘除是个重操作, 所以相对严格一点
func (list *List[T]) Detach(node *Lode[T]) {
	if node == nil || list.Head == nil {
		return
	}

	if node.Next != nil {
		node.Next.Prev = node.Prev
	} else {
		list.Tail = node.Prev
	}

	if node.Prev != nil {
		node.Prev.Next = node.Next
	} else {
		list.Head = node.Next
	}

	node.Next, node.Prev = nil, nil
}

// InsertAfter 在 node 之后插入 next 并返回它。
// next 不应该已在另一个列表中（否则可能破坏另一个列表的结构）
func (list *List[T]) InsertAfter(node *Lode[T], next *Lode[T]) {
	if node == nil || list.Head == nil || node == next {
		return
	}

	// 将 next 插入到 node 后面
	next.Prev = node
	next.Next = node.Next

	if node.Next != nil {
		node.Next.Prev = next
	} else {
		// node 是尾节点，更新列表尾指针
		list.Tail = next
	}
	node.Next = next
}

// InsertBefore 在 node 之前插入 prev 并返回它。
// 注意：prev 不应该已在另一个列表中（否则可能破坏另一个列表的结构）
func (list *List[T]) InsertBefore(node *Lode[T], prev *Lode[T]) {
	if node == nil || list.Head == nil || node == prev {
		return
	}

	// 将 prev 插入到 node 前面
	prev.Next = node
	prev.Prev = node.Prev

	if node.Prev != nil {
		node.Prev.Next = prev
	} else {
		// node 是头节点，更新列表头指针
		list.Head = prev
	}
	node.Prev = prev
}

// Len returns the number of elements (O(n)).
func (list *List[T]) Len() (count int) {
	if list == nil {
		return
	}
	for n := list.Head; n != nil; n = n.Next {
		count++
	}
	return
}

// IsEmpty reports whether the list is empty.
func (list *List[T]) IsEmpty() bool { return list == nil || list.Head == nil }

// PopFront removes and returns the front node.
func (list *List[T]) PopFront() *Lode[T] {
	node := list.Head
	list.Detach(node)
	return node
}

// PopBack removes and returns the back node.
func (list *List[T]) PopBack() *Lode[T] {
	node := list.Tail
	list.Detach(node)
	return node
}

// MoveToFront moves 'node' to the front.
func (list *List[T]) MoveToFront(node *Lode[T]) {
	if list.Head == nil || node == list.Head {
		return
	}
	// 先从当前位置摘除，再放到头部
	list.Detach(node)
	list.PushFrontPartial(node)
}

// MoveToBack moves 'node' to the back.
func (list *List[T]) MoveToBack(node *Lode[T]) {
	if list.Head == nil || node == list.Tail {
		return
	}
	// 先从当前位置摘除，再放到尾部
	list.Detach(node)
	list.PushBackPartial(node)
}
