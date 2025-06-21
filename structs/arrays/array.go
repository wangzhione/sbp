// Package arrays provides utility functions for working with slices and arrays.
package arrays

import "github.com/wangzhione/sbp/structs"

// go std slices 很 strong

// Contains checks if a value exists in a slice
//
// func Contains[S ~[]E, E comparable](s S, v E) bool

// func Clone[S ~[]E, E any](s S) S

// Index returns the index of the first occurrence of value, or -1 if not found
//
// func Index[S ~[]E, E comparable](s S, v E) int
// func IndexFunc[S ~[]E, E any](s S, f func(E) bool) int

// Unique 去重 slice 并保持原顺序
func Unique[T comparable](a []T) (result []T) {
	seen := make(map[T]struct{}) // 存储已见过的元素

	for _, item := range a {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return
}

// Set array to set[]struct{}
func Set[T comparable](a []T) map[T]struct{} {
	seen := make(map[T]struct{}, len(a)) // 存储已见过的元素

	for _, item := range a {
		seen[item] = struct{}{}
	}

	return seen
}

// Filter returns a new slice containing elements that satisfy the predicate function
func Filter[T any](arr []T, predicate func(T) bool) (result []T) {
	for _, v := range arr {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return
}

// SlicePtrs returns a slice of *T from the specified values.
func SlicePtrs[T any](vv ...T) []*T {
	ptrs := make([]*T, len(vv))
	for i := range vv {
		ptrs[i] = structs.Ptr(vv[i])
	}
	return ptrs
}
