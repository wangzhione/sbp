package maps

import (
	"cmp"
	"slices"
)

// Keys returns all keys from a map as a slice
func Keys[K cmp.Ordered, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys) // 对 keys 进行排序
	return keys
}

// RandKeys returns all keys from a map as a slice
func RandKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// RandValues returns all values from a map as a slice
func RandValues[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// Invert swaps keys and values in a map (values must be unique)
func Invert[K comparable, V comparable](m map[K]V) map[V]K {
	inverted := make(map[V]K, len(m))
	for k, v := range m {
		inverted[v] = k
	}
	return inverted
}

// Filter returns a new map containing key-value pairs that satisfy the predicate function
func Filter[K comparable, V any](m map[K]V, predicate func(K, V) bool) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		if predicate(k, v) {
			result[k] = v
		}
	}
	return result
}

// MapValues applies a function to each value and returns a new map with transformed values
func MapValues[K comparable, V any, R any](m map[K]V, mapper func(V) R) map[K]R {
	result := make(map[K]R, len(m))
	for k, v := range m {
		result[k] = mapper(v)
	}
	return result
}
