// Package structs provides utility functions for working with pointers and ordered values.
package structs

import "cmp"

// Ptr returns a pointer to the provided value.
//
// const No = "9527"
// structs.Ptr(BucketNo)
//
// 对于 var Oh string , 更方便走 &Oh
func Ptr[T any](v T) *T {
	return &v
}

func Max[T cmp.Ordered](vals ...T) (maxval T) {
	if len(vals) == 0 {
		return
	}

	maxval = vals[0]
	for i := 1; i < len(vals); i++ {
		if maxval < vals[i] {
			maxval = vals[i]
		}
	}
	return
}

func Min[T cmp.Ordered](vals ...T) (minval T) {
	if len(vals) == 0 {
		return
	}

	minval = vals[0]
	for i := 1; i < len(vals); i++ {
		if minval > vals[i] {
			minval = vals[i]
		}
	}
	return
}

func Ternary[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}

// Coalesce 返回第一个非零值（需要可比较）
func Coalesce[T comparable](v, def T) (zero T) {
	if v == zero {
		return def
	}
	return v
}
