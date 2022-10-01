package stream

// functions to make something from Stream that is not Stream.

import (
	"github.com/rprtr258/go-flow/v2/fun"
)

// ForEach invokes a simple function for each element of the stream.
func ForEach[A any](xs Stream[A], f func(A)) {
	for x := range xs {
		f(x)
	}
}

// DrainAll throws away all values.
func DrainAll[A any](xs Stream[A]) {
	ForEach(xs, func(_ A) {})
}

// CollectToSlice executes the stream and collects all results to a slice.
func CollectToSlice[A any](xs Stream[A]) []A {
	slice := make([]A, 0)
	ForEach(xs, func(a A) { slice = append(slice, a) })
	return slice
}

// CollectToSet executes the stream and collects all results to a set.
func CollectToSet[A comparable](xs Stream[A]) fun.Set[A] {
	set := make(fun.Set[A])
	ForEach(xs, func(a A) { set[a] = fun.Unit1 })
	return set
}

// Head takes the first element if present.
func Head[A any](xs Stream[A]) fun.Option[A] {
	return fun.FromNull(xs.Next())
}

// Reduce reduces stream into one value using given operation.
func Reduce[A, B any](start A, op func(A, B) A, xs Stream[B]) A {
	for x := range xs {
		start = op(start, x)
	}
	return start
}

// Count consumes stream and returns it's length.
func Count[A any](xs Stream[A]) int {
	return Sum(Map(xs, fun.Const[int, A](1)))
}

// Group groups elements by a function that returns a key.
func Group[A any, K comparable](xs Stream[A], by func(A) K) map[K][]A {
	res := make(map[K][]A)
	ForEach(
		xs,
		func(a A) {
			key := by(a)
			vals, ok := res[key]
			if ok {
				vals = append(vals, a)
				res[key] = vals
			} else {
				res[key] = []A{a}
			}
		},
	)
	return res
}

// GroupAggregate is a convenience function that groups and then maps the subslices.
func GroupAggregate[A, B any, K comparable](xs Stream[A], by func(A) K, aggregate func([]A) B) map[K]B {
	tmp := Group(xs, by)
	res := make(map[K]B, len(tmp))
	for k, v := range tmp {
		res[k] = aggregate(v)
	}
	return res
}

// ToCounterBy consumes the stream and returns Counter with count of how many times each key was seen.
func ToCounterBy[A any, K comparable](xs Stream[A], by func(A) K) fun.Counter[K] {
	return GroupAggregate(xs, by, func(ys []A) int { return len(ys) })
}

// CollectCounter consumes the stream makes Counter with count of how many times each element was seen.
func CollectCounter[A comparable](xs Stream[A]) fun.Counter[A] {
	return ToCounterBy(xs, fun.Identity[A])
}

// Any consumes the stream and checks if any of the stream elements matches the predicate.
func Any[A any](xs Stream[A], p func(A) bool) bool {
	return Reduce(false, func(acc bool, a A) bool { return acc || p(a) }, xs)
}
