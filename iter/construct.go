package iter

// functions to make Stream from something that is not stream

import (
	"golang.org/x/exp/constraints"

	"github.com/rprtr258/fun"
)

func FromInt(n int) Seq[int] {
	return func(yield func(int) bool) bool {
		for i := 0; i < int(n); i++ {
			if !yield(i) {
				return false
			}
		}

		return true
	}
}

func FromPullFunc[T any](sf func() (T, error)) Seq[fun.Result[T]] {
	return func(yield func(r fun.Result[T]) bool) bool {
		for {
			x, err := sf()
			if err != nil {
				yield(fun.Result[T]{x, err, false})
				return true
			}

			if !yield(fun.Result[T]{x, nil, true}) {
				return false
			}
		}
	}
}

func FromSlice[T any](s []T) Seq[fun.Pair[int, T]] {
	return func(yield func(fun.Pair[int, T]) bool) bool {
		for i, x := range s {
			if !yield(fun.Pair[int, T]{i, x}) {
				return false
			}
		}

		return true
	}
}

func FromString(s string) Seq[fun.Pair[int, rune]] {
	return func(yield func(fun.Pair[int, rune]) bool) bool {
		for i, r := range s {
			if !yield(fun.Pair[int, rune]{i, r}) {
				return false
			}
		}

		return true
	}
}

func FromDict[K comparable, V any](d map[K]V) Seq[fun.Pair[K, V]] {
	return func(yield func(fun.Pair[K, V]) bool) bool {
		for k, v := range d {
			if !yield(fun.Pair[K, V]{k, v}) {
				return false
			}
		}

		return true
	}
}

// FromMany returns a stream with all the given values.
func FromMany[V any](vs ...V) Seq[V] {
	return func(yield func(V) bool) bool {
		for _, v := range vs {
			if !yield(v) {
				return false
			}
		}
		return true
	}
}

// FromGenerator constructs an infinite stream of values using the production function.
func FromGenerator[A any](x0 A, f func(A) A) Seq[A] {
	return func(yield func(A) bool) bool {
		for cur := x0; ; cur = f(cur) {
			if !yield(cur) {
				return false
			}
		}
	}
}

// FromNothing returns an empty stream.
func FromNothing[A any]() Seq[A] {
	return func(yield func(A) bool) bool {
		return true
	}
}

// FromRange makes stream starting with start, step equal to step and going up to end, but not including end.
func FromRange[N constraints.Ordered](start, end, step N) Seq[N] {
	return func(yield func(N) bool) bool {
		for i := start; i < end; i += step {
			if !yield(i) {
				return false
			}
		}

		return true
	}
}

func FromInfiniteGenerator[T any](f func(func(T))) Seq[T] {
	return func(yield func(T) bool) bool {
		f(func(t T) {
			yield(t)
		})

		return true
	}
}
