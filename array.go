package fun

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/rprtr258/fun/exp/zun"
)

// FromMap makes slice of key/value pairs from map.
func FromMap[A comparable, B any](kv map[A]B) []Pair[A, B] {
	kvs := make([]Pair[A, B], 0, len(kv))
	zun.MapToSlice(&kvs, kv, func(k A, v B) Pair[A, B] {
		return Pair[A, B]{k, v}
	})
	return kvs
}

// Copy slice
func Copy[T any](slice ...T) []T {
	res := make([]T, 0, len(slice))
	copy(res, slice)
	return res
}

// ReverseInplace reverses slice in place.
func ReverseInplace[A any](xs []A) {
	zun.Reverse(xs)
}

// Subslice returns slice from start to end without panicking on out of bounds
func Subslice[T any](start, end int, slice ...T) []T {
	return zun.Subslice(start, end, slice...)
}

// Chunk divides slice into chunks of size chunkSize
func Chunk[T any](chunkSize int, slice ...T) [][]T {
	if chunkSize <= 0 {
		panic(fmt.Errorf("invalid chunkSize: %d", chunkSize))
	}

	res := make([][]T, 0, len(slice)/chunkSize+1)
	for i := 0; i < len(slice); i += chunkSize {
		res = append(res, Subslice(i, i+chunkSize, slice...))
	}
	return res
}

// ConcatMap is like Map but concatenates results
func ConcatMap[T, R any](f func(T) []R, slice ...T) []R {
	res := []R{}
	for _, elem := range slice {
		res = append(res, f(elem)...)
	}
	return res
}

// All returns true if all elements satisfy the condition
func All[T any](condition func(T) bool, slice ...T) bool {
	_, _, ok := zun.Index(func(elem T) bool {
		return !condition(elem)
	}, slice...)
	return !ok
}

// Any returns true if any element satisfies the condition
func Any[T any](condition func(T) bool, slice ...T) bool {
	_, _, ok := zun.Index(func(elem T) bool {
		return condition(elem)
	}, slice...)
	return ok
}

// SortBy sorts slice in place by given function
func SortBy[T any, R cmp.Ordered](by func(T) R, slice ...T) {
	slices.SortFunc(slice, func(i, j T) int {
		return cmp.Compare(by(i), by(j))
	})
}

// GroupBy groups elements by key
func GroupBy[T any, K comparable](by func(T) K, slice ...T) map[K][]T {
	res := map[K][]T{}
	for _, elem := range slice {
		k := by(elem)
		res[k] = append(res[k], elem)
	}
	return res
}
