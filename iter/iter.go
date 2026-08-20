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

type (
	Seq[V any]     iter.Seq[V]
	Seq2[K, V any] iter.Seq2[K, V]
)

func (seq Seq[V]) Slice() []V {
	return slices.Collect(iter.Seq[V](seq))
}

func (seq Seq[V]) Chain(other Seq[V]) Seq[V] {
	return Concat(seq, other)
}

// Map converts values of the stream.
func (seq Seq[I]) Map[O any](f func(I) O) Seq[O] {
	return func(yield func(O) bool) {
		seq(func(a I) bool {
			return yield(f(a))
		})
	}
}

func (seq Seq[T]) MapTo2[K, V any](f func(T) (K, V)) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		seq(func(a T) bool {
			return yield(f(a))
		})
	}
}

func (seq Seq2[K, V]) MapFrom2[T any](f func(K, V) T) Seq[T] {
	return func(yield func(T) bool) {
		for k, v := range seq {
			if !yield(f(k, v)) {
				break
			}
		}
	}
}

func (seq Seq2[A, B]) Map2[K, V any](f func(A, B) (K, V)) Seq2[K, V] {
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
func (x Seq[V]) MergeFunc(y Seq[V], f func(V, V) int) Seq[V] {
	return func(yield func(V) bool) {
		next, stop := y.Pull()
		defer stop()
		vy, ok := next()
		x(func(vx V) bool {
			for ok && f(vx, vy) > 0 {
				if !yield(vy) {
					return false
				}
				vy, ok = next()
			}
			return yield(vx)
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
	return x.MergeFunc(y, cmp.Compare[V])
}

// FlatMap maps stream using function and concatenates result streams into one.
func (seq Seq[I]) FlatMap[O any](f func(I) Seq[O]) Seq[O] {
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
func Chunked[V any](xs Seq[V], n int) Seq[[]V] {
	if n <= 0 {
		panic(fmt.Sprintf("Chunk must be of positive size, but %d given", n))
	}

	return func(yield func([]V) bool) {
		chunk := make([]V, 0, n)
		xs(func(v V) bool {
			chunk = append(chunk, v)
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
func (xs Seq[V]) Intersperse(sep V) Seq[V] {
	return func(yield func(V) bool) {
		isFirst := true
		xs(func(v V) bool {
			if !isFirst && !yield(sep) {
				return false
			}

			isFirst = false

			return yield(v)
		})
	}
}

func (xs Seq2[K, V]) Keys() Seq[K] {
	return xs.MapFrom2(func(k K, _ V) K {
		return k
	})
}

func (xs Seq2[K, V]) Values() Seq[V] {
	return xs.MapFrom2(func(_ K, v V) V {
		return v
	})
}

// Repeat appends the same stream infinitely.
func (xs Seq[V]) Repeat() Seq[V] {
	return func(yield func(V) bool) {
		for {
			cont := true
			xs(func(x V) bool {
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
func (xs Seq[V]) Take(n int) Seq[V] {
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
			return yield(v)
		})
	}
}

// Skip skips n elements in the stream.
func (xs Seq[V]) Skip(n int) Seq[V] {
	return func(yield func(V) bool) {
		skipped := 0
		xs(func(v V) bool {
			if skipped == n {
				if !yield(v) {
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
func (seq Seq[V]) Filter(p func(V) bool) Seq[V] {
	return func(yield func(V) bool) {
		seq(func(v V) bool {
			return !p(v) || yield(v)
		})
	}
}

// TakeWhile takes elements while predicate is true.
func (xs Seq[V]) TakeWhile(p func(V) bool) Seq[V] {
	return func(yield func(V) bool) {
		xs(func(v V) bool {
			return p(v) && yield(v)
		})
	}
}

// DebugSeq prints every processed element, without changing it.
func (xs Seq[V]) DebugSeq() Seq[V] {
	return xs.Map(fun.Debug[V])
}

// DebugSeqP prints every processed element, without changing it.
func (xs Seq[V]) DebugSeqP(prefix string) Seq[V] {
	return xs.Map(fun.DebugP[V](prefix))
}

// Unique makes stream of unique elements.
func Unique[V comparable](xs Seq[V]) Seq[V] {
	seen := set.New[V](0)
	return func(yield func(V) bool) {
		xs(func(x V) bool {
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
func (seq Seq[I]) MapFilter[O any](f func(I) (O, bool)) Seq[O] {
	return func(yield func(O) bool) {
		seq(func(a I) bool {
			b, ok := f(a)
			return !ok || yield(b)
		})
	}
}

// Paged makes stream from stream of pages of elements represented as slices.
func Paged[V any](seq Seq[[]V]) Seq[V] {
	return seq.FlatMap(func(vs []V) Seq[V] {
		return FromMany(vs...)
	})
}
