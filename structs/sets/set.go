package sets

import (
	"encoding/json"
	"fmt"
	"maps"
)

// Set[T] map set
type Set[T comparable] map[T]struct{}

func NewSet[T comparable](vals ...T) Set[T] {
	s := make(Set[T], len(vals))
	for _, elem := range vals {
		s.Add(elem)
	}
	return s
}

/*
	//go:linkname makemap_small
	func makemap_small() *hmap {
		h := new(hmap)
		h.hash0 = uint32(rand())
		return h
	}

	go map 默认返回 *hmap 所以大部分情况, 大部分情况下函数方法值传递更好.
*/

func (s Set[T]) Add(v T) { s[v] = struct{}{} }

func (s Set[T]) Append(vals ...T) {
	for _, key := range vals {
		s[key] = struct{}{}
	}
}

func (s Set[T]) Contains(v T) (ok bool) {
	_, ok = s[v]
	return
}

func (s Set[T]) Len() int { return len(s) }

func (s Set[T]) Delete(v T) { delete(s, v) }

// Remove delete other set & return source set, 主要用于 unit test
func (s Set[T]) Remove(other Set[T]) Set[T] {
	for key := range other {
		s.Delete(key)
	}
	return s
}

func (s Set[T]) Equal(other Set[T]) bool {
	if len(s) != other.Len() {
		return false
	}

	for key := range other {
		if !s.Contains(key) {
			return false
		}
	}
	return true
}

func (s Set[T]) Clone() Set[T] {
	return maps.Clone(s)
}

func (s Set[T]) ToSlice() (keys []T) {
	for elem := range s {
		keys = append(keys, elem)
	}
	return
}

func (s Set[T]) setstring(name string) string {
	if len(s) == 0 {
		return name + "{}"
	}

	if len(s) == 1 {
		for elem := range s {
			return fmt.Sprintf("%s{%v}", name, elem)
		}
	}

	var buf []byte
	buf = append(buf, name...)
	buf = append(buf, '{')
	for elem := range s {
		buf = append(buf, fmt.Sprintf("%v,", elem)...)
	}
	buf[len(buf)-1] = '}'
	return string(buf)
}

func (s Set[T]) String() string {
	return s.setstring("Set")
}

// MarshalJSON Set[T] Go 支持值接收者方法在指针上调用
func (s Set[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.ToSlice())
}

// UnmarshalJSON Go 的 JSON 解码器 json.Unmarshal 只能调用指针接收者的方法
func (s *Set[T]) UnmarshalJSON(buf []byte) (err error) {
	// unmarshal set key slice
	var keys []T
	err = json.Unmarshal(buf, &keys)
	if err != nil {
		return
	}

	*s = make(Set[T], len(keys))
	for _, key := range keys {
		(*s)[key] = struct{}{}
	}
	return
}
