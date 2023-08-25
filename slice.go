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
	func(T) (R, bool) | func(T, int) (R, bool) |
		func(T) Option[R] | func(T, int) Option[R]
}](slice []T, f F) []R {
	res := []R{}
	switch f := any(f).(type) {
	case func(T) (R, bool):
		for _, x := range slice {
			if y, ok := f(x); ok {
				res = append(res, y)
			}
		}
	case func(T, int) (R, bool):
		for i, x := range slice {
			if y, ok := f(x, i); ok {
				res = append(res, y)
			}
		}
	case func(T) Option[R]:
		for _, x := range slice {
			if y, ok := f(x).Unpack(); ok {
				res = append(res, y)
			}
		}
	case func(T, int) Option[R]:
		for i, x := range slice {
			if y, ok := f(x, i).Unpack(); ok {
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
	res := make([]R, 0, len(dict))
	for k, v := range dict {
		res = append(res, f(k, v))
	}
	return res
}

func MapFilterToSlice[K comparable, V, R any](dict map[K]V, f func(K, V) (R, bool)) []R {
	res := make([]R, 0, len(dict))
	for k, v := range dict {
		y, ok := f(k, v)
		if !ok {
			continue
		}
		res = append(res, y)
	}
	return res
}

func Keys[K comparable, V any](dict map[K]V) []K {
	res := make([]K, 0, len(dict))
	for k := range dict {
		res = append(res, k)
	}
	return res
}

func Values[K comparable, V any](dict map[K]V) []V {
	res := make([]V, 0, len(dict))
	for _, v := range dict {
		res = append(res, v)
	}
	return res
}

// FindKeyBy returns the key of the first element predicate returns truthy for.
func FindKeyBy[K comparable, V any](dict map[K]V, predicate func(K, V) bool) (K, bool) {
	for k, v := range dict {
		if predicate(k, v) {
			return k, true
		}
	}

	return Zero[K](), false
}

// Uniq returns unique values of slice.
func Uniq[T comparable](collection []T) []T {
	res := make([]T, 0, len(collection))
	seen := make(map[T]struct{}, len(collection))
	for _, x := range collection {
		if _, ok := seen[x]; !ok {
			seen[x] = struct{}{}
			res = append(res, x)
		}
	}
	return res
}

// Contains returns true if an element is present in a collection.
func Contains[T comparable](slice []T, needle T) bool {
	for _, x := range slice {
		if x == needle {
			return true
		}
	}

	return false
}

// SliceToMap returns a map containing key-value pairs provided by transform function applied to elements of the given slice.
// If any of two pairs would have the same key the last one gets added to the map.
// The order of keys in returned map is not specified and is not guaranteed to be the same from the original array.
// Alias of Associate().
// Play: https://go.dev/play/p/WHa2CfMO3Lr
func SliceToMap[K comparable, V, T any, F interface {
	func(T) (K, V) | func(T, int) (K, V)
}](slice []T, f F) map[K]V {
	res := make(map[K]V, len(slice))
	switch f := any(f).(type) {
	case func(T) (K, V):
		for _, t := range slice {
			k, v := f(t)
			res[k] = v
		}
	case func(T, int) (K, V):
		for i, t := range slice {
			k, v := f(t, i)
			res[k] = v
		}
	default:
		panic("unreachable")
	}
	return res
}
