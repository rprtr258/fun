// Package stream provides a way to construct data processing streams from smaller pieces.
package iter

import (
	"cmp"
	"fmt"
	"iter"
	"slices"

	"github.com/rprtr258/fun"
	"github.com/rprtr258/fun/set"
)

type Seq[V any] iter.Seq[V]

func (seq Seq[V]) Filter(p func(V) bool) Seq[V] {
	return Filter(seq, p)
}

func (seq Seq[V]) Map(f func(V) V) Seq[V] {
	return Map(seq, f)
}

func (seq Seq[V]) MapFilter(f func(V) (V, bool)) Seq[V] {
	return MapFilter(seq, f)
}

func (seq Seq[V]) FlatMap(f func(V) Seq[V]) Seq[V] {
	return FlatMap(seq, f)
}

func (seq Seq[V]) Take(n int) Seq[V] {
	return Take(seq, n)
}

func (seq Seq[V]) Head() (V, bool) {
	return Head(seq)
}

func (seq Seq[V]) ForEach(f func(V)) {
	ForEach(seq, f)
}

func (seq Seq[V]) Any(p func(V) bool) bool {
	return Any(seq, p)
}

func (seq Seq[V]) All(p func(V) bool) bool {
	return All(seq, p)
}

func (seq Seq[V]) Slice() []V {
	return slices.Collect(iter.Seq[V](seq))
}

func (seq Seq[V]) Count() int {
	return Count(seq)
}

func (seq Seq[V]) Chain(other Seq[V]) Seq[V] {
	return Concat(seq, other)
}

// Map converts values of the stream.
func Map[I, O any](seq Seq[I], f func(I) O) Seq[O] {
	return func(yield func(O) bool) {
		seq(func(a I) bool {
			return yield(f(a))
		})
	}
}

func MapTo2[T, K, V any](seq Seq[T], f func(T) (K, V)) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		seq(func(a T) bool {
			return yield(f(a))
		})
	}
}

func MapFrom2[T, K, V any](seq iter.Seq2[K, V], f func(K, V) T) Seq[T] {
	return func(yield func(T) bool) {
		for k, v := range seq {
			if !yield(f(k, v)) {
				break
			}
		}
	}
}

func Map2[A, B, K, V any](seq iter.Seq2[A, B], f func(A, B) (K, V)) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		seq(func(a A, b B) bool {
			return yield(f(a, b))
		})
	}
}

// Concat returns an iterator over the concatenation of the sequences.
func Concat[V any](seqs ...Seq[V]) Seq[V] {
	return func(yield func(V) bool) {
		for _, seq := range seqs {
			cont := true
			seq(func(v V) bool {
				cont = cont && yield(v)
				return cont
			})
			if !cont {
				return
			}
		}
	}
}

// MergeFunc merges two sequences of values ordered by the function f.
// Values appear in the output once for each time they appear in x
// and once for each time they appear in y.
// When equal values appear in both sequences,
// the output contains the values from x before the values from y.
// If the two input sequences are not ordered by f,
// the output sequence will not be ordered by f,
// but it will still contain every value from x and y exactly once.
func MergeFunc[V any](x, y Seq[V], f func(V, V) int) Seq[V] {
	return func(yield func(V) bool) {
		next, stop := Pull(y)
		defer stop()
		vy, ok := next()
		x(func(vx V) bool {
			for ok && f(vx, vy) > 0 {
				if !yield(vy) {
					return false
				}
				vy, ok = next()
			}
			if !yield(vx) {
				return false
			}
			return true
		})

		for ; ok; vy, ok = next() {
			if !yield(vy) {
				return
			}
		}
	}
}

// Merge merges two sequences of ordered values.
// Values appear in the output once for each time they appear in x
// and once for each time they appear in y.
// If the two input sequences are not ordered,
// the output sequence will not be ordered,
// but it will still contain every value from x and y exactly once.
//
// Merge is equivalent to calling MergeFunc with cmp.Compare[V]
// as the ordering function.
func Merge[V cmp.Ordered](x, y Seq[V]) Seq[V] {
	return MergeFunc(x, y, cmp.Compare[V])
}

// FlatMap maps stream using function and concatenates result streams into one.
func FlatMap[I, O any](seq Seq[I], f func(I) Seq[O]) Seq[O] {
	return func(yield func(O) bool) {
		seq(func(in I) bool {
			cont := true
			f(in)(func(out O) bool {
				cont = cont && yield(out)
				return cont
			})
			return cont
		})
	}
}

// Flatten simplifies a stream of streams to just the stream of values by concatenating all inner streams.
func Flatten[V any](seqseq Seq[Seq[V]]) Seq[V] {
	return func(yield func(V) bool) {
		seqseq(func(seq Seq[V]) bool {
			cont := true
			seq(func(v V) bool {
				cont = cont && yield(v)
				return cont
			})
			return cont
		})
	}
}

// Chunked groups elements by n and produces a stream of slices.
// Produced chunks must not be retained.
func Chunked[A any](xs Seq[A], n int) Seq[[]A] {
	if n <= 0 {
		panic(fmt.Sprintf("Chunk must be of positive size, but %d given", n))
	}

	return func(yield func([]A) bool) {
		chunk := make([]A, 0, n)
		xs(func(a A) bool {
			chunk = append(chunk, a)
			if len(chunk) == n {
				if !yield(chunk) {
					return false
				}

				chunk = chunk[:0]
			}
			return true
		})

		if len(chunk) != 0 {
			yield(chunk)
		}
	}
}

// Intersperse adds a separator after each stream element.
func Intersperse[A any](xs Seq[A], sep A) Seq[A] {
	return func(yield func(A) bool) {
		isFirst := true
		xs(func(a A) bool {
			if !isFirst && !yield(sep) {
				return false
			}

			isFirst = false

			if !yield(a) {
				return false
			}

			return true
		})
	}
}

func Keys[K, V any](xs iter.Seq2[K, V]) Seq[K] {
	return MapFrom2(xs, func(k K, _ V) K {
		return k
	})
}

func Values[K, V any](xs iter.Seq2[K, V]) Seq[V] {
	return MapFrom2(xs, func(_ K, v V) V {
		return v
	})
}

// Repeat appends the same stream infinitely.
func Repeat[A any](xs Seq[A]) Seq[A] {
	return func(yield func(A) bool) {
		for {
			cont := true
			xs(func(x A) bool {
				cont = cont && yield(x)
				return cont
			})
			if !cont {
				return
			}
		}
	}
}

// Take cuts the stream after n elements.
func Take[V any](xs Seq[V], n int) Seq[V] {
	if n < 0 {
		panic(fmt.Sprintf("Take size must be non-negative, but %d given", n))
	}

	return func(yield func(V) bool) {
		took := 0
		xs(func(v V) bool {
			if took == n {
				return false
			}

			took++
			if !yield(v) {
				return false
			}

			return true
		})
	}
}

// Skip skips n elements in the stream.
func Skip[A any](xs Seq[A], n int) Seq[A] {
	return func(yield func(A) bool) {
		skipped := 0
		xs(func(a A) bool {
			if skipped == n {
				if !yield(a) {
					return false
				}
			} else {
				skipped++
			}
			return true
		})
	}
}

// Filter leaves in the stream only the elements that satisfy the given predicate.
func Filter[V any](seq Seq[V], p func(V) bool) Seq[V] {
	return func(yield func(V) bool) {
		seq(func(a V) bool {
			return !p(a) || yield(a)
		})
	}
}

// TakeWhile takes elements while predicate is true.
func TakeWhile[A any](xs Seq[A], p func(A) bool) Seq[A] {
	return func(yield func(A) bool) {
		xs(func(a A) bool {
			return p(a) && yield(a)
		})
	}
}

// DebugSeq prints every processed element, without changing it.
func DebugSeq[A any](xs Seq[A]) Seq[A] {
	return Map(xs, fun.Debug[A])
}

// DebugSeqP prints every processed element, without changing it.
func DebugSeqP[A any](prefix string, xs Seq[A]) Seq[A] {
	return Map(xs, fun.DebugP[A](prefix))
}

// Unique makes stream of unique elements.
func Unique[A comparable](xs Seq[A]) Seq[A] {
	seen := set.New[A](0)
	return func(yield func(A) bool) {
		xs(func(x A) bool {
			if !seen.Contains(x) {
				if !yield(x) {
					return false
				}
				seen.Add(x)
			}
			return true
		})
	}
}

// MapFilter applies function to every element and leaves only elements that are not None.
func MapFilter[I, O any](seq Seq[I], f func(I) (O, bool)) Seq[O] {
	return func(yield func(O) bool) {
		seq(func(a I) bool {
			b, ok := f(a)
			return !ok || yield(b)
		})
	}
}

// Paged makes stream from stream of pages of elements represented as slices.
func Paged[V any](seq Seq[[]V]) Seq[V] {
	return FlatMap(seq, func(vs []V) Seq[V] {
		return FromMany(vs...)
	})
}
