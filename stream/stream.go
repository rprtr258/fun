// Package stream provides a way to construct data processing streams from smaller pieces.
package stream

import (
	"fmt"

	"github.com/rprtr258/go-flow/fun"
)

// Stream is a finite or infinite stream of values.
type Stream[A any] <-chan A

// Next gives either value or nothing if stream has ended.
func (stream Stream[A]) Next() (A, bool) {
	a, end := <-stream
	return a, end
}

// Map converts values of the stream.
func Map[A, B any](xs Stream[A], f func(A) B) Stream[B] {
	res := make(chan B)
	go func() {
		for x := range xs {
			res <- f(x)
		}
		close(res)
	}()
	return res
}

// Chain appends another stream after the end of the first one.
func Chain[A any](xss ...Stream[A]) Stream[A] {
	res := make(chan A)
	go func() {
		for _, xs := range xss {
			for x := range xs {
				res <- x
			}
		}
		close(res)
	}()
	return res
}

// FlatMap maps stream using function and concatenates result streams into one.
func FlatMap[A, B any](xs Stream[A], f func(A) Stream[B]) Stream[B] {
	res := make(chan B)
	go func() {
		for x := range xs {
			for y := range f(x) {
				res <- y
			}
		}
		close(res)
	}()
	return res
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
func Chunked[A any](xs Stream[A], n int) Stream[[]A] {
	if n <= 0 {
		panic(fmt.Sprintf("Chunk must be of positive size, but %d given", n))
	}
	res := make(chan []A)
	go func() {
		chunk := make([]A, 0, n)
		for x := range xs {
			chunk = append(chunk, x)
			if len(chunk) == n {
				res <- chunk
				chunk = make([]A, 0, n)
			}
		}
		if len(chunk) > 0 {
			res <- chunk
		}
		close(res)
	}()
	return res
}

// Intersperse adds a separator after each stream element.
func Intersperse[A any](xs Stream[A], sep A) Stream[A] {
	res := make(chan A)
	go func() {
		isFirst := true
		for x := range xs {
			if !isFirst {
				res <- sep
			}
			res <- x
			isFirst = false
		}
		close(res)
	}()
	return res
}

// Repeat appends the same stream infinitely.
func Repeat[A any](xs Stream[A]) Stream[A] {
	res := make(chan A)
	go func() {
		buf := []A{}
		for x := range xs {
			res <- x
			buf = append(buf, x)
		}
		for {
			for _, x := range buf {
				res <- x
			}
		}
	}()
	return res
}

// Take cuts the stream after n elements.
func Take[A any](xs Stream[A], n int) Stream[A] {
	if n < 0 {
		panic(fmt.Sprintf("Cannot take negative number of elements, concretely %d", n))
	}
	res := make(chan A)
	go func() {
		defer close(res)
		if n == 0 {
			return
		}
		for x := range xs {
			res <- x
			n--
			if n == 0 {
				break
			}
		}
		// TODO: exhaust xs?
	}()
	return res
}

// Skip skips n elements in the stream.
func Skip[A any](xs Stream[A], n int) Stream[A] {
	res := make(chan A)
	go func() {
		for i := 0; i < n; i++ {
			<-xs
		}
		for x := range xs {
			res <- x
		}
		close(res)
	}()
	return res
}

// Filter leaves in the stream only the elements that satisfy the given predicate.
func Filter[A any](xs Stream[A], p func(A) bool) Stream[A] {
	res := make(chan A)
	go func() {
		for x := range xs {
			if p(x) {
				res <- x
			}
		}
		close(res)
	}()
	return res
}

// Find searches for first element matching the predicate.
func Find[A any](xs Stream[A], p func(A) bool) fun.Option[A] {
	x, stop := <-Filter(xs, p)
	return fun.FromNull(x, stop)
}

// TakeWhile takes elements while predicate is true.
func TakeWhile[A any](xs Stream[A], p func(A) bool) Stream[A] {
	res := make(chan A)
	go func() {
		for x := range xs {
			if p(x) {
				res <- x
			} else {
				break
			}
		}
		close(res)
	}()
	return res
}

// DebugPrint prints every processed element, without changing it.
func DebugPrint[A any](prefix string, xs Stream[A]) Stream[A] {
	return Map(xs, fun.Debug[A](prefix))
}

// Unique makes stream of unique elements.
func Unique[A comparable](xs Stream[A]) Stream[A] {
	res := make(chan A)
	go func() {
		seen := fun.NewSet[A]()
		for x := range xs {
			if !seen.Contains(x) {
				seen[x] = fun.Unit1 //TODO: put
				res <- x
			}
		}
		close(res)
	}()
	return res
}

// MapFilter applies function to every element and leaves only elements that are not None.
func MapFilter[A, B any](xs Stream[A], f func(A) fun.Option[B]) Stream[B] {
	res := make(chan B)
	go func() {
		for x := range xs {
			y := f(x)
			if y.IsSome() {
				res <- y.Unwrap()
			}
		}
		close(res)
	}()
	return res
}

// Paged makes stream from stream of pages of elements represented as slices.
func Paged[A any](xs Stream[[]A]) Stream[A] {
	return FlatMap(xs, FromSlice[A])
}

// type bufEntry[A any] struct {
// 	value       A
// 	streamsLeft int64
// }

// type scatterCopyImpl[A any] struct {
// 	source       Stream[A]
// 	buf          *[]bufEntry[A]
// 	deleted      *int
// 	mu           *sync.Mutex
// 	totalStreams int64

// 	index int
// }

// func (xs *scatterCopyImpl[A]) advanceBuffer() {
// 	if xs.index-*xs.deleted == len(*xs.buf) {
// 		x := xs.source.Next()
// 		if x.IsNone() {
// 			return
// 		}
// 		bufEntry := bufEntry[A]{
// 			value:       x.Unwrap(),
// 			streamsLeft: xs.totalStreams,
// 		}
// 		if len(*xs.buf) == cap(*xs.buf) {
// 			// reallocate buffer
// 			buf := append(*xs.buf, bufEntry)
// 			*xs.buf = buf
// 		} else {
// 			// insert into buffer without realloc, **xs.buf is not changed
// 			*xs.buf = append(*xs.buf, bufEntry)
// 		}
// 	}
// }

// func (xs *scatterCopyImpl[A]) Next() fun.Option[A] {
// 	xs.mu.Lock()
// 	defer xs.mu.Unlock()

// 	xs.advanceBuffer()

// 	effectiveIndex := xs.index - *xs.deleted
// 	if effectiveIndex == len(*xs.buf) {
// 		return fun.None[A]()
// 	}
// 	atomic.AddInt64(&(*xs.buf)[effectiveIndex].streamsLeft, -1)
// 	x := (*xs.buf)[effectiveIndex].value
// 	if (*xs.buf)[effectiveIndex].streamsLeft == 0 {
// 		*xs.deleted++
// 		*xs.buf = (*xs.buf)[1:]
// 	}
// 	xs.index++
// 	return fun.Some(x)
// }

// // ScatterCopy copies stream into N streams with all source elements.
// func ScatterCopy[A any](xs Stream[A], n uint) []Stream[A] {
// 	buf := make([]bufEntry[A], 0)
// 	var mu sync.Mutex
// 	deleted := 0
// 	res := make([]Stream[A], 0, n)
// 	for i := uint(0); i < n; i++ {
// 		res = append(res, &scatterCopyImpl[A]{
// 			source:       xs,
// 			buf:          &buf,
// 			deleted:      &deleted,
// 			mu:           &mu,
// 			totalStreams: int64(n),
// 			index:        0,
// 		})
// 	}
// 	return res
// }

// Scatter splits stream into N streams, each element from source stream goes to one of the result streams.
func Scatter[A any](xs Stream[A], n uint) []Stream[A] {
	res := make([]Stream[A], 0, n)
	for i := uint(0); i < n; i++ {
		outCh := make(chan A)
		go func() {
			for x := range xs {
				outCh <- x
			}
			close(outCh)
		}()
		res = append(res, outCh)
	}
	return res
}

// // ScatterEvenly splits stream into N streams of source elements using round robin distribution
// func ScatterEvenly[A any](xs Stream[A], n uint) []Stream[A] {
// 	ch := ToChannel(xs)

// 	chans := make([]chan A, 0, n)
// 	for i := uint(0); i < n; i++ {
// 		chans = append(chans, make(chan A))
// 	}
// 	go func() {
// 		for {
// 			for _, outCh := range chans {
// 				x, ok := <-ch
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
// 		func(c chan A) Stream[A] { return FromChannel(c) },
// 	))
// }

// ScatterRoute routes elements from source stream into first matching predicate stream or last stream if non matched.
// Routes are functions from source element index and element to bool: does element match the route or not.
func ScatterRoute[A any](xs Stream[A], routes []func(uint, A) bool) []Stream[A] {
	n := len(routes)
	chans := make([]chan A, 0, n+1)
	for i := 0; i < n+1; i++ {
		chans = append(chans, make(chan A))
	}
	go func() {
		idx := uint(0)
		for {
			x, ok := <-xs
			if !ok {
				goto END
			}
			for i := 0; i < n; i++ {
				if routes[i](idx, x) {
					chans[i] <- x
					goto MATCH_FOUND
				}
			}
			chans[n] <- x
		MATCH_FOUND:
			idx++
		}
	END:
		for _, outCh := range chans {
			close(outCh)
		}
	}()

	return CollectToSlice(Map(
		FromSlice(chans),
		func(c chan A) Stream[A] { return c },
	))
}

// Gather gets elements from all input streams into single stream.
func Gather[A any](xss []Stream[A]) Stream[A] {
	res := make(chan A)
	done := make(chan fun.Unit, len(xss))
	for _, xs := range xss {
		go func(xs Stream[A]) {
			for x := range xs {
				res <- x
			}
			done <- fun.Unit1
		}(xs)
	}
	go func() {
		for range xss {
			<-done
		}
		close(res)
	}()
	return res
}
