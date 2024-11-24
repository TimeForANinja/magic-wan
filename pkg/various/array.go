package various

func ArrayFind[T any](array []*T, checker func(*T) bool) *T {
	for _, item := range array {
		if checker(item) {
			return item
		}
	}
	return nil
}

func ArrayIncludes[T any](array []T, checker func(T) bool) bool {
	for _, item := range array {
		if checker(item) {
			return true
		}
	}
	return false
}

func ReflectMap[T any, J comparable, F any](m map[J]T, mapper func(T) F) []F {
	result := make([]F, 0, len(m))
	for _, item := range m {
		result = append(result, mapper(item))
	}
	return result
}
