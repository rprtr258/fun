package fun

import "github.com/rprtr258/fun/exp/zun"

func Map[R, T any, F interface {
	func(T) R | func(T, int) R
}](f F, slice ...T) []R {
	if slice == nil {
		return nil
	}

	res := make([]R, len(slice))
	switch f := any(f).(type) {
	case func(T) R:
		zun.Map(res, f, slice...)
	case func(T, int) R:
		zun.MapI(res, f, slice...)
	default:
		panic("unreachable")
	}
	return res
}

func Filter[T any, F interface {
	func(T) bool | func(T, int) bool
}](f F, slice ...T) []T {
	res := make([]T, 0, len(slice))
	switch f := any(f).(type) {
	case func(T) bool:
		zun.Filter(&res, f, slice...)
	case func(T, int) bool:
		zun.FilterI(&res, f, slice...)
	default:
		panic("unreachable")
	}
	return res
}

func FilterMap[R, T any, F interface {
	func(T) (R, bool) | func(T, int) (R, bool) |
		func(T) Option[R] | func(T, int) Option[R]
}](f F, slice ...T) []R {
	res := make([]R, 0, len(slice))
	switch f := any(f).(type) {
	case func(T) (R, bool):
		zun.FilterMap(&res, f, slice...)
	case func(T, int) (R, bool):
		zun.FilterMapI(&res, f, slice...)
	case func(T) Option[R]:
		zun.FilterMap(&res, func(x T) (R, bool) {
			return f(x).Unpack()
		}, slice...)
	case func(T, int) Option[R]:
		zun.FilterMapI(&res, func(x T, i int) (R, bool) {
			return f(x, i).Unpack()
		}, slice...)
	default:
		panic("unreachable")
	}
	return res
}

func MapDict[T comparable, R any](dict map[T]R, collection ...T) []R {
	result := make([]R, len(collection))
	for i, item := range collection {
		result[i] = dict[item]
	}
	return result
}

func MapErr[R, T any, F interface {
	func(T) (R, error) | func(T, int) (R, error)
}](f F, slice ...T) ([]R, error) {
	res := make([]R, len(slice))
	switch f := any(f).(type) {
	case func(T) (R, error):
		if err := zun.MapErr(res, f, slice...); err != nil {
			return nil, err
		}
	case func(T, int) (R, error):
		if err := zun.MapErrI(res, f, slice...); err != nil {
			return nil, err
		}
	default:
		panic("unreachable")
	}
	return res, nil
}

func Deref[T any](ptr *T) T {
	if ptr == nil {
		return Zero[T]()
	}
	return *ptr
}

func MapToSlice[K comparable, V, R any](dict map[K]V, f func(K, V) R) []R {
	res := make([]R, 0, len(dict))
	zun.MapToSlice(&res, dict, f)
	return res
}

func MapFilterToSlice[K comparable, V, R any](dict map[K]V, f func(K, V) (R, bool)) []R {
	res := make([]R, 0, len(dict))
	zun.MapFilterToSlice(&res, dict, f)
	return res
}

func Keys[K comparable, V any](dict map[K]V) []K {
	res := make([]K, 0, len(dict))
	zun.Keys(&res, dict)
	return res
}

func Values[K comparable, V any](dict map[K]V) []V {
	res := make([]V, 0, len(dict))
	zun.Values(&res, dict)
	return res
}

// FindKeyBy returns the key of the first element predicate returns truthy for.
func FindKeyBy[K comparable, V any](dict map[K]V, predicate func(K, V) bool) (K, bool) {
	return zun.FindKeyBy(dict, predicate)
}

// Uniq returns unique values of slice.
func Uniq[T comparable](collection ...T) []T {
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

// Index returns first found element by predicate along with it's index
func Index[T any](find func(T) bool, slice ...T) (T, int, bool) {
	return zun.Index(find, slice...)
}

// Contains returns true if an element is present in a collection.
func Contains[T comparable](needle T, slice ...T) bool {
	return zun.Contains(needle, slice...)
}

// SliceToMap returns a map containing key-value pairs provided by transform function applied to elements of the given slice.
// If any of two pairs would have the same key the last one gets added to the map.
// The order of keys in returned map is not specified and is not guaranteed to be the same from the original array.
// Alias of Associate().
// Play: https://go.dev/play/p/WHa2CfMO3Lr
func SliceToMap[K comparable, V, T any, F interface {
	func(T) (K, V) | func(T, int) (K, V)
}](f F, slice ...T) map[K]V {
	res := make(map[K]V, len(slice))
	switch f := any(f).(type) {
	case func(T) (K, V):
		zun.SliceToMap(res, f, slice...)
	case func(T, int) (K, V):
		zun.SliceToMapI(res, f, slice...)
	default:
		panic("unreachable")
	}
	return res
}
