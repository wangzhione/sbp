package sets

import (
	"encoding/json"
	"sync"
)

// LockSet is a thread-safe map set.
// 如果你需要去使用, 需要区分对 set 的 read 和 write 操作
//
// r := NewLockSet()
//
// read step
// r.RLock()
// defer r.RUnlock()
// read r.S
//
// write step
// r.Lock()
// defer r.Unlock()
// write r.S
type LockSet[T comparable] struct {
	sync.RWMutex
	S Set[T]
}

func NewLockSet[T comparable](vals ...T) *LockSet[T] {
	return &LockSet[T]{S: New(vals...)}
}

func (r *LockSet[T]) Add(vals ...T) {
	r.Lock()
	defer r.Unlock()
	r.S.Add(vals...)
}

func (r *LockSet[T]) Clear() {
	r.Lock()
	defer r.Unlock()
	clear(r.S)
}

func (r *LockSet[T]) Len() int {
	r.RLock()
	defer r.RUnlock()
	return r.S.Len()
}

func (r *LockSet[T]) Exists(v T) bool {
	r.RLock()
	defer r.RUnlock()
	return r.S.Exists(v)
}

func (r *LockSet[T]) Delete(vals ...T) {
	r.Lock()
	defer r.Unlock()

	for _, key := range vals {
		r.S.Delete(key)
	}
}

func (r *LockSet[T]) Remove(other *LockSet[T]) Set[T] {
	r.Lock()
	defer r.Unlock()

	other.RLock()
	defer other.RUnlock()

	return r.S.Remove(other.S)
}

func (r *LockSet[T]) Equal(other *LockSet[T]) bool {
	r.RLock()
	defer r.RUnlock()

	other.RLock()
	defer other.RUnlock()

	return r.S.Equal(other.S)
}

func (r *LockSet[T]) Clone() (news *LockSet[T]) {
	r.RLock()
	defer r.RUnlock()

	news = &LockSet[T]{
		S: r.S.Clone(),
	}
	return
}

func (r *LockSet[T]) ToSlice() []T {
	r.RLock()
	defer r.RUnlock()
	return r.S.Slice()
}

func (r *LockSet[T]) String() string {
	r.RLock()
	defer r.RUnlock()
	return r.S.String()
}

func (r *LockSet[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.ToSlice())
}

func (r *LockSet[T]) UnmarshalJSON(buf []byte) error {
	var keys []T
	err := json.Unmarshal(buf, &keys)
	if err != nil {
		return err
	}

	r.Lock()
	defer r.Unlock()

	r.S = New(keys...)
	return nil
}
