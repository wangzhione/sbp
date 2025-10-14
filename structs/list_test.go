package structs

import (
	"testing"
)

// TestNewList 测试创建新链表
func TestNewList(t *testing.T) {
	list := NewList[int]()

	// 验证新创建的链表为空
	if list.Front != nil {
		t.Errorf("新链表的 Front 应该为 nil，实际为 %v", list.Front)
	}
	if list.Back != nil {
		t.Errorf("新链表的 Back 应该为 nil，实际为 %v", list.Back)
	}
}

// TestNewListNode 测试创建新节点
func TestNewListNode(t *testing.T) {
	value := 42
	node := NewListNode(value)

	// 验证节点值正确
	if node.Value != value {
		t.Errorf("节点值应该为 %d，实际为 %v", value, node.Value)
	}

	// 验证新节点的前后指针为 nil
	if node.Prev != nil {
		t.Errorf("新节点的 Prev 应该为 nil，实际为 %v", node.Prev)
	}
	if node.Next != nil {
		t.Errorf("新节点的 Next 应该为 nil，实际为 %v", node.Next)
	}
}

// TestPushBack 测试在链表尾部添加元素
func TestPushBack(t *testing.T) {
	list := NewList[int]()

	// 测试向空链表添加元素
	list.PushBack(1)
	if list.Front == nil || list.Back == nil {
		t.Error("向空链表添加元素后，Front 和 Back 都不应该为 nil")
	}
	if list.Front.Value != 1 || list.Back.Value != 1 {
		t.Errorf("向空链表添加元素后，Front 和 Back 的值都应该为 1，实际 Front=%v, Back=%v",
			list.Front.Value, list.Back.Value)
	}

	// 测试向非空链表添加元素
	list.PushBack(2)
	list.PushBack(3)

	// 验证链表结构
	if list.Front.Value != 1 {
		t.Errorf("Front 值应该为 1，实际为 %v", list.Front.Value)
	}
	if list.Back.Value != 3 {
		t.Errorf("Back 值应该为 3，实际为 %v", list.Back.Value)
	}

	// 验证链表连接关系
	if list.Front.Next.Value != 2 {
		t.Errorf("第二个元素应该为 2，实际为 %v", list.Front.Next.Value)
	}
	if list.Back.Prev.Value != 2 {
		t.Errorf("倒数第二个元素应该为 2，实际为 %v", list.Back.Prev.Value)
	}
}

// TestPushFront 测试在链表头部添加元素
func TestPushFront(t *testing.T) {
	list := NewList[int]()

	// 测试向空链表添加元素
	list.PushFront(1)
	if list.Front == nil || list.Back == nil {
		t.Error("向空链表添加元素后，Front 和 Back 都不应该为 nil")
	}
	if list.Front.Value != 1 || list.Back.Value != 1 {
		t.Errorf("向空链表添加元素后，Front 和 Back 的值都应该为 1，实际 Front=%v, Back=%v",
			list.Front.Value, list.Back.Value)
	}

	// 测试向非空链表添加元素
	list.PushFront(2)
	list.PushFront(3)

	// 验证链表结构
	if list.Front.Value != 3 {
		t.Errorf("Front 值应该为 3，实际为 %v", list.Front.Value)
	}
	if list.Back.Value != 1 {
		t.Errorf("Back 值应该为 1，实际为 %v", list.Back.Value)
	}

	// 验证链表连接关系
	if list.Front.Next.Value != 2 {
		t.Errorf("第二个元素应该为 2，实际为 %v", list.Front.Next.Value)
	}
	if list.Back.Prev.Value != 2 {
		t.Errorf("倒数第二个元素应该为 2，实际为 %v", list.Back.Prev.Value)
	}
}

// TestPushBackNode 测试在链表尾部添加节点
func TestPushBackNode(t *testing.T) {
	list := NewList[int]()

	// 测试添加 nil 节点
	list.PushBackNode(nil)
	if list.Front != nil || list.Back != nil {
		t.Error("添加 nil 节点后，链表应该仍然为空")
	}

	// 测试向空链表添加节点
	node1 := NewListNode(1)
	list.PushBackNode(node1)
	if list.Front != node1 || list.Back != node1 {
		t.Error("向空链表添加节点后，Front 和 Back 都应该指向该节点")
	}

	// 测试向非空链表添加节点
	node2 := NewListNode(2)
	list.PushBackNode(node2)
	if list.Back != node2 {
		t.Error("添加节点后，Back 应该指向新添加的节点")
	}
	if list.Front.Next != node2 {
		t.Error("新添加的节点应该成为第二个节点")
	}
}

// TestPushFrontNode 测试在链表头部添加节点
func TestPushFrontNode(t *testing.T) {
	list := NewList[int]()

	// 测试添加 nil 节点
	list.PushFrontNode(nil)
	if list.Front != nil || list.Back != nil {
		t.Error("添加 nil 节点后，链表应该仍然为空")
	}

	// 测试向空链表添加节点
	node1 := NewListNode(1)
	list.PushFrontNode(node1)
	if list.Front != node1 || list.Back != node1 {
		t.Error("向空链表添加节点后，Front 和 Back 都应该指向该节点")
	}

	// 测试向非空链表添加节点
	node2 := NewListNode(2)
	list.PushFrontNode(node2)
	if list.Front != node2 {
		t.Error("添加节点后，Front 应该指向新添加的节点")
	}
	if list.Back.Prev != node2 {
		t.Error("新添加的节点应该成为第二个节点")
	}
}

// TestRemoveNode 测试从链表中移除节点
func TestRemoveNode(t *testing.T) {
	list := NewList[int]()

	// 测试移除 nil 节点
	list.RemoveNode(nil)
	if list.Front != nil || list.Back != nil {
		t.Error("移除 nil 节点后，链表应该仍然为空")
	}

	// 创建测试数据
	list.PushBack(1)
	list.PushBack(2)
	list.PushBack(3)

	// 测试移除中间节点
	middleNode := list.Front.Next
	list.RemoveNode(middleNode)

	// 验证移除后的链表结构
	if list.Front.Value != 1 {
		t.Errorf("移除中间节点后，Front 值应该为 1，实际为 %v", list.Front.Value)
	}
	if list.Back.Value != 3 {
		t.Errorf("移除中间节点后，Back 值应该为 3，实际为 %v", list.Back.Value)
	}
	if list.Front.Next != list.Back {
		t.Error("移除中间节点后，Front.Next 应该直接指向 Back")
	}
	if list.Back.Prev != list.Front {
		t.Error("移除中间节点后，Back.Prev 应该直接指向 Front")
	}

	// 测试移除头节点
	list.RemoveNode(list.Front)
	if list.Front.Value != 3 {
		t.Errorf("移除头节点后，Front 值应该为 3，实际为 %v", list.Front.Value)
	}
	if list.Back.Value != 3 {
		t.Errorf("移除头节点后，Back 值应该为 3，实际为 %v", list.Back.Value)
	}

	// 测试移除尾节点（也是最后一个节点）
	list.RemoveNode(list.Back)
	if list.Front != nil || list.Back != nil {
		t.Error("移除最后一个节点后，链表应该为空")
	}
}

// TestInsertAfter 测试在指定节点后插入节点
func TestInsertAfter(t *testing.T) {
	list := NewList[int]()

	// 创建测试数据
	list.PushBack(1)
	list.PushBack(3)

	// 在第一个节点后插入新节点
	newNode := NewListNode(2)
	insertedNode := list.InsertAfter(list.Front, newNode)

	// 验证插入的节点
	if insertedNode != newNode {
		t.Error("InsertAfter 应该返回插入的节点")
	}

	// 验证链表结构
	if list.Front.Value != 1 {
		t.Errorf("Front 值应该为 1，实际为 %v", list.Front.Value)
	}
	if list.Front.Next.Value != 2 {
		t.Errorf("第二个节点值应该为 2，实际为 %v", list.Front.Next.Value)
	}
	if list.Back.Value != 3 {
		t.Errorf("Back 值应该为 3，实际为 %v", list.Back.Value)
	}

	// 验证连接关系
	if newNode.Prev != list.Front {
		t.Error("新节点的 Prev 应该指向第一个节点")
	}
	if newNode.Next != list.Back {
		t.Error("新节点的 Next 应该指向最后一个节点")
	}
}

// TestInsertBefore 测试在指定节点前插入节点
func TestInsertBefore(t *testing.T) {
	list := NewList[int]()

	// 创建测试数据
	list.PushBack(1)
	list.PushBack(3)

	// 在第二个节点前插入新节点
	newNode := NewListNode(2)
	insertedNode := list.InsertBefore(list.Back, newNode)

	// 验证插入的节点
	if insertedNode != newNode {
		t.Error("InsertBefore 应该返回插入的节点")
	}

	// 验证链表结构
	if list.Front.Value != 1 {
		t.Errorf("Front 值应该为 1，实际为 %v", list.Front.Value)
	}
	if list.Front.Next.Value != 2 {
		t.Errorf("第二个节点值应该为 2，实际为 %v", list.Front.Next.Value)
	}
	if list.Back.Value != 3 {
		t.Errorf("Back 值应该为 3，实际为 %v", list.Back.Value)
	}

	// 验证连接关系
	if newNode.Prev != list.Front {
		t.Error("新节点的 Prev 应该指向第一个节点")
	}
	if newNode.Next != list.Back {
		t.Error("新节点的 Next 应该指向最后一个节点")
	}
}

// TestComplexOperations 测试复杂操作组合
func TestComplexOperations(t *testing.T) {
	list := NewList[string]()

	// 构建一个复杂的链表
	list.PushBack("A")
	list.PushBack("B")
	list.PushBack("C")
	list.PushFront("0")

	// 验证初始状态
	expected := []string{"0", "A", "B", "C"}
	current := list.Front
	for i, exp := range expected {
		if current == nil {
			t.Errorf("位置 %d 的节点不应该为 nil", i)
			break
		}
		if current.Value != exp {
			t.Errorf("位置 %d 的值应该为 %s，实际为 %s", i, exp, current.Value)
		}
		current = current.Next
	}

	// 测试插入操作
	newNode := NewListNode("X")
	list.InsertAfter(list.Front, newNode)

	// 验证插入后的状态
	expected = []string{"0", "X", "A", "B", "C"}
	current = list.Front
	for i, exp := range expected {
		if current == nil {
			t.Errorf("插入后位置 %d 的节点不应该为 nil", i)
			break
		}
		if current.Value != exp {
			t.Errorf("插入后位置 %d 的值应该为 %s，实际为 %s", i, exp, current.Value)
		}
		current = current.Next
	}

	// 测试移除操作
	list.RemoveNode(newNode)

	// 验证移除后的状态
	expected = []string{"0", "A", "B", "C"}
	current = list.Front
	for i, exp := range expected {
		if current == nil {
			t.Errorf("移除后位置 %d 的节点不应该为 nil", i)
			break
		}
		if current.Value != exp {
			t.Errorf("移除后位置 %d 的值应该为 %s，实际为 %s", i, exp, current.Value)
		}
		current = current.Next
	}
}

// TestEmptyListOperations 测试空链表的操作
func TestEmptyListOperations(t *testing.T) {
	list := NewList[int]()

	// 测试对空链表进行各种操作
	list.PushBack(1)
	if list.Front == nil || list.Back == nil {
		t.Error("向空链表 PushBack 后，Front 和 Back 不应该为 nil")
	}

	// 清空链表
	list.RemoveNode(list.Front)

	// 测试对空链表进行 PushFront
	list.PushFront(2)
	if list.Front == nil || list.Back == nil {
		t.Error("向空链表 PushFront 后，Front 和 Back 不应该为 nil")
	}
}
