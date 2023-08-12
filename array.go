package fun

import (
	"fmt"
	"sort"

	"golang.org/x/exp/constraints"
)

func Subslice[T any](slice []T, start, end int) []T {
	if start >= end {
		return nil
	}

	start = Max(start, 0)
	end = Min(end, len(slice))
	return slice[start:end]
}

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

func ConcatMap[T, R any](slice []T, f func(T) []R) []R {
	res := []R{}
	for _, elem := range slice {
		res = append(res, f(elem)...)
	}
	return res
}

func All[T any](slice []T, condition func(T) bool) bool {
	for _, elem := range slice {
		if !condition(elem) {
			return false
		}
	}
	return true
}

func Any[T any](slice []T, condition func(T) bool) bool {
	for _, elem := range slice {
		if condition(elem) {
			return true
		}
	}
	return false
}

func SortBy[T any, R constraints.Ordered](slice []T, by func(T) R) {
	sort.Slice(slice, func(i, j int) bool {
		return by(slice[i]) < by(slice[j])
	})
}

func GroupBy[T any, K comparable](slice []T, by func(T) K) map[K][]T {
	res := map[K][]T{}
	for _, elem := range slice {
		k := by(elem)
		res[k] = append(res[k], elem)
	}
	return res
}
