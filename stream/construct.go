package stream

// functions to make Stream from something that is not stream

import (
	"golang.org/x/exp/constraints"
)

// Once returns a stream of one value.
func Once[A any](a A) Stream[A] {
	return FromMany(a)
}

// FromMany returns a stream with all the given values.
func FromMany[A any](as ...A) Stream[A] {
	// cannot inline as compiler cant infer s type
	return Values(Stream2[int, A](Slice[A](as)))
}

// Generate constructs an infinite stream of values using the production function.
func Generate[A any](x0 A, f func(A) A) Stream[A] {
	return PushFunc[A](func(yield func(A) bool) bool {
		for cur := x0; ; cur = f(cur) {
			if !yield(cur) {
				return false
			}
		}
	})
}

// NewStreamEmpty returns an empty stream.
func NewStreamEmpty[A any]() Stream[A] {
	return PushFunc[A](func(yield func(A) bool) bool {
		return true
	})
}

// Range makes stream starting with start, step equal to step and going up to end, but not including end.
func Range[N constraints.Ordered](start, end, step N) Stream[N] {
	return PushFunc[N](func(yield func(N) bool) bool {
		for i := start; i < end; i += step {
			if !yield(i) {
				return false
			}
		}

		return true
	})
}

func NewGenerator[T any](f func(func(T))) Stream[T] {
	return PushFunc[T](func(yield func(T) bool) bool {
		f(func(t T) {
			yield(t)
		})

		return true
	})
}
