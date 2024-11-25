package various

func MapProject[T any, J comparable, F any](m map[J]T, mapper func(T) F) []F {
	result := make([]F, 0, len(m))
	for _, item := range m {
		result = append(result, mapper(item))
	}
	return result
}

func MapFilter[T any, J comparable](array map[J]T, checker func(T) bool) []J {
	result := make([]J, 0, len(array))
	for key, item := range array {
		if checker(item) {
			result = append(result, key)
		}
	}
	return result
}
