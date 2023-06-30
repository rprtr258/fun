// Package stream provides a way to construct data processing streams from smaller pieces.
package stream

import (
	"errors"
	"fmt"

	"github.com/rprtr258/go-flow/v2/fun"
)

var EOF = errors.New("stream ended")

// Stream is a finite or infinite stream of values.
type Stream[T any] interface {
	// Next returns the next value in stream. Returned error is:
	// - nil if value is got with no errors.
	// - EOF is returned if stream has ended.
	// - error if some error happened.
	Next() (T, error)
}

type Iterator[T any] interface {
	// Next returns next value and false, if iterator was not ended.
	// Returns zero value and true otherwise.
	Next() (T, bool)
	// Err returns first error happened during Next calls.
	Err() error
}

// Iterator2 is iterator taking function which handles iteratee and returns
// whether to stop iteration.
type Iterator2[T any] func(func(T, error) bool)

type StreamFunc[T any] func() (T, error)

func (f StreamFunc[T]) Next() (T, error) {
	return f()
}

// Map converts values of the stream.
func Map[A, B any](xs Stream[A], f func(A) B) Stream[B] {
	return StreamFunc[B](func() (B, error) {
		x, err := xs.Next()
		if err == EOF {
			var b B
			return b, EOF
		}

		return f(x), err
	})
}

// Chain appends another stream after the end of the first one.
func Chain[A any](xss ...Stream[A]) Stream[A] {
	var zero A
	i := 0
	return StreamFunc[A](func() (A, error) {
		if i == len(xss) {
			return zero, EOF
		}

		for ; i < len(xss); i++ {
			x, err := xss[i].Next()
			switch err {
			case nil:
				return x, nil
			case EOF:
				continue
			default:
				return x, err
			}
		}

		return zero, EOF
	})
}

// FlatMap maps stream using function and concatenates result streams into one.
func FlatMap[A, B any](xs Stream[A], f func(A) Stream[B]) Stream[B] {
	var zero B
	var batch Stream[B]
	return StreamFunc[B](func() (B, error) {
		for {
			if batch != nil {
				x, err := batch.Next()
				switch err {
				case nil:
					return x, nil
				case EOF:
				default:
					return zero, err
				}
			}

			x, err := xs.Next()
			switch err {
			case nil:
				batch = f(x)
			case EOF:
				return zero, EOF
			default:
				return zero, err
			}
		}
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
func Chunked[A any](xs Stream[A], n int) Stream[[]A] {
	if n <= 0 {
		panic(fmt.Sprintf("Chunk must be of positive size, but %d given", n))
	}

	return StreamFunc[[]A](func() ([]A, error) {
		res := make([]A, 0, n)
		for i := 0; i < n; i++ {
			x, err := xs.Next()
			switch err {
			case nil:
				res = append(res, x)
			case EOF:
				if len(res) == 0 {
					return nil, EOF
				}

				return res, nil
			default:
				return res, err
			}
		}

		return res, nil
	})
}

// Intersperse adds a separator after each stream element.
func Intersperse[A any](xs Stream[A], sep A) Stream[A] {
	isSep := false
	ended := false
	var zero A
	return StreamFunc[A](func() (A, error) {
		if ended {
			return zero, EOF
		}

		if isSep {
			isSep = false
			return sep, nil
		}

		isSep = true
		switch x, err := xs.Next(); err {
		case nil:
			return x, nil
		case EOF:
			ended = true
			isSep = false
			return zero, EOF
		default:
			return zero, err
		}
	})
}

// Repeat appends the same stream infinitely.
func Repeat[A any](xs Stream[A]) Stream[A] {
	buf := []A{}
	var ierr error
X:
	for {
		x, err := xs.Next()
		switch err {
		case nil:
			buf = append(buf, x)
		case EOF:
			break X
		default:
			ierr = err
			break X
		}
	}

	var zero A
	i := 0
	return StreamFunc[A](func() (A, error) {
		if ierr != nil {
			return zero, ierr
		}

		res := buf[i]
		i = (i + 1) % len(buf)
		return res, nil
	})
}

// Take cuts the stream after n elements.
func Take[A any](xs Stream[A], n int) Stream[A] {
	if n < 0 {
		panic(fmt.Sprintf("Cannot take negative number of elements, concretely %d", n))
	}

	took := 0
	var zero A
	return StreamFunc[A](func() (A, error) {
		if took == n {
			return zero, EOF
		}

		took++
		return xs.Next()
	})
}

// Skip skips n elements in the stream.
func Skip[A any](xs Stream[A], n int) Stream[A] {
	for i := 0; i < n; i++ {
		xs.Next()
	}

	return xs
}

// Filter leaves in the stream only the elements that satisfy the given predicate.
func Filter[A any](xs Stream[A], p func(A) bool) Stream[A] {
	var zero A
	return StreamFunc[A](func() (A, error) {
		for {
			switch x, err := xs.Next(); err {
			case nil:
				if p(x) {
					return x, nil
				}
			case EOF:
				return zero, EOF
			default:
				return zero, err
			}
		}
	})
}

// Find searches for first element matching the predicate.
func Find[A any](xs Stream[A], p func(A) bool) fun.Option[A] {
	x, err := Filter(xs, p).Next()
	return fun.FromNull(x, err == nil)
}

// TakeWhile takes elements while predicate is true.
func TakeWhile[A any](xs Stream[A], p func(A) bool) Stream[A] {
	var zero A
	ended := false
	return StreamFunc[A](func() (A, error) {
		if ended {
			return zero, EOF
		}

		for {
			switch x, err := xs.Next(); err {
			case nil:
				if p(x) {
					return x, nil
				} else {
					ended = true
					return zero, EOF
				}
			case EOF:
				return zero, EOF
			default:
				return zero, err
			}
		}
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
	var zero B
	return StreamFunc[B](func() (B, error) {
		for {
			switch x, err := xs.Next(); err {
			case nil:
				y, ok := f(x).Unpack()
				if ok {
					return y, nil
				}
			case EOF:
				return zero, EOF
			default:
				return zero, err
			}
		}
	})
}

// Paged makes stream from stream of pages of elements represented as slices.
func Paged[A any](xs Stream[[]A]) Stream[A] {
	return FlatMap(xs, FromSlice[A])
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
