// Package stream provides a way to construct data processing streams from smaller pieces.
package iter

import (
	"fmt"

	"github.com/rprtr258/fun"
)

type (
	Seq[V any] func(yield func(V) bool) bool
)

// Map converts values of the stream.
func Map[I, O any](seq Seq[I], f func(I) O) Seq[O] {
	return func(yield func(O) bool) bool {
		// for a := range seq {
		return seq(func(a I) bool {
			return yield(f(a))
		})
	}
}

// Concat returns an iterator over the concatenation of the sequences.
func Concat[V any](seqs ...Seq[V]) Seq[V] {
	return func(yield func(V) bool) bool {
		for _, seq := range seqs {
			if !seq(yield) {
				return false
			}
		}

		return true
	}
}

// FlatMap maps stream using function and concatenates result streams into one.
func FlatMap[I, O any](seq Seq[I], f func(I) Seq[O]) Seq[O] {
	return func(yield func(O) bool) bool {
		return seq(func(in I) bool {
			return f(in)(yield)
		})
	}
}

// Flatten simplifies a stream of streams to just the stream of values by concatenating all inner streams.
func Flatten[A any](xs Seq[Seq[A]]) Seq[A] {
	return FlatMap(xs, fun.Identity[Seq[A]])
}

// Chunked groups elements by n and produces a stream of slices.
// Produced chunks must not be retained.
func Chunked[A any](xs Seq[A], n int) Seq[[]A] {
	if n <= 0 {
		panic(fmt.Sprintf("Chunk must be of positive size, but %d given", n))
	}

	return func(yield func([]A) bool) bool {
		chunk := make([]A, 0, n)
		if !xs(func(a A) bool {
			chunk = append(chunk, a)
			if len(chunk) == n {
				if !yield(append([]A(nil), chunk...)) {
					return false
				}

				chunk = chunk[:0]
			}

			return true
		}) {
			return false
		}

		return len(chunk) == 0 || yield(chunk)
	}
}

// Intersperse adds a separator after each stream element.
func Intersperse[A any](xs Seq[A], sep A) Seq[A] {
	return func(yield func(A) bool) bool {
		isFirst := true
		return xs(func(a A) bool {
			if !isFirst && !yield(sep) {
				return false
			}

			isFirst = false

			return yield(a)
		})
	}
}

func Keys[K, V any](xs Seq[fun.Pair[K, V]]) Seq[K] {
	return Map(xs, func(p fun.Pair[K, V]) K {
		return p.K
	})
}

func Values[K, V any](xs Seq[fun.Pair[K, V]]) Seq[V] {
	return Map(xs, func(p fun.Pair[K, V]) V {
		return p.V
	})
}

// Repeat appends the same stream infinitely.
func Repeat[A any](xs Seq[A]) Seq[A] {
	return func(yield func(A) bool) bool {
		for xs(yield) {
		}
		return false
	}
}

// Take cuts the stream after n elements.
func Take[V any](xs Seq[V], n int) Seq[V] {
	if n < 0 {
		panic(fmt.Sprintf("Take size must be non-negative, but %d given", n))
	}

	return func(yield func(V) bool) bool {
		took := 0
		return xs(func(v V) bool {
			if took == n {
				return false
			}

			took++
			return yield(v)
		})
	}
}

// Skip skips n elements in the stream.
func Skip[A any](xs Seq[A], n int) Seq[A] {
	return func(yield func(A) bool) bool {
		skipped := 0
		return xs(func(a A) bool {
			if skipped == n {
				return yield(a)
			}

			skipped++
			return true
		})
	}
}

// Filter leaves in the stream only the elements that satisfy the given predicate.
func Filter[V any](seq Seq[V], p func(V) bool) Seq[V] {
	return func(yield func(V) bool) bool {
		// for a := range seq {
		return seq(func(a V) bool {
			if p(a) && !yield(a) {
				return false
			}

			return true
		})
	}
}

// Find searches for first element matching the predicate.
func Find[A any](xs Seq[A], p func(A) bool) fun.Option[A] {
	var res fun.Option[A]
	xs(func(a A) bool {
		if p(a) {
			res = fun.Some(a)
			return false
		}

		return true
	})
	return res
}

// TakeWhile takes elements while predicate is true.
func TakeWhile[A any](xs Seq[A], p func(A) bool) Seq[A] {
	return func(yield func(A) bool) bool {
		return xs(func(a A) bool {
			if !p(a) {
				return false
			}

			yield(a)
			return true
		})
	}
}

// DebugPrint prints every processed element, without changing it.
func DebugPrint[A any](prefix string, xs Seq[A]) Seq[A] {
	return Map(xs, fun.Debug[A](prefix))
}

// Unique makes stream of unique elements.
func Unique[A comparable](xs Seq[A]) Seq[A] {
	seen := fun.NewSet[A]()
	var a A
	return MapFilter(xs, func(x A) (A, bool) {
		if seen.Contains(x) {
			return a, false
		}

		seen.Add(x)
		return x, true
	})
}

// MapFilter applies function to every element and leaves only elements that are not None.
func MapFilter[I, O any](seq Seq[I], f func(I) (O, bool)) Seq[O] {
	return func(yield func(O) bool) bool {
		return seq(func(a I) bool {
			if b, ok := f(a); ok {
				return yield(b)
			}

			return true
		})
	}
}

// Paged makes stream from stream of pages of elements represented as slices.
func Paged[V any](seq Seq[[]V]) Seq[V] {
	return FlatMap(seq, func(vs []V) Seq[V] {
		return FromMany(vs...)
	})
}
