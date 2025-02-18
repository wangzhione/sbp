package array

// RemoveDuplicates 去重 slice 并保持原顺序
func RemoveDuplicates[T comparable](a []T) (result []T) {
	seen := make(map[T]struct{}) // 存储已见过的元素

	for _, item := range a {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}
