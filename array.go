package fun

import (
	"cmp"
	"fmt"
	"slices"
)

// FromMap makes slice of key/value pairs from map.
func FromMap[A comparable, B any](kv map[A]B) []Pair[A, B] {
	kvs := make([]Pair[A, B], 0, len(kv))
	for k, v := range kv {
		kvs = append(kvs, Pair[A, B]{k, v})
	}
	return kvs
}

// Copy slice
func Copy[T any](slice []T) []T {
	res := make([]T, 0, len(slice))
	copy(res, slice)
	return res
}

// ReverseInplace reverses slice in place.
func ReverseInplace[A any](xs []A) {
	for i, j := 0, len(xs)-1; i < j; i, j = i+1, j-1 {
		xs[i], xs[j] = xs[j], xs[i]
	}
}

// Subslice returns slice from start to end without panicking on out of bounds
func Subslice[T any](slice []T, start, end int) []T {
	if start >= end {
		return nil
	}

	start = Max(start, 0)
	end = Min(end, len(slice))
	return slice[start:end]
}

// Chunk divides slice into chunks of size chunkSize
func Chunk[T any](slice []T, chunkSize int) [][]T {
	if chunkSize <= 0 {
		panic(fmt.Errorf("invalid chunkSize: %d", chunkSize))
	}

	res := make([][]T, 0, len(slice)/chunkSize+1)
	for i := 0; i < len(slice); i += chunkSize {
		res = append(res, Subslice(slice, i, i+chunkSize))
	}
	return res
}

// ConcatMap is like Map but concatenates results
func ConcatMap[T, R any](slice []T, f func(T) []R) []R {
	res := []R{}
	for _, elem := range slice {
		res = append(res, f(elem)...)
	}
	return res
}

// All returns true if all elements satisfy the condition
func All[T any](slice []T, condition func(T) bool) bool {
	for _, elem := range slice {
		if !condition(elem) {
			return false
		}
	}
	return true
}

// Any returns true if any element satisfies the condition
func Any[T any](slice []T, condition func(T) bool) bool {
	for _, elem := range slice {
		if condition(elem) {
			return true
		}
	}
	return false
}

// SortBy sorts slice in place by given function
func SortBy[T any, R cmp.Ordered](slice []T, by func(T) R) {
	slices.SortFunc(slice, func(i, j T) int {
		return cmp.Compare(by(i), by(j))
	})
}

// GroupBy groups elements by key
func GroupBy[T any, K comparable](slice []T, by func(T) K) map[K][]T {
	res := map[K][]T{}
	for _, elem := range slice {
		k := by(elem)
		res[k] = append(res[k], elem)
	}
	return res
}
