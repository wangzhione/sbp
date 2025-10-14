package structs

import (
	"testing"
)

// LRUItem 表示LRU缓存中的一个条目
type LRUItem[K comparable, V any] struct {
	Key   K
	Value V
}

// LRUCache 基于双向链表和哈希表实现的LRU缓存
// 复用 list.go 中的 List 和 ListNode 结构
type LRUCache[K comparable, V any] struct {
	capacity int                            // 缓存容量
	list     *List[LRUItem[K, V]]           // 双向链表，用于维护访问顺序
	cache    map[K]*ListNode[LRUItem[K, V]] // 哈希表，用于O(1)查找
}

// NewLRUCache 创建一个新的LRU缓存
func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	if capacity <= 0 {
		panic("LRU缓存容量必须大于0")
	}

	return &LRUCache[K, V]{
		capacity: capacity,
		list:     NewList[LRUItem[K, V]](),
		cache:    make(map[K]*ListNode[LRUItem[K, V]]),
	}
}

// Get 获取缓存中的值，如果存在则将其标记为最近使用
func (lru *LRUCache[K, V]) Get(key K) (V, bool) {
	if node, exists := lru.cache[key]; exists {
		// 将节点移动到头部，标记为最近使用
		lru.list.MoveToFront(node)
		return node.Value.Value, true
	}
	var zero V
	return zero, false
}

// Put 添加或更新缓存中的键值对
func (lru *LRUCache[K, V]) Put(key K, value V) {
	if node, exists := lru.cache[key]; exists {
		// 更新已存在的节点
		node.Value.Value = value
		lru.list.MoveToFront(node)
	} else {
		// 创建新条目
		item := LRUItem[K, V]{Key: key, Value: value}

		if lru.list.Len() >= lru.capacity {
			// 缓存已满，移除最久未使用的节点
			tailNode := lru.list.PopBackNode()
			if tailNode != nil {
				delete(lru.cache, tailNode.Value.Key)
			}
		}

		// 添加新节点到头部
		newNode := NewListNode(item)
		lru.list.PushFrontNode(newNode)
		lru.cache[key] = newNode
	}
}

// Delete 删除缓存中的键值对
func (lru *LRUCache[K, V]) Delete(key K) bool {
	if node, exists := lru.cache[key]; exists {
		lru.list.RemoveNode(node)
		delete(lru.cache, key)
		return true
	}
	return false
}

// Len 返回当前缓存中的元素数量
func (lru *LRUCache[K, V]) Len() int {
	return lru.list.Len()
}

// Cap 返回缓存容量
func (lru *LRUCache[K, V]) Cap() int {
	return lru.capacity
}

// IsEmpty 检查缓存是否为空
func (lru *LRUCache[K, V]) IsEmpty() bool {
	return lru.list.IsEmpty()
}

// IsFull 检查缓存是否已满
func (lru *LRUCache[K, V]) IsFull() bool {
	return lru.list.Len() >= lru.capacity
}

// Clear 清空缓存
func (lru *LRUCache[K, V]) Clear() {
	lru.cache = make(map[K]*ListNode[LRUItem[K, V]])
	lru.list = NewList[LRUItem[K, V]]()
}

// Keys 返回缓存中所有的键（按访问顺序，最近使用的在前）
func (lru *LRUCache[K, V]) Keys() []K {
	keys := make([]K, 0, lru.list.Len())
	current := lru.list.Front
	for current != nil {
		keys = append(keys, current.Value.Key)
		current = current.Next
	}
	return keys
}

// Values 返回缓存中所有的值（按访问顺序，最近使用的在前）
func (lru *LRUCache[K, V]) Values() []V {
	values := make([]V, 0, lru.list.Len())
	current := lru.list.Front
	for current != nil {
		values = append(values, current.Value.Value)
		current = current.Next
	}
	return values
}

// Contains 检查缓存中是否包含指定的键
func (lru *LRUCache[K, V]) Contains(key K) bool {
	_, exists := lru.cache[key]
	return exists
}

// 测试代码

// TestNewLRUCache 测试创建LRU缓存
func TestNewLRUCache(t *testing.T) {
	// 测试正常创建
	lru := NewLRUCache[int, string](3)
	if lru == nil {
		t.Error("LRU缓存不应该为nil")
	}
	if lru.Cap() != 3 {
		t.Errorf("容量应该为3，实际为%d", lru.Cap())
	}
	if lru.Len() != 0 {
		t.Errorf("初始长度应该为0，实际为%d", lru.Len())
	}
	if !lru.IsEmpty() {
		t.Error("新创建的缓存应该为空")
	}

	// 测试容量为0的情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("容量为0时应该panic")
		}
	}()
	NewLRUCache[int, string](0)
}

// TestLRUGet 测试Get操作
func TestLRUGet(t *testing.T) {
	lru := NewLRUCache[int, string](3)

	// 测试获取不存在的键
	value, exists := lru.Get(1)
	if exists {
		t.Error("不存在的键不应该返回true")
	}
	if value != "" {
		t.Errorf("不存在的键应该返回零值，实际为%v", value)
	}

	// 添加一些数据
	lru.Put(1, "one")
	lru.Put(2, "two")
	lru.Put(3, "three")

	// 测试获取存在的键
	value, exists = lru.Get(1)
	if !exists {
		t.Error("存在的键应该返回true")
	}
	if value != "one" {
		t.Errorf("值应该为'one'，实际为%v", value)
	}

	// 验证Get操作会更新访问顺序
	keys := lru.Keys()
	expected := []int{1, 3, 2} // 1被访问后应该移到最前面
	if len(keys) != len(expected) {
		t.Errorf("键的数量应该为%d，实际为%d", len(expected), len(keys))
	}
	for i, key := range keys {
		if key != expected[i] {
			t.Errorf("位置%d的键应该为%d，实际为%d", i, expected[i], key)
		}
	}
}

// TestLRUPut 测试Put操作
func TestLRUPut(t *testing.T) {
	lru := NewLRUCache[int, string](3)

	// 测试添加新键值对
	lru.Put(1, "one")
	if lru.Len() != 1 {
		t.Errorf("长度应该为1，实际为%d", lru.Len())
	}
	if !lru.Contains(1) {
		t.Error("应该包含键1")
	}

	// 测试更新已存在的键
	lru.Put(1, "updated")
	value, exists := lru.Get(1)
	if !exists || value != "updated" {
		t.Errorf("更新后的值应该为'updated'，实际为%v", value)
	}

	// 测试容量限制
	lru.Put(2, "two")
	lru.Put(3, "three")
	lru.Put(4, "four") // 这应该移除键1

	if lru.Len() != 3 {
		t.Errorf("长度应该为3，实际为%d", lru.Len())
	}
	if lru.Contains(1) {
		t.Error("键1应该被移除")
	}
	if !lru.Contains(4) {
		t.Error("应该包含键4")
	}

	// 验证访问顺序
	keys := lru.Keys()
	expected := []int{4, 3, 2}
	for i, key := range keys {
		if key != expected[i] {
			t.Errorf("位置%d的键应该为%d，实际为%d", i, expected[i], key)
		}
	}
}

// TestLRUDelete 测试Delete操作
func TestLRUDelete(t *testing.T) {
	lru := NewLRUCache[int, string](3)

	// 测试删除不存在的键
	if lru.Delete(1) {
		t.Error("删除不存在的键应该返回false")
	}

	// 添加一些数据
	lru.Put(1, "one")
	lru.Put(2, "two")
	lru.Put(3, "three")

	// 测试删除存在的键
	if !lru.Delete(2) {
		t.Error("删除存在的键应该返回true")
	}
	if lru.Len() != 2 {
		t.Errorf("删除后长度应该为2，实际为%d", lru.Len())
	}
	if lru.Contains(2) {
		t.Error("键2应该被删除")
	}

	// 验证剩余键的顺序
	keys := lru.Keys()
	expected := []int{3, 1}
	for i, key := range keys {
		if key != expected[i] {
			t.Errorf("位置%d的键应该为%d，实际为%d", i, expected[i], key)
		}
	}
}

// TestLRUClear 测试Clear操作
func TestLRUClear(t *testing.T) {
	lru := NewLRUCache[int, string](3)

	// 添加一些数据
	lru.Put(1, "one")
	lru.Put(2, "two")
	lru.Put(3, "three")

	// 清空缓存
	lru.Clear()

	if lru.Len() != 0 {
		t.Errorf("清空后长度应该为0，实际为%d", lru.Len())
	}
	if !lru.IsEmpty() {
		t.Error("清空后缓存应该为空")
	}
	if lru.Contains(1) {
		t.Error("清空后不应该包含任何键")
	}
}

// TestLRUAccessOrder 测试访问顺序
func TestLRUAccessOrder(t *testing.T) {
	lru := NewLRUCache[int, string](3)

	// 添加数据
	lru.Put(1, "one")
	lru.Put(2, "two")
	lru.Put(3, "three")

	// 验证初始顺序
	keys := lru.Keys()
	expected := []int{3, 2, 1}
	for i, key := range keys {
		if key != expected[i] {
			t.Errorf("初始位置%d的键应该为%d，实际为%d", i, expected[i], key)
		}
	}

	// 访问键2，应该移到最前面
	lru.Get(2)
	keys = lru.Keys()
	expected = []int{2, 3, 1}
	for i, key := range keys {
		if key != expected[i] {
			t.Errorf("访问后位置%d的键应该为%d，实际为%d", i, expected[i], key)
		}
	}

	// 更新键1，应该移到最前面
	lru.Put(1, "updated")
	keys = lru.Keys()
	expected = []int{1, 2, 3}
	for i, key := range keys {
		if key != expected[i] {
			t.Errorf("更新后位置%d的键应该为%d，实际为%d", i, expected[i], key)
		}
	}
}

// TestLRUCapacity 测试容量限制
func TestLRUCapacity(t *testing.T) {
	lru := NewLRUCache[int, string](2)

	// 添加数据直到超过容量
	lru.Put(1, "one")
	lru.Put(2, "two")
	lru.Put(3, "three") // 应该移除键1

	if lru.Len() != 2 {
		t.Errorf("长度应该为2，实际为%d", lru.Len())
	}
	if lru.Contains(1) {
		t.Error("键1应该被移除")
	}
	if !lru.Contains(2) || !lru.Contains(3) {
		t.Error("键2和3应该存在")
	}

	// 验证最久未使用的键被移除
	keys := lru.Keys()
	expected := []int{3, 2}
	for i, key := range keys {
		if key != expected[i] {
			t.Errorf("位置%d的键应该为%d，实际为%d", i, expected[i], key)
		}
	}
}

// TestLRUEdgeCases 测试边界情况
func TestLRUEdgeCases(t *testing.T) {
	// 测试容量为1的缓存
	lru := NewLRUCache[int, string](1)

	lru.Put(1, "one")
	if !lru.Contains(1) {
		t.Error("应该包含键1")
	}

	lru.Put(2, "two") // 应该移除键1
	if lru.Contains(1) {
		t.Error("键1应该被移除")
	}
	if !lru.Contains(2) {
		t.Error("应该包含键2")
	}

	// 测试空缓存的各种操作
	emptyLRU := NewLRUCache[int, string](3)
	if !emptyLRU.IsEmpty() {
		t.Error("空缓存应该为空")
	}
	if emptyLRU.IsFull() {
		t.Error("空缓存不应该为满")
	}
	if len(emptyLRU.Keys()) != 0 {
		t.Error("空缓存的键列表应该为空")
	}
	if len(emptyLRU.Values()) != 0 {
		t.Error("空缓存的值列表应该为空")
	}
}

// TestLRUDifferentTypes 测试不同数据类型
func TestLRUDifferentTypes(t *testing.T) {
	// 测试字符串键
	stringLRU := NewLRUCache[string, int](2)
	stringLRU.Put("a", 1)
	stringLRU.Put("b", 2)
	value, exists := stringLRU.Get("a")
	if !exists || value != 1 {
		t.Errorf("字符串键测试失败，期望1，实际%v", value)
	}

	// 测试结构体值
	type Person struct {
		Name string
		Age  int
	}
	structLRU := NewLRUCache[int, Person](2)
	person := Person{Name: "Alice", Age: 30}
	structLRU.Put(1, person)
	retrieved, exists := structLRU.Get(1)
	if !exists || retrieved.Name != "Alice" || retrieved.Age != 30 {
		t.Error("结构体值测试失败")
	}
}

// TestLRUConcurrentSafety 测试并发安全性（基础测试）
// 注意：当前实现不是线程安全的，这个测试主要用于检测panic
func TestLRUConcurrentSafety(t *testing.T) {
	// 跳过并发测试，因为当前实现不是线程安全的
	t.Skip("跳过并发测试，当前实现不是线程安全的")
}

// 基准测试

// BenchmarkLRUGet 测试Get操作的性能
func BenchmarkLRUGet(b *testing.B) {
	lru := NewLRUCache[int, int](1000)

	// 预填充数据
	for i := 0; i < 1000; i++ {
		lru.Put(i, i*2)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lru.Get(i % 1000)
	}
}

// BenchmarkLRUPut 测试Put操作的性能
func BenchmarkLRUPut(b *testing.B) {
	lru := NewLRUCache[int, int](1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lru.Put(i%1000, i*2)
	}
}

// BenchmarkLRUDelete 测试Delete操作的性能
func BenchmarkLRUDelete(b *testing.B) {
	lru := NewLRUCache[int, int](1000)

	// 预填充数据
	for i := 0; i < 1000; i++ {
		lru.Put(i, i*2)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lru.Delete(i % 1000)
	}
}

// BenchmarkLRUMixedOperations 测试混合操作的性能
func BenchmarkLRUMixedOperations(b *testing.B) {
	lru := NewLRUCache[int, int](1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := i % 1000
		switch i % 4 {
		case 0:
			lru.Put(key, i)
		case 1:
			lru.Get(key)
		case 2:
			lru.Delete(key)
		case 3:
			lru.Contains(key)
		}
	}
}

// BenchmarkLRULargeCapacity 测试大容量缓存的性能
func BenchmarkLRULargeCapacity(b *testing.B) {
	lru := NewLRUCache[int, int](10000)

	// 预填充数据
	for i := 0; i < 10000; i++ {
		lru.Put(i, i*2)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lru.Get(i % 10000)
	}
}

// BenchmarkLRUSmallCapacity 测试小容量缓存的性能
func BenchmarkLRUSmallCapacity(b *testing.B) {
	lru := NewLRUCache[int, int](10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lru.Put(i, i*2)
	}
}

// 辅助函数：验证LRU缓存的内部状态
func verifyLRUState[K comparable, V any](t *testing.T, lru *LRUCache[K, V], expectedSize int, expectedKeys []K) {
	t.Helper()

	if lru.Len() != expectedSize {
		t.Errorf("LRU长度应该为%d，实际为%d", expectedSize, lru.Len())
	}

	actualKeys := lru.Keys()
	if len(actualKeys) != len(expectedKeys) {
		t.Errorf("键的数量应该为%d，实际为%d", len(expectedKeys), len(actualKeys))
		return
	}

	for i, key := range actualKeys {
		if key != expectedKeys[i] {
			t.Errorf("位置%d的键应该为%v，实际为%v", i, expectedKeys[i], key)
		}
	}
}

// TestLRUStateVerification 测试状态验证辅助函数
func TestLRUStateVerification(t *testing.T) {
	lru := NewLRUCache[int, string](3)

	// 测试空状态
	verifyLRUState(t, lru, 0, []int{})

	// 添加数据并验证
	lru.Put(1, "one")
	verifyLRUState(t, lru, 1, []int{1})

	lru.Put(2, "two")
	verifyLRUState(t, lru, 2, []int{2, 1})

	lru.Put(3, "three")
	verifyLRUState(t, lru, 3, []int{3, 2, 1})

	// 访问键2，验证顺序变化
	lru.Get(2)
	verifyLRUState(t, lru, 3, []int{2, 3, 1})
}
