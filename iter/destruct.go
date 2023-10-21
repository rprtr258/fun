package iter

// functions to make something from Seq that is not Seq.

import (
	"github.com/rprtr258/fun"
	"github.com/rprtr258/fun/set"
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
func ToSet[A comparable](seq Seq[A]) set.Set[A] {
	set := set.New[A](0)
	ForEach(seq, func(a A) { set.Add(a) })
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
func Reduce[A, B any](seq Seq[B], start A, op func(A, B) A) A {
	acc := start
	seq(func(b B) bool {
		acc = op(acc, b)
		return true
	})
	return acc
}

// Sum finds sum of elements in stream.
func Sum[A fun.Number](xs Seq[A]) A {
	var zero A
	return Reduce(xs, zero, func(x A, y A) A {
		return x + y
	})
}

// Count returns seq length.
func Count[V any](seq Seq[V]) int {
	res := 0
	seq(func(V) bool {
		res++
		return true
	})
	return res
}

// Group groups elements by a function that returns a key.
func Group[K comparable, V any](seq Seq[V], by func(V) K) map[K][]V {
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
func ToCounterBy[K comparable, V any](seq Seq[V], by func(V) K) map[K]int {
	return GroupAggregate(seq, by, func(ys []V) int { return len(ys) })
}

// ToCounter consumes the seq makes Counter with count of how many times each element was seen.
func ToCounter[V comparable](seq Seq[V]) map[V]int {
	return ToCounterBy(seq, func(v V) V { return v })
}

// Any consumes the seq and checks if any of the seq elements matches the predicate
func Any[V any](seq Seq[V], p func(V) bool) bool {
	found := false
	seq(func(v V) bool {
		found = p(v)
		return !found
	})
	return found
}

// All consumes the seq and checks if all of the seq elements match the predicate
func All[V any](seq Seq[V], p func(V) bool) bool {
	res := true
	seq(func(v V) bool {
		if !p(v) {
			res = false
		}
		return res
	})
	return res
}

func Pull[V any](push Seq[V]) (pull func() (V, bool), stop func()) {
	copush := func(more bool, yield func(V) bool) V {
		if more {
			push(yield)
		}
		var zero V
		return zero
	}

	cin := make(chan bool)
	cout := make(chan V)
	running := true
	resume := func(in bool) (out V, ok bool) {
		if !running {
			return
		}
		cin <- in
		out = <-cout
		return out, running
	}
	yield := func(out V) bool {
		cout <- out
		return <-cin
	}
	go func() {
		out := copush(<-cin, yield)
		running = false
		cout <- out
	}()
	pull = func() (V, bool) {
		return resume(true)
	}
	stop = func() {
		resume(false)
	}
	return pull, stop
}
