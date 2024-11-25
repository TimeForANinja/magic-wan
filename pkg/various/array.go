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

func ArrayProject[T any, F any](m []T, mapper func(T) F) []F {
	result := make([]F, 0, len(m))
	for _, item := range m {
		result = append(result, mapper(item))
	}
	return result
}

func ArrayFilter[T any](m []T, checker func(T) bool) []T {
	result := make([]T, 0, len(m))
	for _, item := range m {
		if checker(item) {
			result = append(result, item)
		}
	}
	return result
}

func ArrayEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
