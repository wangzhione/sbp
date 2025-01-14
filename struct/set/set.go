package set

// Set is the primary interface provided by the mapset package.  It
// represents an unordered set of data and a large number of
// operations that can be applied to that set.
type Set[T comparable] interface {
	// Add adds an element to the set.
	Add(val T)

	// Append multiple elements to the set.
	Append(val ...T)

	// Len returns the number of elements in the set.
	Len() int

	// Contains returns whether the given items
	// are all in the set.
	Contains(val ...T) bool

	// Contain returns whether the given item
	// is in the set.
	//
	// maybe Contains may cause the argument to escape to the heap.
	Contain(val T) bool

	// ContainSet returns whether at least one of the
	// given element are in the set.
	ContainSet(other Set[T]) bool

	// Remove removes a single element from the set.
	Remove(i T)

	// Removes removes multiple elements from the set.
	Removes(i ...T)

	// Each iterates over elements and executes the passed func against each element.
	// If passed func returns true, stop iteration at the time.
	Each(func(T) bool)

	// Equal determines if two sets are equal to each
	// other. If they have the same cardinality
	// and contain the same elements, they are
	// considered equal. The order in which
	// the elements were added is irrelevant.
	//
	// Note that the argument to Equal must be
	// of the same type as the receiver of the
	// method. Otherwise, Equal will panic.
	Equal(other Set[T]) bool

	// Clone returns a clone of the set using the same
	// implementation, duplicating all keys.
	Clone() Set[T]

	// Difference returns the difference between this set
	// and other. The returned set will contain
	// all elements of this set that are not also
	// elements of other.
	//
	// Note that the argument to Difference
	// must be of the same type as the receiver
	// of the method. Otherwise, Difference will
	// panic.
	Difference(other Set[T]) Set[T]

	// Intersect returns a new set containing only the elements
	// that exist only in both sets.
	//
	// Note that the argument to Intersect
	// must be of the same type as the receiver
	// of the method. Otherwise, Intersect will
	// panic.
	Intersect(other Set[T]) Set[T]

	// IsSubSet determines if every element in this set is in
	// the other set.
	//
	// Note that the argument to IsSubSet
	// must be of the same type as the receiver
	// of the method. Otherwise, IsSubSet will
	// panic.
	IsSubSet(other Set[T]) bool

	// SymmetricDifference returns a new set with all elements which are
	// in either this set or the other set but not in both.
	//
	// Note that the argument to SymmetricDifference
	// must be of the same type as the receiver
	// of the method. Otherwise, SymmetricDifference
	// will panic.
	SymmetricDifference(other Set[T]) Set[T]

	// Union returns a new set with all elements in both sets.
	//
	// Note that the argument to Union must be of the
	// same type as the receiver of the method.
	// Otherwise, Union will panic.
	Union(other Set[T]) Set[T]

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
