package arrays

// Unique 去重 slice 并保持原顺序
func Unique[T comparable](a []T) (result []T) {
	seen := make(map[T]struct{}) // 存储已见过的元素

	for _, item := range a {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

// Contains checks if a value exists in a slice
func Contains[T comparable](arr []T, value T) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

// IndexOf returns the index of the first occurrence of value, or -1 if not found
func IndexOf[T comparable](arr []T, value T) int {
	for i, v := range arr {
		if v == value {
			return i
		}
	}
	return -1
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
