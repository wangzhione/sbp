package sets

import (
	"encoding/json"
	"sync"
)

// SetSaft[T] thread safe map set
// 如果你需要去使用, 需要区分对 set 的 read 和 write 操作
//
// var r *SetSaft[T]
//
// read step
// r.Lock()
// defer r.Unlock()
// read r.S
//
// write step
// r.RLock()
// defer r.RUnlock()
// write r.S
type SetSaft[T comparable] struct {
	sync.RWMutex
	S Set[T]
}

func NewTSet[T comparable]() *SetSaft[T] { return &SetSaft[T]{S: NewSet[T]()} }

func NewTSetWithValue[T comparable](vals ...T) *SetSaft[T] {
	return &SetSaft[T]{S: NewSetWithValue(vals...)}
}

func (r *SetSaft[T]) Add(v T) {
	r.Lock()
	defer r.Unlock()
	r.S.Add(v)
}

func (r *SetSaft[T]) Append(vals ...T) {
	r.Lock()
	defer r.Unlock()
	r.S.Append(vals...)
}

func (r *SetSaft[T]) AddSet(other Set[T]) {
	r.Lock()
	defer r.Unlock()
	r.S.AddSet(other)
}

func (r *SetSaft[T]) Clear() {
	r.Lock()
	defer r.Unlock()
	r.S = NewSet[T]()
}

func (r *SetSaft[T]) Len() int {
	r.RLock()
	defer r.RUnlock()
	return r.S.Len()
}

func (r *SetSaft[T]) Exists(vals ...T) bool {
	r.RLock()
	defer r.RUnlock()
	return r.S.Exists(vals...)
}

func (r *SetSaft[T]) Contains(v T) bool {
	r.RLock()
	defer r.RUnlock()
	return r.S.Contains(v)
}

func (r *SetSaft[T]) ContainSet(other Set[T]) bool {
	r.RLock()
	defer r.RUnlock()
	return r.S.ContainSet(other)
}

func (r *SetSaft[T]) Delete(v T) {
	r.Lock()
	defer r.Unlock()
	r.S.Delete(v)
}

func (r *SetSaft[T]) Remove(vals ...T) {
	r.Lock()
	defer r.Unlock()
	r.S.Remove(vals...)
}

func (r *SetSaft[T]) RemoveSet(other *SetSaft[T]) Set[T] {
	r.Lock()
	defer r.Unlock()

	other.RLock()
	defer other.RUnlock()

	return r.S.RemoveSet(other.S)
}

func (r *SetSaft[T]) EQual(other *SetSaft[T]) bool {
	r.RLock()
	defer r.RUnlock()

	other.RLock()
	defer other.RUnlock()

	return r.S.EQual(other.S)
}

func (r *SetSaft[T]) Clone() Set[T] {
	r.RLock()
	defer r.RUnlock()
	return r.S.Clone()
}

func (r *SetSaft[T]) ToSlice() []T {
	r.RLock()
	defer r.RUnlock()
	return r.S.ToSlice()
}

func (r *SetSaft[T]) String() string {
	r.RLock()
	defer r.RUnlock()
	return r.S.String()
}

func (r *SetSaft[T]) MarshalJSON() ([]byte, error) {
	r.RLock()
	defer r.RUnlock()
	return json.Marshal(r.S.ToSlice())
}

func (r *SetSaft[T]) UnmarshalJSON(buf []byte) error {
	var keys []T
	err := json.Unmarshal(buf, &keys)
	if err != nil {
		return err
	}

	r.Lock()
	defer r.Unlock()

	r.S = make(Set[T], len(keys))
	for _, key := range keys {
		r.S[key] = struct{}{}
	}
	return nil
}
