package structs

import (
	"testing"
)

// 辅助函数：验证链表结构
func verifyListStructure[T comparable](t *testing.T, list *List[T], expected []T) {
	t.Helper()

	if len(expected) == 0 {
		if list.Front != nil || list.Back != nil {
			t.Error("空链表应该 Front 和 Back 都为 nil")
		}
		return
	}

	// 验证长度
	if list.Len() != len(expected) {
		t.Errorf("链表长度应该为 %d，实际为 %d", len(expected), list.Len())
	}

	// 验证 Front 和 Back
	if list.Front == nil || list.Back == nil {
		t.Error("非空链表不应该有 nil 的 Front 或 Back")
		return
	}

	if list.Front.Value != expected[0] {
		t.Errorf("Front 值应该为 %v，实际为 %v", expected[0], list.Front.Value)
	}

	if list.Back.Value != expected[len(expected)-1] {
		t.Errorf("Back 值应该为 %v，实际为 %v", expected[len(expected)-1], list.Back.Value)
	}

	// 验证正向遍历
	current := list.Front
	for i, exp := range expected {
		if current == nil {
			t.Errorf("位置 %d 的节点不应该为 nil", i)
			break
		}
		if current.Value != exp {
			t.Errorf("位置 %d 的值应该为 %v，实际为 %v", i, exp, current.Value)
		}
		current = current.Next
	}

	// 验证反向遍历
	current = list.Back
	for i := len(expected) - 1; i >= 0; i-- {
		if current == nil {
			t.Errorf("反向遍历位置 %d 的节点不应该为 nil", i)
			break
		}
		if current.Value != expected[i] {
			t.Errorf("反向遍历位置 %d 的值应该为 %v，实际为 %v", i, expected[i], current.Value)
		}
		current = current.Prev
	}
}

// 辅助函数：创建测试链表
func createTestList[T any](values ...T) *List[T] {
	list := NewList[T]()
	for _, value := range values {
		list.PushBack(value)
	}
	return list
}

// 辅助函数：获取链表所有值
func getListValues[T any](list *List[T]) []T {
	var values []T
	for current := list.Front; current != nil; current = current.Next {
		values = append(values, current.Value)
	}
	return values
}

// TestNewList 测试创建新链表
func TestNewList(t *testing.T) {
	list := NewList[int]()

	// 验证新创建的链表为空
	verifyListStructure(t, list, []int{})
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
	verifyListStructure(t, list, []int{1})

	// 测试向非空链表添加元素
	list.PushBack(2)
	list.PushBack(3)
	verifyListStructure(t, list, []int{1, 2, 3})
}

// TestPushFront 测试在链表头部添加元素
func TestPushFront(t *testing.T) {
	list := NewList[int]()

	// 测试向空链表添加元素
	list.PushFront(1)
	verifyListStructure(t, list, []int{1})

	// 测试向非空链表添加元素
	list.PushFront(2)
	list.PushFront(3)
	verifyListStructure(t, list, []int{3, 2, 1})
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
	list.Detach(nil)
	if list.Front != nil || list.Back != nil {
		t.Error("移除 nil 节点后，链表应该仍然为空")
	}

	// 创建测试数据
	list.PushBack(1)
	list.PushBack(2)
	list.PushBack(3)

	// 测试移除中间节点
	middleNode := list.Front.Next
	list.Detach(middleNode)

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
	list.Detach(list.Front)
	if list.Front.Value != 3 {
		t.Errorf("移除头节点后，Front 值应该为 3，实际为 %v", list.Front.Value)
	}
	if list.Back.Value != 3 {
		t.Errorf("移除头节点后，Back 值应该为 3，实际为 %v", list.Back.Value)
	}

	// 测试移除尾节点（也是最后一个节点）
	list.Detach(list.Back)
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
	list.InsertAfter(list.Front, newNode)

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
	list.InsertBefore(list.Back, newNode)

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
	list.Detach(newNode)

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
	list.Detach(list.Front)

	// 测试对空链表进行 PushFront
	list.PushFront(2)
	if list.Front == nil || list.Back == nil {
		t.Error("向空链表 PushFront 后，Front 和 Back 不应该为 nil")
	}
}

// TestLen 测试获取链表长度
func TestLen(t *testing.T) {
	list := NewList[int]()

	// 测试空链表长度
	if list.Len() != 0 {
		t.Errorf("空链表长度应该为 0，实际为 %d", list.Len())
	}

	// 逐步添加元素并测试长度
	for i := 1; i <= 5; i++ {
		list.PushBack(i)
		if list.Len() != i {
			t.Errorf("添加 %d 个元素后，长度应该为 %d，实际为 %d", i, i, list.Len())
		}
	}

	// 逐步删除元素并测试长度
	for i := 4; i >= 0; i-- {
		list.Detach(list.Front)
		if list.Len() != i {
			t.Errorf("删除元素后，长度应该为 %d，实际为 %d", i, list.Len())
		}
	}
}

// TestIsEmpty 测试检查链表是否为空
func TestIsEmpty(t *testing.T) {
	list := NewList[int]()

	// 测试空链表
	if !list.IsEmpty() {
		t.Error("新创建的链表应该为空")
	}

	// 添加元素后测试
	list.PushBack(1)
	if list.IsEmpty() {
		t.Error("添加元素后，链表不应该为空")
	}

	// 删除元素后测试
	list.Detach(list.Front)
	if !list.IsEmpty() {
		t.Error("删除所有元素后，链表应该为空")
	}
}

// TestPopFront 测试弹出并返回首元素
func TestPopFront(t *testing.T) {
	list := NewList[int]()

	// 测试空链表
	_, ok := list.PopFront()
	if ok {
		t.Error("空链表不应该能弹出元素")
	}

	// 创建测试数据
	list.PushBack(1)
	list.PushBack(2)
	list.PushBack(3)

	// 测试弹出首元素
	value, ok := list.PopFront()
	if !ok {
		t.Error("应该能成功弹出首元素")
	}
	if value != 1 {
		t.Errorf("弹出的值应该为 1，实际为 %d", value)
	}

	// 验证弹出后的链表结构
	if list.Front.Value != 2 {
		t.Errorf("弹出后 Front 值应该为 2，实际为 %d", list.Front.Value)
	}
	if list.Back.Value != 3 {
		t.Errorf("弹出后 Back 值应该为 3，实际为 %d", list.Back.Value)
	}

	// 测试弹出所有元素
	list.PopFront()
	value, ok = list.PopFront()
	if !ok || value != 3 {
		t.Errorf("最后一个元素应该为 3，实际为 %d", value)
	}

	// 验证链表为空
	if !list.IsEmpty() {
		t.Error("弹出所有元素后，链表应该为空")
	}
}

// TestPopBack 测试弹出并返回尾元素
func TestPopBack(t *testing.T) {
	list := NewList[int]()

	// 测试空链表
	_, ok := list.PopBack()
	if ok {
		t.Error("空链表不应该能弹出元素")
	}

	// 创建测试数据
	list.PushBack(1)
	list.PushBack(2)
	list.PushBack(3)

	// 测试弹出尾元素
	value, ok := list.PopBack()
	if !ok {
		t.Error("应该能成功弹出尾元素")
	}
	if value != 3 {
		t.Errorf("弹出的值应该为 3，实际为 %d", value)
	}

	// 验证弹出后的链表结构
	if list.Front.Value != 1 {
		t.Errorf("弹出后 Front 值应该为 1，实际为 %d", list.Front.Value)
	}
	if list.Back.Value != 2 {
		t.Errorf("弹出后 Back 值应该为 2，实际为 %d", list.Back.Value)
	}

	// 测试弹出所有元素
	list.PopBack()
	value, ok = list.PopBack()
	if !ok || value != 1 {
		t.Errorf("最后一个元素应该为 1，实际为 %d", value)
	}

	// 验证链表为空
	if !list.IsEmpty() {
		t.Error("弹出所有元素后，链表应该为空")
	}
}

// TestPopFrontNode 测试弹出并返回首节点
func TestPopFrontNode(t *testing.T) {
	list := NewList[int]()

	// 测试空链表
	node := list.PopFrontNode()
	if node != nil {
		t.Error("空链表不应该能弹出节点")
	}

	// 创建测试数据
	list.PushBack(1)
	list.PushBack(2)
	list.PushBack(3)

	// 测试弹出首节点
	node = list.PopFrontNode()
	if node == nil {
		t.Error("应该能成功弹出首节点")
		return
	}
	if node.Value != 1 {
		t.Errorf("弹出节点的值应该为 1，实际为 %d", node.Value)
	}

	// 验证节点已完全分离
	if node.Next != nil || node.Prev != nil {
		t.Error("弹出的节点应该完全分离")
	}

	// 验证弹出后的链表结构
	if list.Front.Value != 2 {
		t.Errorf("弹出后 Front 值应该为 2，实际为 %d", list.Front.Value)
	}
	if list.Back.Value != 3 {
		t.Errorf("弹出后 Back 值应该为 3，实际为 %d", list.Back.Value)
	}
}

// TestPopBackNode 测试弹出并返回尾节点
func TestPopBackNode(t *testing.T) {
	list := NewList[int]()

	// 测试空链表
	node := list.PopBackNode()
	if node != nil {
		t.Error("空链表不应该能弹出节点")
	}

	// 创建测试数据
	list.PushBack(1)
	list.PushBack(2)
	list.PushBack(3)

	// 测试弹出尾节点
	node = list.PopBackNode()
	if node == nil {
		t.Error("应该能成功弹出尾节点")
		return
	}
	if node.Value != 3 {
		t.Errorf("弹出节点的值应该为 3，实际为 %d", node.Value)
	}

	// 验证节点已完全分离
	if node.Next != nil || node.Prev != nil {
		t.Error("弹出的节点应该完全分离")
	}

	// 验证弹出后的链表结构
	if list.Front.Value != 1 {
		t.Errorf("弹出后 Front 值应该为 1，实际为 %d", list.Front.Value)
	}
	if list.Back.Value != 2 {
		t.Errorf("弹出后 Back 值应该为 2，实际为 %d", list.Back.Value)
	}
}

// TestMoveToFront 测试移动节点到链表头部
func TestMoveToFront(t *testing.T) {
	list := NewList[int]()

	// 测试空链表
	list.MoveToFront(nil)
	if list.Front != nil || list.Back != nil {
		t.Error("空链表操作后应该仍然为空")
	}

	// 创建测试数据
	list.PushBack(1)
	list.PushBack(2)
	list.PushBack(3)
	list.PushBack(4)

	// 测试移动尾节点到头部
	tailNode := list.Back
	list.MoveToFront(tailNode)

	// 验证移动后的链表结构 [4, 1, 2, 3]
	expected := []int{4, 1, 2, 3}
	current := list.Front
	for i, exp := range expected {
		if current == nil {
			t.Errorf("位置 %d 的节点不应该为 nil", i)
			break
		}
		if current.Value != exp {
			t.Errorf("位置 %d 的值应该为 %d，实际为 %d", i, exp, current.Value)
		}
		current = current.Next
	}

	// 验证头尾指针
	if list.Front.Value != 4 {
		t.Errorf("Front 值应该为 4，实际为 %d", list.Front.Value)
	}
	if list.Back.Value != 3 {
		t.Errorf("Back 值应该为 3，实际为 %d", list.Back.Value)
	}

	// 测试移动已经是头部的节点
	headNode := list.Front
	list.MoveToFront(headNode)
	if list.Front != headNode {
		t.Error("移动头部节点后，Front 应该仍然指向该节点")
	}
}

// TestMoveToBack 测试移动节点到链表尾部
func TestMoveToBack(t *testing.T) {
	list := NewList[int]()

	// 测试空链表
	list.MoveToBack(nil)
	if list.Front != nil || list.Back != nil {
		t.Error("空链表操作后应该仍然为空")
	}

	// 创建测试数据
	list.PushBack(1)
	list.PushBack(2)
	list.PushBack(3)
	list.PushBack(4)

	// 测试移动头节点到尾部
	headNode := list.Front
	list.MoveToBack(headNode)

	// 验证移动后的链表结构 [2, 3, 4, 1]
	expected := []int{2, 3, 4, 1}
	current := list.Front
	for i, exp := range expected {
		if current == nil {
			t.Errorf("位置 %d 的节点不应该为 nil", i)
			break
		}
		if current.Value != exp {
			t.Errorf("位置 %d 的值应该为 %d，实际为 %d", i, exp, current.Value)
		}
		current = current.Next
	}

	// 验证头尾指针
	if list.Front.Value != 2 {
		t.Errorf("Front 值应该为 2，实际为 %d", list.Front.Value)
	}
	if list.Back.Value != 1 {
		t.Errorf("Back 值应该为 1，实际为 %d", list.Back.Value)
	}

	// 测试移动已经是尾部的节点
	tailNode := list.Back
	list.MoveToBack(tailNode)
	if list.Back != tailNode {
		t.Error("移动尾部节点后，Back 应该仍然指向该节点")
	}
}

// TestEdgeCases 测试边界情况和错误处理
func TestEdgeCases(t *testing.T) {
	list := NewList[int]()

	// 测试 InsertAfter 的边界情况
	// 1. 插入 nil 节点
	list.InsertAfter(nil, nil) // 应该不执行任何操作

	// 2. 向空链表插入
	list.PushBack(1)
	list.InsertAfter(nil, NewListNode(2)) // 应该不执行任何操作

	// 3. 插入自己
	node := list.Front
	list.InsertAfter(node, node) // 应该不执行任何操作

	// 测试 InsertBefore 的边界情况
	// 1. 插入 nil 节点
	list.InsertBefore(nil, nil) // 应该不执行任何操作

	// 2. 向空链表插入
	list.Detach(list.Front)
	list.InsertBefore(nil, NewListNode(2)) // 应该不执行任何操作

	// 3. 插入自己
	list.PushBack(1)
	node = list.Front
	list.InsertBefore(node, node) // 应该不执行任何操作

	// 测试 RemoveNode 的边界情况
	// 1. 移除 nil 节点
	list.Detach(nil)

	// 2. 从空链表移除
	list.Detach(list.Front)
	list.Detach(NewListNode(99))

	// 测试 MoveToFront 的边界情况
	// 1. 移动 nil 节点
	list.MoveToFront(nil)

	// 2. 从空链表移动
	list.MoveToFront(NewListNode(99))

	// 测试 MoveToBack 的边界情况
	// 1. 移动 nil 节点
	list.MoveToBack(nil)

	// 2. 从空链表移动
	list.MoveToBack(NewListNode(99))
}

// TestComplexScenarios 测试复杂场景
func TestComplexScenarios(t *testing.T) {
	list := NewList[string]()

	// 场景1：频繁的插入和删除操作
	// 构建链表: A -> B -> C -> D
	list.PushBack("A")
	list.PushBack("B")
	list.PushBack("C")
	list.PushBack("D")

	// 在 B 后插入 X: A -> B -> X -> C -> D
	newNode := NewListNode("X")
	list.InsertAfter(list.Front.Next, newNode)

	// 在 C 前插入 Y: A -> B -> X -> Y -> C -> D
	newNode2 := NewListNode("Y")
	list.InsertBefore(list.Back.Prev, newNode2)

	// 验证最终结构
	expected := []string{"A", "B", "X", "Y", "C", "D"}
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

	// 场景2：移动操作组合
	// 移动 Y 到头部: Y -> A -> B -> X -> C -> D
	list.MoveToFront(newNode2)
	if list.Front.Value != "Y" {
		t.Errorf("移动后 Front 应该为 Y，实际为 %s", list.Front.Value)
	}

	// 移动 X 到尾部: Y -> A -> B -> C -> D -> X
	list.MoveToBack(newNode)
	if list.Back.Value != "X" {
		t.Errorf("移动后 Back 应该为 X，实际为 %s", list.Back.Value)
	}

	// 场景3：手动删除特定元素
	// 删除所有元音字母 A 和 Y
	current = list.Front
	for current != nil {
		next := current.Next
		if current.Value == "A" || current.Value == "Y" {
			list.Detach(current)
		}
		current = next
	}

	// 验证删除后的结构: B -> C -> D -> X
	expected = []string{"B", "C", "D", "X"}
	current = list.Front
	for i, exp := range expected {
		if current == nil {
			t.Errorf("删除后位置 %d 的节点不应该为 nil", i)
			break
		}
		if current.Value != exp {
			t.Errorf("删除后位置 %d 的值应该为 %s，实际为 %s", i, exp, current.Value)
		}
		current = current.Next
	}
}

// TestConcurrentSafety 测试并发安全性（基础测试）
func TestConcurrentSafety(t *testing.T) {
	// 注意：这个测试只是基础测试，真正的并发测试需要更复杂的设置
	list := NewList[int]()

	// 测试在并发环境下的基本操作不会panic
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("协程 %d 发生 panic: %v", id, r)
				}
				done <- true
			}()

			// 执行一些基本操作
			for j := 0; j < 100; j++ {
				list.PushBack(id*100 + j)
				if list.Len() > 0 {
					list.PopFront()
				}
			}
		}(i)
	}

	// 等待所有协程完成
	for i := 0; i < 10; i++ {
		<-done
	}
}

// 基准测试

// BenchmarkPushBack 测试 PushBack 操作的性能
func BenchmarkPushBack(b *testing.B) {
	list := NewList[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.PushBack(i)
	}
}

// BenchmarkPushFront 测试 PushFront 操作的性能
func BenchmarkPushFront(b *testing.B) {
	list := NewList[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.PushFront(i)
	}
}

// BenchmarkPopFront 测试 PopFront 操作的性能
func BenchmarkPopFront(b *testing.B) {
	list := NewList[int]()
	// 预填充数据
	for i := 0; i < b.N; i++ {
		list.PushBack(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.PopFront()
	}
}

// BenchmarkPopBack 测试 PopBack 操作的性能
func BenchmarkPopBack(b *testing.B) {
	list := NewList[int]()
	// 预填充数据
	for i := 0; i < b.N; i++ {
		list.PushBack(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.PopBack()
	}
}

// BenchmarkLen 测试 Len 操作的性能
func BenchmarkLen(b *testing.B) {
	list := NewList[int]()
	// 预填充数据
	for i := 0; i < 1000; i++ {
		list.PushBack(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.Len()
	}
}

// BenchmarkMoveToFront 测试 MoveToFront 操作的性能
func BenchmarkMoveToFront(b *testing.B) {
	list := NewList[int]()
	// 预填充数据
	for i := 0; i < 1000; i++ {
		list.PushBack(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 移动尾节点到头部
		if list.Back != nil {
			list.MoveToFront(list.Back)
		}
	}
}

// BenchmarkMoveToBack 测试 MoveToBack 操作的性能
func BenchmarkMoveToBack(b *testing.B) {
	list := NewList[int]()
	// 预填充数据
	for i := 0; i < 1000; i++ {
		list.PushBack(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 移动头节点到尾部
		if list.Front != nil {
			list.MoveToBack(list.Front)
		}
	}
}

// BenchmarkInsertAfter 测试 InsertAfter 操作的性能
func BenchmarkInsertAfter(b *testing.B) {
	list := NewList[int]()
	// 预填充数据
	for i := 0; i < 1000; i++ {
		list.PushBack(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 在中间位置插入
		if list.Front != nil && list.Front.Next != nil {
			newNode := NewListNode(i)
			list.InsertAfter(list.Front, newNode)
		}
	}
}

// BenchmarkInsertBefore 测试 InsertBefore 操作的性能
func BenchmarkInsertBefore(b *testing.B) {
	list := NewList[int]()
	// 预填充数据
	for i := 0; i < 1000; i++ {
		list.PushBack(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 在中间位置插入
		if list.Back != nil && list.Back.Prev != nil {
			newNode := NewListNode(i)
			list.InsertBefore(list.Back, newNode)
		}
	}
}

// BenchmarkMixedOperations 测试混合操作的性能
func BenchmarkMixedOperations(b *testing.B) {
	list := NewList[int]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 混合各种操作
		list.PushBack(i)
		if i%2 == 0 {
			list.PushFront(i * 2)
		}
		if i%3 == 0 && !list.IsEmpty() {
			list.PopFront()
		}
		if i%5 == 0 && !list.IsEmpty() {
			list.PopBack()
		}
		if i%7 == 0 && list.Len() > 10 {
			// 手动删除一些元素
			current := list.Front
			for current != nil && list.Len() > 5 {
				next := current.Next
				if current.Value%10 == 0 {
					list.Detach(current)
				}
				current = next
			}
		}
	}
}

// TestTableDrivenTests 表驱动测试，测试多种数据类型
func TestTableDrivenTests(t *testing.T) {
	tests := []struct {
		name     string
		values   []interface{}
		expected []interface{}
	}{
		{
			name:     "整数类型",
			values:   []interface{}{1, 2, 3, 4, 5},
			expected: []interface{}{1, 2, 3, 4, 5},
		},
		{
			name:     "字符串类型",
			values:   []interface{}{"hello", "world", "test"},
			expected: []interface{}{"hello", "world", "test"},
		},
		{
			name:     "浮点数类型",
			values:   []interface{}{1.1, 2.2, 3.3},
			expected: []interface{}{1.1, 2.2, 3.3},
		},
		{
			name:     "布尔类型",
			values:   []interface{}{true, false, true},
			expected: []interface{}{true, false, true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewList[interface{}]()
			for _, value := range tt.values {
				list.PushBack(value)
			}

			// 验证链表结构
			if list.Len() != len(tt.expected) {
				t.Errorf("链表长度应该为 %d，实际为 %d", len(tt.expected), list.Len())
			}

			// 验证值
			current := list.Front
			for i, exp := range tt.expected {
				if current == nil {
					t.Errorf("位置 %d 的节点不应该为 nil", i)
					break
				}
				if current.Value != exp {
					t.Errorf("位置 %d 的值应该为 %v，实际为 %v", i, exp, current.Value)
				}
				current = current.Next
			}
		})
	}
}

// TestMemoryLeaks 测试内存泄漏（基础测试）
func TestMemoryLeaks(t *testing.T) {
	list := NewList[int]()

	// 大量添加和删除操作
	for i := 0; i < 10000; i++ {
		list.PushBack(i)
		if i%2 == 0 {
			list.PopFront()
		}
	}

	// 清空链表
	for !list.IsEmpty() {
		list.PopFront()
	}

	// 验证链表为空
	verifyListStructure(t, list, []int{})
}

// TestStressTest 压力测试
func TestStressTest(t *testing.T) {
	list := NewList[int]()

	// 大量操作
	for i := 0; i < 1000; i++ {
		// 随机操作
		switch i % 6 {
		case 0:
			list.PushBack(i)
		case 1:
			list.PushFront(i)
		case 2:
			if !list.IsEmpty() {
				list.PopFront()
			}
		case 3:
			if !list.IsEmpty() {
				list.PopBack()
			}
		case 4:
			if !list.IsEmpty() && list.Len() > 1 {
				// 手动删除一些元素
				current := list.Front
				for current != nil && list.Len() > 1 {
					next := current.Next
					if current.Value%10 == 0 {
						list.Detach(current)
					}
					current = next
				}
			}
		case 5:
			if !list.IsEmpty() && list.Len() > 1 {
				if list.Back != nil {
					list.MoveToFront(list.Back)
				}
			}
		}

		// 每100次操作验证一次链表结构
		if i%100 == 0 {
			// 验证链表没有循环引用
			visited := make(map[*ListNode[int]]bool)
			current := list.Front
			for current != nil {
				if visited[current] {
					t.Error("检测到循环引用")
					return
				}
				visited[current] = true
				current = current.Next
			}
		}
	}
}

// TestIteratorPattern 测试迭代器模式
func TestIteratorPattern(t *testing.T) {
	list := createTestList(1, 2, 3, 4, 5)

	// 正向迭代
	expected := []int{1, 2, 3, 4, 5}
	actual := getListValues(list)

	if len(actual) != len(expected) {
		t.Errorf("迭代结果长度应该为 %d，实际为 %d", len(expected), len(actual))
	}

	for i, exp := range expected {
		if i >= len(actual) || actual[i] != exp {
			t.Errorf("位置 %d 的值应该为 %d，实际为 %v", i, exp, actual[i])
		}
	}

	// 反向迭代
	reverseExpected := []int{5, 4, 3, 2, 1}
	var reverseActual []int
	for current := list.Back; current != nil; current = current.Prev {
		reverseActual = append(reverseActual, current.Value)
	}

	if len(reverseActual) != len(reverseExpected) {
		t.Errorf("反向迭代结果长度应该为 %d，实际为 %d", len(reverseExpected), len(reverseActual))
	}

	for i, exp := range reverseExpected {
		if i >= len(reverseActual) || reverseActual[i] != exp {
			t.Errorf("反向位置 %d 的值应该为 %d，实际为 %v", i, exp, reverseActual[i])
		}
	}
}
