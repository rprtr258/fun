package fun

import "fmt"

type IthError struct {
	Index int
	Op    string
	Err   error
}

func (e *IthError) Error() string {
	return fmt.Sprintf("index: %d, op: %s, err: %s", e.Index, e.Op, e.Err)
}

func Map[T, R any](slice []T, f func(T) R) []R {
	res := make([]R, len(slice))
	for i, elem := range slice {
		res[i] = f(elem)
	}
	return res
}

func MapErr[T, R any](slice []T, f func(T) (R, error)) ([]R, *IthError) {
	res := make([]R, len(slice))
	for i, elem := range slice {
		y, err := f(elem)
		if err != nil {
			return nil, &IthError{i, "map", err}
		}

		res[i] = y
	}
	return res, nil
}

func Filter[T any](slice []T, f func(T) bool) []T {
	res := make([]T, 0, len(slice))
	for _, elem := range slice {
		if f(elem) {
			res = append(res, elem)
		}
	}
	return res
}

func FilterErr[T any](slice []T, f func(T) (bool, error)) ([]T, *IthError) {
	res := make([]T, len(slice))
	for i, elem := range slice {
		ok, err := f(elem)
		if err != nil {
			return nil, &IthError{i, "filter", err}
		}

		if ok {
			res[i] = elem
		}
	}
	return res, nil
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

func Const[T any](t T) func() T {
	return func() T { return t }
}

func Cond[R any](defaultValue R, cases ...func() (R, bool)) R {
	for _, case_ := range cases {
		if res, ok := case_(); ok {
			return res
		}
	}

	return defaultValue
}

func If[T any](condition bool, ifTrue, ifFalse T) T {
	if condition {
		return ifTrue
	}
	return ifFalse
}

func Has[K comparable, V any](dict map[K]V, key K) bool {
	_, ok := dict[key]
	return ok
}
