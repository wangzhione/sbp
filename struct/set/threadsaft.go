package set

import (
	"encoding/json"
	"sync"
)

type TSet[T comparable] struct {
	sync.RWMutex
	S Set[T]
}

// Assert concrete type:TSet adheres to ISet interface.
var _ ISet[string] = (*TSet[string])(nil)

func NewTSet[T comparable]() *TSet[T] { return &TSet[T]{S: NewSet[T]()} }

func NewTSetWithSize[T comparable](size int) *TSet[T] { return &TSet[T]{S: NewSetWithSize[T](size)} }

func (r *TSet[T]) Add(v T) {
	r.Lock()
	defer r.Unlock()
	r.S.Add(v)
}

func (r *TSet[T]) Append(vals ...T) {
	r.Lock()
	defer r.Unlock()
	r.S.Append(vals...)
}

func (r *TSet[T]) AddSet(other Set[T]) {
	r.Lock()
	defer r.Unlock()
	r.S.AddSet(other)
}

func (r *TSet[T]) Len() int {
	r.RLock()
	defer r.RUnlock()
	return r.S.Len()
}

func (r *TSet[T]) Exists(vals ...T) bool {
	r.RLock()
	defer r.RUnlock()
	return r.S.Exists(vals...)
}

func (r *TSet[T]) Contain(v T) bool {
	r.RLock()
	defer r.RUnlock()
	return r.S.Contain(v)
}

func (r *TSet[T]) ContainSet(other Set[T]) bool {
	r.RLock()
	defer r.RUnlock()
	return r.S.ContainSet(other)
}

func (r *TSet[T]) Delete(v T) {
	r.Lock()
	defer r.Unlock()
	r.S.Delete(v)
}

func (r *TSet[T]) Remove(vals ...T) {
	r.Lock()
	defer r.Unlock()
	r.S.Remove(vals...)
}

func (r *TSet[T]) RemoveSet(other Set[T]) Set[T] {
	r.Lock()
	defer r.Unlock()
	return r.S.RemoveSet(other)
}

func (r *TSet[T]) EQual(other Set[T]) bool {
	r.RLock()
	defer r.RUnlock()
	return r.S.EQual(other)
}

func (r *TSet[T]) Clone() Set[T] {
	r.RLock()
	defer r.RUnlock()
	return r.S.Clone()
}

func (r *TSet[T]) ToSlice() []T {
	r.RLock()
	defer r.RUnlock()
	return r.S.ToSlice()
}

func (r *TSet[T]) String() string {
	r.RLock()
	defer r.RUnlock()
	return r.S.String()
}

func (r *TSet[T]) MarshalJSON() ([]byte, error) {
	r.RLock()
	defer r.RUnlock()
	return json.Marshal(r.S.ToSlice())
}

func (r *TSet[T]) UnmarshalJSON(buf []byte) error {
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
