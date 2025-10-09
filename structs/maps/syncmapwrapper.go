package maps

import (
	"sync"
)

// Map 是 sync.Map 的 泛型 wrapper 封装，K 必须是可比较类型

type Map[K comparable, V any] struct {
	Data sync.Map
}

func (m *Map[K, V]) Store(key K, value V) {
	m.Data.Store(key, value)
}

func (m *Map[K, V]) Load(key K) (V, bool) {
	v, ok := m.Data.Load(key)
	if !ok {
		var zero V
		return zero, false
	}
	return v.(V), true
}

func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	// loaded == true：
	// 	表示键已存在，函数返回的是已经存在的旧值 v。
	// loaded == false：
	// 	表示键原先不存在，本次调用把你传入的 value 存进去了，并返回它；
	// 	也就是返回的 result 就是你传入的那个 value（对引用类型就是同一指针/引用，对值类型就是那个值的拷贝）。
	v, loaded := m.Data.LoadOrStore(key, value)
	return v.(V), loaded
}

func (m *Map[K, V]) Delete(key K) {
	m.Data.Delete(key)
}

func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.Data.Range(func(k, v any) bool {
		return f(k.(K), v.(V))
	})
}
