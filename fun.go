package fun

func Map[T, R any](slice []T, f func(T) R) []R {
	res := make([]R, len(slice))
	for i, elem := range slice {
		res[i] = f(elem)
	}
	return res
}

func FilterMap[T, R any](slice []T, f func(T) (R, bool)) []R {
	res := make([]R, 0, len(slice))
	for _, elem := range slice {
		y, ok := f(elem)
		if ok {
			res = append(res, y)
		}
	}
	return res
}

func ToMap[K comparable, V, T any](slice []T, f func(T) (K, V)) map[K]V {
	res := make(map[K]V, len(slice))
	for _, elem := range slice {
		k, v := f(elem)
		res[k] = v
	}
	return res
}

func ToSlice[K comparable, V, T any](dict map[K]V, f func(K, V) T) []T {
	res := make([]T, 0, len(dict))
	for k, v := range dict {
		res = append(res, f(k, v))
	}
	return res
}

func IterSlice[T any](slice []T, f func(int, T)) {
	for i, elem := range slice {
		f(i, elem)
	}
}

func IterMap[K comparable, V any](dict map[K]V, f func(K, V)) {
	for i, elem := range dict {
		f(i, elem)
	}
}

func Keys[K comparable, V any](dict map[K]V) []K {
	return ToSlice(dict, func(k K, _ V) K { return k })
}

func Zero[T any]() T {
	var zero T
	return zero
}

func Ptr[T any](t T) *T {
	return &t
}

func If[T any](condition bool, ifTrue, ifFalse T) T {
	if condition {
		return ifTrue
	}
	return ifFalse
}
