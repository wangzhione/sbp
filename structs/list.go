// Package structs list provides an implementation of a doubly-linked list with a front
// and back. The individual nodes of the list are publicly exposed so that the
// user can have fine-grained control over the list.
package structs

// ListNode is a node in the linked list.
type ListNode[T any] struct {
	Prev, Next *ListNode[T]
	Value      T
}

func NewListNode[T any](value T) *ListNode[T] { return &ListNode[T]{Value: value} }

func (node *ListNode[T]) GetValue() (value T) {
	if node != nil {
		return node.Value
	}
	return
}

// List implements a doubly-linked list.
type List[T any] struct {
	Front, Back *ListNode[T]
}

// List[T]{} or &List[T]{} 都可以, 推荐前者 struct 方式声明

func NewList[T any](values ...T) *List[T] {
	if len(values) == 0 {
		return &List[T]{}
	}
	list := &List[T]{Front: &ListNode[T]{Value: values[0]}}
	list.Back = list.Front
	for i := 1; i < len(values); i++ {
		list.Back.Next = &ListNode[T]{Prev: list.Back, Value: values[i]}
		list.Back = list.Back.Next
	}
	return list
}

// PushBackPartial adds the node 'node' to the back of the list.
func (list *List[T]) PushBackPartial(node *ListNode[T]) {
	if list.Back == nil {
		list.Front = node
		list.Back = node
		return
	}

	node.Prev = list.Back
	list.Back.Next = node
	list.Back = node
}

// PushBackNode adds the node 'node' to the back of the list.
func (list *List[T]) PushBackNode(node *ListNode[T]) {
	if node == nil {
		return
	}

	node.Next = nil
	list.PushBackPartial(node)
}

// PushBack adds 'value' to the end of the list.
func (list *List[T]) PushBack(value T) {
	list.PushBackPartial(NewListNode(value))
}

// PushFrontPartial adds the node 'node' to the front of the list.
func (list *List[T]) PushFrontPartial(node *ListNode[T]) {
	if list.Front == nil {
		list.Front = node
		list.Back = node
		return
	}

	node.Next = list.Front
	list.Front.Prev = node
	list.Front = node
}

// PushFrontNode adds the node 'node' to the front of the list.
func (list *List[T]) PushFrontNode(node *ListNode[T]) {
	if node == nil {
		return
	}

	node.Prev = nil
	list.PushFrontPartial(node)
}

// PushFront adds 'value' to the beginning of the list.
func (list *List[T]) PushFront(value T) {
	list.PushFrontPartial(NewListNode(value))
}

// Detach 从链表中摘除节点 'node'。摘除是个重操作, 所以相对严格一点
func (list *List[T]) Detach(node *ListNode[T]) {
	if node == nil || list.Front == nil {
		return
	}

	if node.Next != nil {
		node.Next.Prev = node.Prev
	} else {
		list.Back = node.Prev
	}

	if node.Prev != nil {
		node.Prev.Next = node.Next
	} else {
		list.Front = node.Next
	}

	node.Next, node.Prev = nil, nil
}

// InsertAfter 在 node 之后插入 next 并返回它。
// next 不应该已在另一个列表中（否则可能破坏另一个列表的结构）
func (list *List[T]) InsertAfter(node *ListNode[T], next *ListNode[T]) {
	if node == nil || list.Front == nil || node == next {
		return
	}

	// 将 next 插入到 node 后面
	next.Prev = node
	next.Next = node.Next

	if node.Next != nil {
		node.Next.Prev = next
	} else {
		// node 是尾节点，更新列表尾指针
		list.Back = next
	}
	node.Next = next
}

// InsertBefore 在 node 之前插入 prev 并返回它。
// 注意：prev 不应该已在另一个列表中（否则可能破坏另一个列表的结构）
func (list *List[T]) InsertBefore(node *ListNode[T], prev *ListNode[T]) {
	if node == nil || list.Front == nil || node == prev {
		return
	}

	// 将 prev 插入到 node 前面
	prev.Next = node
	prev.Prev = node.Prev

	if node.Prev != nil {
		node.Prev.Next = prev
	} else {
		// node 是头节点，更新列表头指针
		list.Front = prev
	}
	node.Prev = prev
}

// Len returns the number of elements (O(n)).
func (list *List[T]) Len() (count int) {
	if list == nil {
		return
	}
	for n := list.Front; n != nil; n = n.Next {
		count++
	}
	return
}

// IsEmpty reports whether the list is empty.
func (list *List[T]) IsEmpty() bool { return list == nil || list.Front == nil }

// PopFront removes and returns the front node.
func (list *List[T]) PopFront() *ListNode[T] {
	node := list.Front
	list.Detach(node)
	return node
}

// PopBack removes and returns the back node.
func (list *List[T]) PopBack() *ListNode[T] {
	node := list.Back
	list.Detach(node)
	return node
}

// MoveToFront moves 'node' to the front.
func (list *List[T]) MoveToFront(node *ListNode[T]) {
	if list.Front == nil || node == list.Front {
		return
	}
	// 先从当前位置摘除，再放到头部
	list.Detach(node)
	list.PushFrontPartial(node)
}

// MoveToBack moves 'node' to the back.
func (list *List[T]) MoveToBack(node *ListNode[T]) {
	if list.Front == nil || node == list.Back {
		return
	}
	// 先从当前位置摘除，再放到尾部
	list.Detach(node)
	list.PushBackPartial(node)
}
