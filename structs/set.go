// Package structs sets provides a generic Set implementation for comparable types.
package structs

import (
	"encoding/json"
	"fmt"
	"maps"
	"strings"
)

// Set is a generic set implementation for comparable types.
type Set[T comparable] map[T]struct{}

func NewSet[T comparable](value ...T) Set[T] {
	s := make(Set[T], len(value))
	s.Add(value...)
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

func (s Set[T]) Add(value ...T) {
	for _, key := range value {
		s[key] = struct{}{}
	}
}

func (s Set[T]) Exists(value T) bool {
	_, ok := s[value]
	return ok
}

func (s Set[T]) Len() int { return len(s) }

func (s Set[T]) Delete(value T) { delete(s, value) }

// Remove delete other set & return source set, 主要用于 unit test
func (s Set[T]) Remove(other Set[T]) Set[T] {
	for value := range other {
		delete(s, value)
	}
	return s
}

func (s Set[T]) Equal(other Set[T]) bool {
	if len(s) != other.Len() {
		return false
	}

	for key := range other {
		if !s.Exists(key) {
			return false
		}
	}

	return true
}

func (s Set[T]) Clone() Set[T] {
	return maps.Clone(s)
}

func (s Set[T]) Slice() []T {
	keys := make([]T, 0, len(s)) // 预分配容量
	for elem := range s {
		keys = append(keys, elem)
	}
	return keys
}

func (s Set[T]) String() string {
	n := len(s)
	if n == 0 {
		return "{}"
	}

	if n == 1 {
		for elem := range s {
			return fmt.Sprintf("{%v}", elem)
		}
	}

	// @see https://github.com/golang/go/issues/73189
	var b strings.Builder
	b.WriteByte('{')

	i := 0
	for elem := range s {
		if i > 0 {
			b.WriteByte(',')
		}
		fmt.Fprintf(&b, "%v", elem)
		i++
	}

	b.WriteByte('}')
	return b.String()
}

// MarshalJSON Set[T] Go 支持值接收者方法在指针上调用
func (s Set[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Slice())
}

// UnmarshalJSON Go 的 JSON 解码器 json.Unmarshal 只能调用指针接收者的方法
func (s *Set[T]) UnmarshalJSON(data []byte) (err error) {
	// unmarshal set key slice
	var keys []T
	err = json.Unmarshal(data, &keys)
	if err != nil {
		return
	}

	*s = make(Set[T], len(keys))
	for _, key := range keys {
		(*s)[key] = struct{}{}
	}
	return
}
