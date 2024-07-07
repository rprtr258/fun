package zun

// len(res) >= len(slice)
func Map[R, T any](
	res []R,
	f func(T) R,
	slice ...T,
) {
	for i, x := range slice {
		res[i] = f(x)
	}
}

// len(res) >= len(slice)
func MapI[R, T any](
	res []R,
	f func(T, int) R,
	slice ...T,
) {
	for i, x := range slice {
		res[i] = f(x, i)
	}
}

// len(res) == 0 && cap(res) >= len(slice)
// returns number of elements
func Filter[T any](
	res []T,
	f func(T) bool,
	slice ...T,
) int {
	for _, x := range slice {
		if f(x) {
			res = append(res, x)
		}
	}
	return len(res)
}

// len(res) == 0 && cap(res) >= len(slice)
// returns number of elements
func FilterI[T any](
	res []T,
	f func(T, int) bool,
	slice ...T,
) int {
	for i, x := range slice {
		if f(x, i) {
			res = append(res, x)
		}
	}
	return len(res)
}

// len(res) == 0 && cap(res) >= len(slice)
// returns number of elements
func FilterMap[R, T any](
	res []R,
	f func(T) (R, bool),
	slice ...T,
) int {
	for _, x := range slice {
		if y, ok := f(x); ok {
			res = append(res, y)
		}
	}
	return len(res)
}

// len(res) == 0 && cap(res) >= len(slice)
// returns number of elements
func FilterMapI[R, T any](
	res []R,
	f func(T, int) (R, bool),
	slice ...T,
) int {
	for i, x := range slice {
		if y, ok := f(x, i); ok {
			res = append(res, y)
		}
	}
	return len(res)
}

// len(res) >= len(slice)
// returns first error
func MapErr[R, T any](
	res []R,
	f func(T) (R, error),
	slice ...T,
) error {
	for i, x := range slice {
		y, err := f(x)
		if err != nil {
			return err
		}

		res[i] = y
	}
	return nil
}

// len(res) >= len(slice)
// returns first error
func MapErrI[R, T any](
	res []R,
	f func(T, int) (R, error),
	slice ...T,
) error {
	for i, x := range slice {
		y, err := f(x, i)
		if err != nil {
			return err
		}

		res[i] = y
	}
	return nil
}

// len(res) == 0 && cap(res) >= len(dict)
func MapToSlice[K comparable, V, R any](
	res []R,
	dict map[K]V,
	f func(K, V) R,
) {
	for k, v := range dict {
		res = append(res, f(k, v))
	}
}

// len(res) == 0 && cap(res) >= len(dict)
// returns number of elements
func MapFilterToSlice[K comparable, V, R any](
	res []R,
	dict map[K]V,
	f func(K, V) (R, bool),
) int {
	for k, v := range dict {
		y, ok := f(k, v)
		if !ok {
			continue
		}
		res = append(res, y)
	}
	return len(res)
}

// len(res) == 0 && cap(res) >= len(dict)
func Keys[K comparable, V any](
	res []K,
	dict map[K]V,
) {
	for k := range dict {
		res = append(res, k)
	}
}

// len(res) == 0 && cap(res) >= len(dict)
func Values[K comparable, V any](
	res []V,
	dict map[K]V,
) {
	for _, v := range dict {
		res = append(res, v)
	}
}

// FindKeyBy returns the key of the first element predicate returns truthy for.
func FindKeyBy[K comparable, V any](
	dict map[K]V,
	predicate func(K, V) bool,
) (K, bool) {
	for k, v := range dict {
		if predicate(k, v) {
			return k, true
		}
	}

	return *new(K), false
}

// Index returns first found element by predicate along with it's index
func Index[T any](
	find func(T) bool,
	slice ...T,
) (T, int, bool) {
	for i, x := range slice {
		if find(x) {
			return x, i, true
		}
	}

	return *new(T), -1, false
}

// Contains returns true if an element is present in a collection.
func Contains[T comparable](needle T, slice ...T) bool {
	_, _, ok := Index(func(x T) bool {
		return x == needle
	}, slice...)
	return ok
}

// cap(res) >= len(slice)
func SliceToMap[K comparable, V, T any](
	res map[K]V,
	f func(T) (K, V),
	slice ...T,
) {
	for _, t := range slice {
		k, v := f(t)
		res[k] = v
	}
}

// cap(res) >= len(slice)
func SliceToMapI[K comparable, V, T any](
	res map[K]V,
	f func(T, int) (K, V),
	slice ...T,
) {
	for i, t := range slice {
		k, v := f(t, i)
		res[k] = v
	}
}

func Reverse[A any](xs []A) {
	for i, j := 0, len(xs)-1; i < j; i, j = i+1, j-1 {
		xs[i], xs[j] = xs[j], xs[i]
	}
}

func Subslice[T any](start, end int, slice ...T) []T {
	if start >= end {
		return nil
	}

	start = max(start, 0)
	end = min(end, len(slice))
	return slice[start:end]
}

func Fold[T, R any](op func(T, R) R, initial R, slice ...T) R {
	res := initial
	for _, elem := range slice {
		res = op(elem, res)
	}
	return res
}
