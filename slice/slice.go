// Package slice provides reusable slice utility functions.
package slice

import "github.com/rprtr258/goflow/fun"

// FromMap makes slice of key/value pairs from map.
func FromMap[A comparable, B any](kv map[A]B) []fun.Pair[A, B] {
	kvs := make([]fun.Pair[A, B], 0, len(kv))
	for k, v := range kv {
		kvs = append(kvs, fun.NewPair(k, v))
	}
	return kvs
}

// Reversed makes reversed slice.
func Reversed[A any](xs []A) []A {
	res := make([]A, 0, len(xs))
	copy(res, xs)
	ReverseInplace(res)
	return res
}

// ReverseInplace reverses slice in place.
func ReverseInplace[A any](xs []A) {
	for i, j := 0, len(xs)-1; i < j; i, j = i+1, j-1 {
		xs[i], xs[j] = xs[j], xs[i]
	}
}
