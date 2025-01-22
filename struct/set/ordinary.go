package set

import (
	"encoding/json"
	"fmt"
)

// Set[T] map set
type Set[T comparable] map[T]struct{}

// Assert concrete type:Set adheres to ISet interface.
var _ Seter[string] = (*Set[string])(nil)

func NewSet[T comparable]() Set[T] { return make(Set[T]) }

func NewSetWithSize[T comparable](size int) Set[T] { return make(Set[T], size) }

func NewSetWithValue[T comparable](vals ...T) Set[T] {
	s := NewSetWithSize[T](len(vals))
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

func (s Set[T]) AddSet(other Set[T]) {
	for key := range other {
		s[key] = struct{}{}
	}
}

func (s Set[T]) Len() int { return len(s) }

func (s Set[T]) Exists(vals ...T) bool {
	for _, key := range vals {
		if _, ok := s[key]; !ok {
			return false
		}
	}
	return true
}

func (s Set[T]) Contains(v T) bool {
	_, ok := s[v]
	return ok
}

func (s Set[T]) ContainSet(other Set[T]) bool {
	if len(s) < len(other) {
		return false
	}

	for key := range other {
		if _, ok := s[key]; !ok {
			return false
		}
	}
	return true
}

func (s Set[T]) Delete(v T) { delete(s, v) }

func (s Set[T]) Remove(vals ...T) {
	for _, key := range vals {
		delete(s, key)
	}
}

func (s Set[T]) RemoveSet(other Set[T]) Set[T] {
	for key := range other {
		if s.Contains(key) {
			delete(s, key)
		}
	}
	return s
}

func (s Set[T]) EQual(other Set[T]) bool {
	if len(s) != other.Len() {
		return false
	}

	for key := range s {
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

func (s Set[T]) Clone() Set[T] {
	newset := make(Set[T], len(s))
	for elem := range s {
		newset[elem] = struct{}{}
	}
	return newset
}

func (s Set[T]) ToSlice() []T {
	keys := make([]T, 0, s.Len())
	for elem := range s {
		keys = append(keys, elem)
	}
	return keys
}

func (s Set[T]) String() string {
	if len(s) == 0 {
		return "Set{}"
	}
	if len(s) == 1 {
		for elem := range s {
			return fmt.Sprintf("Set{%v}", elem)
		}
	}

	var buf []byte
	buf = append(buf, "Set{"...)
	for elem := range s {
		buf = append(buf, fmt.Sprintf("%v,", elem)...)
	}
	buf[len(buf)-1] = '}'
	return string(buf)
}

func (s Set[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.ToSlice())
}

func (s *Set[T]) UnmarshalJSON(buf []byte) error {
	var keys []T
	err := json.Unmarshal(buf, &keys)
	if err != nil {
		return err
	}

	*s = make(Set[T], len(keys))
	for _, key := range keys {
		(*s)[key] = struct{}{}
	}
	return nil
}
