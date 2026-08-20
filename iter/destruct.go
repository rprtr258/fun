package iter

// functions to make something from Seq that is not Seq.

import (
	"github.com/rprtr258/fun"
	"github.com/rprtr258/fun/set"
)

// ForEach invokes a simple function for each element of the seq.
func (seq Seq[V]) ForEach(f func(V)) {
	seq(func(v V) bool {
		f(v)
		return true
	})
}

// ToSet executes the seq and collects all results to a set.
func ToSet[V comparable](seq Seq[V]) set.Set[V] {
	set := set.New[V](0)
	seq.ForEach(func(a V) { set.Add(a) })
	return set
}

// Head takes the first element if present.
func (seq Seq[V]) Head() (V, bool) {
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
func (seq Seq[B]) Reduce[A any](start A, op func(A, B) A) A {
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
	return xs.Reduce(zero, func(x A, y A) A {
		return x + y
	})
}

// Count returns seq length.
func (seq Seq[V]) Count() int {
	res := 0
	seq(func(V) bool {
		res++
		return true
	})
	return res
}

// Group groups elements by a function that returns a key.
func (seq Seq[V]) Group[K comparable](by func(V) K) map[K][]V {
	res := make(map[K][]V)
	seq(func(v V) bool {
		key := by(v)
		res[key] = append(res[key], v)
		return true
	})
	return res
}

// GroupAggregate is a convenience function that groups and then maps the subslices.
func (seq Seq[A]) GroupAggregate[B any, K comparable](by func(A) K, aggregate func([]A) B) map[K]B {
	tmp := seq.Group(by)
	res := make(map[K]B, len(tmp))
	for k, v := range tmp {
		res[k] = aggregate(v)
	}
	return res
}

// ToCounterBy consumes the seq and returns Counter with count of how many times each key was seen.
func (seq Seq[V]) ToCounterBy[K comparable](by func(V) K) map[K]int {
	return seq.GroupAggregate(by, func(ys []V) int { return len(ys) })
}

// ToCounter consumes the seq makes Counter with count of how many times each element was seen.
func ToCounter[V comparable](seq Seq[V]) map[V]int {
	return seq.ToCounterBy(func(v V) V { return v })
}

// Any consumes the seq and checks if any of the seq elements matches the predicate
func (seq Seq[V]) Any(p func(V) bool) bool {
	found := false
	seq(func(v V) bool {
		found = p(v)
		return !found
	})
	return found
}

// All consumes the seq and checks if all of the seq elements match the predicate
func (seq Seq[V]) All(p func(V) bool) bool {
	res := true
	seq(func(v V) bool {
		if !p(v) {
			res = false
		}
		return res
	})
	return res
}

func (push Seq[V]) Pull() (pull func() (V, bool), stop func()) {
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

// Find searches for first element matching the predicate.
func (xs Seq[V]) Find(p func(V) bool) (V, bool) {
	var aa V
	found := false
	xs(func(a V) bool {
		if p(a) {
			found = true
			aa = a
			return false
		}
		return true
	})

	return aa, found
}
