package iter

// functions to make Iter from something that is not Iter

import (
	"cmp"
	"iter"
)

func FromInt(n int) Seq[int] {
	return func(yield func(int) bool) {
		for i := 0; i < int(n); i++ {
			if !yield(i) {
				return
			}
		}
	}
}

func FromPullFunc[T any](sf func() (T, error)) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for {
			x, err := sf()
			if err != nil {
				yield(x, err)
				return
			}

			if !yield(x, nil) {
				return
			}
		}
	}
}

func FromString(s string) iter.Seq2[int, rune] {
	return func(yield func(int, rune) bool) {
		for i, r := range s {
			if !yield(i, r) {
				return
			}
		}
	}
}

// FromMany returns a stream with all the given values.
func FromMany[V any](vs ...V) Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range vs {
			if !yield(v) {
				return
			}
		}
	}
}

// FromGenerator constructs an infinite stream of values using the production function.
func FromGenerator[A any](x0 A, f func(A) A) Seq[A] {
	return func(yield func(A) bool) {
		for cur := x0; yield(cur); cur = f(cur) {
		}
	}
}

// FromNothing returns an empty stream.
func FromNothing[A any]() Seq[A] {
	return func(yield func(A) bool) {
	}
}

// FromRange makes stream starting with start, step equal to step and going up to end, but not including end.
func FromRange[N cmp.Ordered](start, end, step N) Seq[N] {
	return func(yield func(N) bool) {
		for i := start; i < end && yield(i); i += step {
		}
	}
}

func FromInfiniteGenerator[T any](f func(func(T))) Seq[T] {
	return func(yield func(T) bool) {
		f(func(t T) {
			yield(t)
		})

		return
	}
}
