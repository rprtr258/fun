package fun

func Map[R, T any, F interface {
	func(T) R | func(T, int) R
}](slice []T, f F) []R {
	res := make([]R, len(slice))
	switch f := any(f).(type) {
	case func(T) R:
		for i, x := range slice {
			res[i] = f(x)
		}
	case func(T, int) R:
		for i, x := range slice {
			res[i] = f(x, i)
		}
	default:
		panic("unreachable")
	}
	return res
}

func FilterMap[R, T any, F interface {
	func(T) (R, bool) | func(T, int) (R, bool)
}](slice []T, f F) []R {
	res := []R{}
	switch f := any(f).(type) {
	case func(T) (R, bool):
		for _, x := range slice {
			y, ok := f(x)
			if ok {
				res = append(res, y)
			}
		}
	case func(T, int) (R, bool):
		for i, x := range slice {
			y, ok := f(x, i)
			if ok {
				res = append(res, y)
			}
		}
	default:
		panic("unreachable")
	}
	return res
}

func MapDict[T comparable, R any](collection []T, dict map[T]R) []R {
	result := make([]R, len(collection))

	for i, item := range collection {
		result[i] = dict[item]
	}

	return result
}

func MapErr[R, T any, E interface {
	error
	comparable
}, FE interface {
	func(T) (R, E) | func(T, int) (R, E)
}](slice []T, f FE) ([]R, E) {
	res := make([]R, len(slice))
	switch f := any(f).(type) {
	case func(T) (R, E):
		for i, x := range slice {
			y, err := f(x)
			if err != Zero[E]() {
				return nil, err
			}

			res[i] = y
		}
	case func(T, int) (R, E):
		for i, x := range slice {
			y, err := f(x, i)
			if err != Zero[E]() {
				return nil, err
			}

			res[i] = y
		}
	default:
		panic("unreachable")
	}
	return res, Zero[E]()
}

func Deref[T any](ptr *T) T {
	if ptr == nil {
		return Zero[T]()
	}
	return *ptr
}

func MapToSlice[K comparable, V, R any](dict map[K]V, f func(K, V) R) []R {
	result := make([]R, 0, len(dict))
	for k, v := range dict {
		result = append(result, f(k, v))
	}
	return result
}

func MapFilterToSlice[K comparable, V, R any](dict map[K]V, f func(K, V) (R, bool)) []R {
	result := make([]R, 0, len(dict))
	for k, v := range dict {
		y, ok := f(k, v)
		if !ok {
			continue
		}
		result = append(result, y)
	}
	return result
}
