package stream

// functions to make something from Seq that is not Seq.

import (
	"github.com/rprtr258/go-flow/v2/fun"
)

// ForEach invokes a simple function for each element of the seq.
func ForEach[V any](seq Seq[V], f func(V)) {
	seq(func(v V) bool {
		f(v)
		return true
	})
}

// ToSlice executes the seq and collects all results to a slice.
func ToSlice[A any](seq Seq[A]) []A {
	slice := make([]A, 0, Count(seq))
	ForEach(seq, func(a A) { slice = append(slice, a) })
	return slice
}

// ToSet executes the seq and collects all results to a set.
func ToSet[A comparable](seq Seq[A]) fun.Set[A] {
	set := make(fun.Set[A])
	ForEach(seq, func(a A) { set[a] = fun.Unit1 })
	return set
}

// Head takes the first element if present.
func Head[V any](seq Seq[V]) (V, bool) {
	var (
		res V
		ok  bool
	)
	seq(func(v V) bool {
		res, ok = v, true
		return false
	})
	return res, ok
}

// Reduce reduces seq into one value using given operation.
func Reduce[A, B any](start A, op func(A, B) A, seq Seq[B]) A {
	acc := start
	seq(func(b B) bool {
		acc = op(acc, b)
		return true
	})
	return acc
}

// Count returns seq length.
func Count[A any](seq Seq[A]) int {
	return Sum(Map(seq, fun.Const[int, A](1)))
}

// Count2 returns seq length.
func Count2[A, B any](seq Seq2[A, B]) int {
	return Count(Keys(seq))
}

// Group groups elements by a function that returns a key.
func Group[V any, K comparable](seq Seq[V], by func(V) K) map[K][]V {
	res := make(map[K][]V)
	seq(func(v V) bool {
		key := by(v)
		res[key] = append(res[key], v)
		return true
	})
	return res
}

// GroupAggregate is a convenience function that groups and then maps the subslices.
func GroupAggregate[A, B any, K comparable](seq Seq[A], by func(A) K, aggregate func([]A) B) map[K]B {
	tmp := Group(seq, by)
	res := make(map[K]B, len(tmp))
	for k, v := range tmp {
		res[k] = aggregate(v)
	}
	return res
}

// ToCounterBy consumes the seq and returns Counter with count of how many times each key was seen.
func ToCounterBy[A any, K comparable](seq Seq[A], by func(A) K) fun.Counter[K] {
	return GroupAggregate(seq, by, func(ys []A) int { return len(ys) })
}

// CollectCounter consumes the seq makes Counter with count of how many times each element was seen.
func CollectCounter[A comparable](seq Seq[A]) fun.Counter[A] {
	return ToCounterBy(seq, fun.Identity[A])
}

// Any consumes the seq and checks if any of the seq elements matches the predicate.
func Any[A any](seq Seq[A], p func(A) bool) bool {
	return Reduce(false, func(acc bool, a A) bool { return acc || p(a) }, seq)
}
