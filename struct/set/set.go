package set

import (
	"cmp"
	"slices"
)

// ISet is the primary interface provided by the mapset package.  It
// represents an unordered set of data and a large number of
// operations that can be applied to that set.
type ISet[T comparable] interface {
	// Add adds an element to the set.
	Add(v T)

	// Append multiple elements to the set.
	Append(vals ...T)

	// Len returns the number of elements in the set.
	Len() int

	// Exists returns whether the given items
	// are all in the set.
	Exists(vals ...T) bool

	// Contain returns whether the given item
	// is in the set.
	//
	// maybe Exists may cause the argument to escape to the heap.
	Contain(v T) bool

	// Delete removes a single element from the set.
	Delete(v T)

	// Remove removes multiple elements from the set.
	Remove(vals ...T)

	// ToSlice returns the members of the set as a slice.
	ToSlice() []T

	// String provides a convenient string representation
	// of the current state of the set.
	String() string

	// MarshalJSON will marshal the set into a JSON-based representation.
	MarshalJSON() ([]byte, error)

	// UnmarshalJSON will unmarshal a JSON-based byte slice into a full Set datastructure.
	// For this to work, set subtypes must implemented the Marshal/Unmarshal interface.
	UnmarshalJSON(b []byte) error
}

// Sorted returns a sorted slice of a set of any ordered type in ascending order.
// When sorting floating-point numbers, NaNs are ordered before other values.
func Sorted[T cmp.Ordered](s ISet[T]) []T {
	keys := s.ToSlice()
	slices.Sort(keys)
	return keys
}

func NewSetFromSlice[T comparable](keys []T) Set[T] {
	s := NewSetWithSize[T](len(keys))

	for _, key := range keys {
		s.Add(key)
	}

	return s
}

// NewSetFromMapKey creates and returns a new set with the given keys of the map.
// Operations on the resulting set are not thread-safe.
func NewSetFromMapKey[T comparable, V any](m map[T]V) Set[T] {
	s := NewSetWithSize[T](len(m))

	for key := range m {
		s.Add(key)
	}

	return s
}

func NewSetFromMapValue[T comparable, V comparable](m map[T]V) Set[V] {
	s := NewSetWithSize[V](len(m))

	for _, value := range m {
		s.Add(value)
	}

	return s
}
