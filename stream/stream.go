// Package stream provides a way to construct data processing streams from smaller pieces.
package stream

import (
	"errors"
	"fmt"

	"github.com/rprtr258/go-flow/v2/fun"
)

var ErrBreak = errors.New("break prematurely")

// Iterator2 is iterator taking function which handles iteratee and returns
// whether to stop iteration. Returns error happened during iteration or ErrBreak
// if exited prematurely.
type Stream[T any] interface {
	For(func(T) bool) error
}

type Chan[T any] <-chan T

func (c Chan[T]) For(f func(T) bool) error {
	for x := range c {
		if !f(x) {
			return ErrBreak
		}
	}

	return nil
}

type Int int

func (n Int) For(f func(int) bool) error {
	for i := 0; i < int(n); i++ {
		if !f(i) {
			return ErrBreak
		}
	}

	return nil
}

// Must return true if yield has not returned false
type PushFunc[T any] func(func(T) bool) bool

func (sf PushFunc[T]) For(f func(T) bool) error {
	if !sf(f) {
		return ErrBreak
	}

	return nil
}

var ErrPullEnd = errors.New("pull end prematurely")

type PullFunc[T any] func() (T, error)

func (sf PullFunc[T]) For(f func(T) bool) error {
	for {
		x, err := sf()
		if err != nil {
			if err == ErrPullEnd {
				return nil
			}

			return err
		}

		if !f(x) {
			return ErrBreak
		}
	}
}

type Stream2[K, V any] interface {
	For(func(K, V) bool) error
}

type Slice[T any] []T

func (s Slice[T]) For(f func(int, T) bool) error {
	for i, x := range s {
		if !f(i, x) {
			return ErrBreak
		}
	}

	return nil
}

type String string

func (s String) For(f func(int, rune) bool) error {
	for i, r := range s {
		if !f(i, r) {
			return ErrBreak
		}
	}

	return nil
}

type Dict[K comparable, V any] map[K]V

func (d Dict[K, V]) For(f func(K, V) bool) error {
	for k, v := range d {
		if !f(k, v) {
			return ErrBreak
		}
	}

	return nil
}

// Must return true if yield has not returned false
type PushFunc2[K, V any] func(func(K, V) bool) bool

func (sf PushFunc2[K, V]) For(f func(K, V) bool) error {
	if !sf(f) {
		return ErrBreak
	}

	return nil
}

// Map converts values of the stream.
func Map[A, B any](xs Stream[A], f func(A) B) Stream[B] {
	return PushFunc[B](func(yield func(B) bool) bool {
		return xs.For(func(a A) bool {
			return yield(f(a))
		}) != ErrBreak
	})
}

// Chain appends another stream after the end of the first one.
func Chain[A any](xss ...Stream[A]) Stream[A] {
	return PushFunc[A](func(yield func(A) bool) bool {
		for _, xs := range xss {
			if err := xs.For(func(a A) bool {
				return yield(a)
			}); err != nil {
				return false
			}
		}

		return true
	})
}

// FlatMap maps stream using function and concatenates result streams into one.
func FlatMap[A, B any](xs Stream[A], f func(A) Stream[B]) Stream[B] {
	return PushFunc[B](func(yield func(B) bool) bool {
		return xs.For(func(a A) bool {
			return f(a).For(yield) == nil
		}) == nil
	})
}

// Flatten simplifies a stream of streams to just the stream of values by concatenating all inner streams.
func Flatten[A any](xs Stream[Stream[A]]) Stream[A] {
	return FlatMap(xs, fun.Identity[Stream[A]])
}

// Sum finds sum of elements in stream.
func Sum[A fun.Number](xs Stream[A]) A {
	var zero A
	return Reduce(zero,
		func(x A, y A) A {
			return x + y
		},
		xs,
	)
}

// Chunked groups elements by n and produces a stream of slices.
// Produced chunks must not be retained.
func Chunked[A any](xs Stream[A], n int) Stream[[]A] {
	if n <= 0 {
		panic(fmt.Sprintf("Chunk must be of positive size, but %d given", n))
	}

	return PushFunc[[]A](func(f func([]A) bool) bool {
		chunk := make([]A, 0, n)
		if err := xs.For(func(a A) bool {
			chunk = append(chunk, a)
			if len(chunk) == n {
				if !f(chunk) {
					return false
				}

				chunk = chunk[:0]
			}

			return true
		}); err != nil {
			return false
		}

		return len(chunk) == 0 || f(chunk)
	})
}

// Intersperse adds a separator after each stream element.
func Intersperse[A any](xs Stream[A], sep A) Stream[A] {
	return PushFunc[A](func(yield func(A) bool) bool {
		isFirst := true
		return xs.For(func(a A) bool {
			if !isFirst && !yield(sep) {
				return false
			}

			isFirst = false

			return yield(a)
		}) == nil
	})
}

func Keys[K, V any](xs Stream2[K, V]) Stream[K] {
	return PushFunc[K](func(yield func(K) bool) bool {
		return xs.For(func(k K, _ V) bool {
			return yield(k)
		}) == nil
	})
}

func Values[K, V any](xs Stream2[K, V]) Stream[V] {
	return PushFunc[V](func(yield func(V) bool) bool {
		return xs.For(func(_ K, v V) bool {
			return yield(v)
		}) == nil
	})
}

// Repeat appends the same stream infinitely.
func Repeat[A any](xs Stream[A]) Stream[A] {
	return PushFunc[A](func(yield func(A) bool) bool {
		for {
			if err := xs.For(func(a A) bool {
				return yield(a)
			}); err != nil {
				return false
			}
		}
	})
}

// Take cuts the stream after n elements.
func Take[A any](xs Stream[A], n int) Stream[A] {
	if n < 0 {
		panic(fmt.Sprintf("Take size must be non-negative, but %d given", n))
	}

	return PushFunc[A](func(yield func(A) bool) bool {
		took := 0
		return xs.For(func(a A) bool {
			if took == n {
				return false
			}

			took++
			return yield(a)
		}) == nil
	})
}

// Skip skips n elements in the stream.
func Skip[A any](xs Stream[A], n int) Stream[A] {
	return PushFunc[A](func(yield func(A) bool) bool {
		skipped := 0
		return xs.For(func(a A) bool {
			if skipped == n {
				return yield(a)
			}

			skipped++
			return true
		}) == nil
	})
}

// Filter leaves in the stream only the elements that satisfy the given predicate.
func Filter[A any](xs Stream[A], p func(A) bool) Stream[A] {
	return PushFunc[A](func(yield func(A) bool) bool {
		return xs.For(func(a A) bool {
			if p(a) {
				return yield(a)
			}

			return true
		}) == nil
	})
}

// Find searches for first element matching the predicate.
func Find[A any](xs Stream[A], p func(A) bool) fun.Option[A] {
	var res fun.Option[A]
	xs.For(func(a A) bool {
		if p(a) {
			res = fun.Some(a)
			return false
		}

		return true
	})
	return res
}

// TakeWhile takes elements while predicate is true.
func TakeWhile[A any](xs Stream[A], p func(A) bool) Stream[A] {
	return PushFunc[A](func(yield func(A) bool) bool {
		return xs.For(func(a A) bool {
			if !p(a) {
				return false
			}

			yield(a)
			return true
		}) == nil
	})
}

// DebugPrint prints every processed element, without changing it.
func DebugPrint[A any](prefix string, xs Stream[A]) Stream[A] {
	return Map(xs, fun.Debug[A](prefix))
}

// Unique makes stream of unique elements.
func Unique[A comparable](xs Stream[A]) Stream[A] {
	seen := fun.NewSet[A]()
	return MapFilter(xs, func(x A) fun.Option[A] {
		if !seen.Contains(x) {
			seen.Add(x)
			return fun.Some(x)
		}

		return fun.None[A]()
	})
}

// MapFilter applies function to every element and leaves only elements that are not None.
func MapFilter[A, B any](xs Stream[A], f func(A) fun.Option[B]) Stream[B] {
	return PushFunc[B](func(yield func(B) bool) bool {
		return xs.For(func(a A) bool {
			if b, ok := f(a).Unpack(); ok {
				return yield(b)
			}

			return true
		}) == nil
	})
}

// Paged makes stream from stream of pages of elements represented as slices.
func Paged[A any](xs Stream[[]A]) Stream[A] {
	return FlatMap(xs, func(as []A) Stream[A] {
		return FromMany(as...)
	})
}

// // ScatterCopy copies stream into N streams with all source elements.
// // ALL resulting streams must be consumed concurrently, use with caution.
// func ScatterCopy[A any](xs Stream[A], n int) []Stream[A] {
// 	chans := make([]chan A, n)
// 	for i := range chans {
// 		chans[i] = make(chan A)
// 	}

// 	go func() {
// 		for x := range xs {
// 			for _, ch := range chans {
// 				ch <- x
// 			}
// 		}
// 		for _, ch := range chans {
// 			close(ch)
// 		}
// 	}()

// 	res := make([]Stream[A], n)
// 	for i := range res {
// 		res[i] = chans[i]
// 	}
// 	return res
// }

// // Scatter splits stream into N streams, each element from source stream goes to one of the result streams.
// func Scatter[A any](xs Stream[A], n int) []Stream[A] {
// 	res := make([]Stream[A], n)
// 	for i := range res {
// 		outCh := make(chan A)
// 		go func() {
// 			for x := range xs {
// 				outCh <- x
// 			}
// 			close(outCh)
// 		}()
// 		res[i] = outCh
// 	}
// 	return res
// }

// // ScatterEvenly splits stream into N streams of source elements using round robin distribution
// func ScatterEvenly[A any](xs Stream[A], n int) []Stream[A] {
// 	chans := make([]chan A, 0, n)
// 	for i := 0; i < n; i++ {
// 		chans = append(chans, make(chan A))
// 	}
// 	go func() {
// 		for {
// 			for _, outCh := range chans {
// 				x, ok := <-xs
// 				if !ok {
// 					goto END
// 				}
// 				outCh <- x
// 			}
// 		}
// 	END:
// 		for _, outCh := range chans {
// 			close(outCh)
// 		}
// 	}()

// 	return CollectToSlice(Map(
// 		FromSlice(chans),
// 		func(c chan A) Stream[A] { return c },
// 	))
// }

// // ScatterRoute routes elements from source stream into first matching predicate stream or last stream if non matched.
// // Routes are functions from source element index and element to bool: does element match the route or not.
// func ScatterRoute[A any](xs Stream[A], routes []func(int, A) bool) []Stream[A] {
// 	n := len(routes)
// 	chans := make([]chan A, 0, n+1)
// 	for i := 0; i < n+1; i++ {
// 		chans = append(chans, make(chan A))
// 	}
// 	go func() {
// 		idx := 0
// 		for {
// 			x, ok := <-xs
// 			if !ok {
// 				goto END
// 			}
// 			for i := 0; i < n; i++ {
// 				if routes[i](idx, x) {
// 					chans[i] <- x
// 					goto MATCH_FOUND
// 				}
// 			}
// 			chans[n] <- x
// 		MATCH_FOUND:
// 			idx++
// 		}
// 	END:
// 		for _, outCh := range chans {
// 			close(outCh)
// 		}
// 	}()

// 	return CollectToSlice(Map(
// 		FromSlice(chans),
// 		func(c chan A) Stream[A] { return c },
// 	))
// }

// // Gather gets elements from all input streams into single stream.
// func Gather[A any](xss []Stream[A]) Stream[A] {
// 	res := make(chan A)
// 	done := make(chan fun.Unit, len(xss))
// 	for _, xs := range xss {
// 		go func(xs Stream[A]) {
// 			for x := range xs {
// 				res <- x
// 			}
// 			done <- fun.Unit1
// 		}(xs)
// 	}
// 	go func() {
// 		for range xss {
// 			<-done
// 		}
// 		close(res)
// 	}()
// 	return res
// }
