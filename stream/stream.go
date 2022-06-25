// Package stream provides a way to construct data processing streams from smaller pieces.
package stream

// TODO: see if returning structs but not interfaces is faster

import (
	"sync"
	"sync/atomic"

	"github.com/rprtr258/goflow/fun"
)

// Stream is a finite or infinite stream of values.
type Stream[A any] interface {
	// Next gives either value or nothing if stream has ended.
	Next() fun.Option[A]
}

type chanStream[A any] interface {
	channel() <-chan A
}

type mapImpl[A, B any] struct {
	Stream[A]
	f func(A) B
}

func (xs *mapImpl[A, B]) Next() fun.Option[B] {
	return fun.Map(xs.Stream.Next(), xs.f)
}

// Map converts values of the stream.
func Map[A, B any](xs Stream[A], f func(A) B) Stream[B] {
	return &mapImpl[A, B]{xs, f}
}

type chainImpl[A any] []Stream[A]

func (xs *chainImpl[A]) Next() fun.Option[A] {
	for len(*xs) > 0 {
		x := (*xs)[0].Next()
		if x.IsSome() {
			return x
		}
		*xs = (*xs)[1:]
	}
	return fun.None[A]()
}

// Chain appends another stream after the end of the first one.
func Chain[A any](xss ...Stream[A]) Stream[A] {
	res := chainImpl[A](xss)
	return &res
}

type flatMapImpl[A, B any] struct {
	Stream[A]
	f    func(A) Stream[B]
	last Stream[B]
}

func (xs *flatMapImpl[A, B]) Next() fun.Option[B] {
	y := xs.last.Next()
	if y.IsNone() {
		xs.last = fun.FoldOption(fun.Map(xs.Stream.Next(), xs.f), fun.Identity[Stream[B]], NewStreamEmpty[B])
		y = xs.last.Next()
	}
	return y
}

// FlatMap maps stream using function and concatenates result streams into one.
func FlatMap[A, B any](xs Stream[A], f func(A) Stream[B]) Stream[B] {
	return &flatMapImpl[A, B]{xs, f, NewStreamEmpty[B]()}
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

type chunkedImpl[A any] struct {
	Stream[A]
	chunkSize int
}

func (xs *chunkedImpl[A]) Next() fun.Option[[]A] {
	x := xs.Stream.Next()
	if x.IsNone() {
		return fun.None[[]A]()
	}
	chunk := make([]A, 1, xs.chunkSize)
	chunk[0] = x.Unwrap()
	for i := 1; i < xs.chunkSize; i++ {
		x := xs.Stream.Next()
		if x.IsNone() {
			break
		}
		chunk = append(chunk, x.Unwrap())
	}
	return fun.Some(chunk)
}

// Chunked groups elements by n and produces a stream of slices.
func Chunked[A any](xs Stream[A], n int) Stream[[]A] {
	return &chunkedImpl[A]{xs, n}
}

type intersperseImpl[A any] struct {
	Stream[A]
	sep       A
	nextElem  fun.Option[A]
	isSepNext bool
	isFirst   bool
}

func (xs *intersperseImpl[A]) Next() fun.Option[A] {
	switch {
	case xs.isFirst:
		xs.isFirst = false
		x := xs.Stream.Next()
		xs.nextElem = xs.Stream.Next()
		xs.isSepNext = true
		return x
	case xs.isSepNext && xs.nextElem.IsSome():
		xs.isSepNext = false
		return fun.Some(xs.sep)
	default:
		var x fun.Option[A]
		x, xs.nextElem = xs.nextElem, xs.Stream.Next()
		xs.isSepNext = true
		return x
	}
}

// Intersperse adds a separator after each stream element.
func Intersperse[A any](xs Stream[A], sep A) Stream[A] {
	return &intersperseImpl[A]{
		Stream:    xs,
		sep:       sep,
		nextElem:  fun.None[A](),
		isFirst:   true,
		isSepNext: false,
	}
}

type repeatImpl[A any] struct {
	Stream[A]
	i   int
	buf []A
}

func (xs *repeatImpl[A]) Next() fun.Option[A] {
	x := xs.Stream.Next()
	if x.IsNone() {
		res := xs.buf[xs.i]
		xs.i = (xs.i + 1) % len(xs.buf)
		return fun.Some(res)
	}
	xs.buf = append(xs.buf, x.Unwrap())
	return x
}

// Repeat appends the same stream infinitely.
func Repeat[A any](xs Stream[A]) Stream[A] {
	return &repeatImpl[A]{xs, 0, make([]A, 0)}
}

type takeImpl[A any] struct {
	Stream[A]
	n uint
}

func (xs *takeImpl[A]) Next() fun.Option[A] {
	if xs.n == 0 {
		return fun.None[A]()
	}
	xs.n--
	return xs.Stream.Next()
}

// Take cuts the stream after n elements.
func Take[A any](xs Stream[A], n uint) Stream[A] {
	return &takeImpl[A]{xs, n}
}

// Skip skips n elements in the stream.
func Skip[A any](xs Stream[A], n int) Stream[A] {
	for i := 0; i < n; i++ {
		if x := xs.Next(); x.IsNone() {
			break
		}
	}
	return xs
}

type filterImpl[A any] struct {
	Stream[A]
	p func(A) bool
}

func (xs *filterImpl[A]) Next() fun.Option[A] {
	for {
		x := xs.Stream.Next()
		if x.IsNone() {
			break
		}
		if xs.p(x.Unwrap()) {
			return x
		}
	}
	return fun.None[A]()
}

// Filter leaves in the stream only the elements that satisfy the given predicate.
func Filter[A any](xs Stream[A], p func(A) bool) Stream[A] {
	return &filterImpl[A]{xs, p}
}

// Find searches for first element matching the predicate.
func Find[A any](xs Stream[A], p func(A) bool) fun.Option[A] {
	return Filter(xs, p).Next()
}

type takeWhileImpl[A any] struct {
	Stream[A]
	p     func(A) bool
	ended bool
}

func (xs *takeWhileImpl[A]) Next() fun.Option[A] {
	if xs.ended {
		return fun.None[A]()
	}
	if x := xs.Stream.Next(); x.IsSome() && xs.p(x.Unwrap()) {
		return x
	}
	xs.ended = true
	return fun.None[A]()
}

// TakeWhile takes elements while predicate is true.
func TakeWhile[A any](xs Stream[A], p func(A) bool) Stream[A] {
	return &takeWhileImpl[A]{xs, p, false}
}

// DebugPrint prints every processed element, without changing it.
func DebugPrint[A any](prefix string, xs Stream[A]) Stream[A] {
	return Map(xs, fun.Debug[A](prefix))
}

type uniqueImpl[A comparable] struct {
	Stream[A]
	seen fun.Set[A]
}

func (xs *uniqueImpl[A]) Next() fun.Option[A] {
	for {
		x := xs.Stream.Next()
		if x.IsNone() {
			return fun.None[A]()
		}
		xVal := x.Unwrap()
		if !xs.seen.Contains(xVal) {
			xs.seen[xVal] = fun.Unit1
			return x
		}
	}
}

// Unique makes stream of unique elements.
func Unique[A comparable](xs Stream[A]) Stream[A] {
	return &uniqueImpl[A]{xs, make(fun.Set[A])}
}

type mapFilterImpl[A, B any] struct {
	Stream[A]
	f func(A) fun.Option[B]
}

func (xs *mapFilterImpl[A, B]) Next() fun.Option[B] {
	for {
		x := xs.Stream.Next()
		if x.IsNone() {
			return fun.None[B]()
		}
		y := xs.f(x.Unwrap())
		if y.IsSome() {
			return y
		}
	}
}

// MapFilter applies function to every element and leaves only elements that are not None.
func MapFilter[A, B any](xs Stream[A], f func(A) fun.Option[B]) Stream[B] {
	return &mapFilterImpl[A, B]{xs, f}
}

// Paged makes stream from stream of pages of elements represented as slices.
func Paged[A any](xs Stream[[]A]) Stream[A] {
	return FlatMap(xs, FromSlice[A])
}

type bufEntry[A any] struct {
	value       A
	streamsLeft int64
}

type scatterCopyImpl[A any] struct {
	source       Stream[A]
	buf          *[]bufEntry[A]
	deleted      *int
	mu           *sync.Mutex
	totalStreams int64

	index int
}

func (xs *scatterCopyImpl[A]) advanceBuffer() {
	if xs.index-*xs.deleted == len(*xs.buf) {
		x := xs.source.Next()
		if x.IsNone() {
			return
		}
		bufEntry := bufEntry[A]{
			value:       x.Unwrap(),
			streamsLeft: xs.totalStreams,
		}
		if len(*xs.buf) == cap(*xs.buf) {
			// reallocate buffer
			buf := append(*xs.buf, bufEntry)
			*xs.buf = buf
		} else {
			// insert into buffer without realloc, **xs.buf is not changed
			*xs.buf = append(*xs.buf, bufEntry)
		}
	}
}

func (xs *scatterCopyImpl[A]) Next() fun.Option[A] {
	xs.mu.Lock()
	defer xs.mu.Unlock()

	xs.advanceBuffer()

	effectiveIndex := xs.index - *xs.deleted
	if effectiveIndex == len(*xs.buf) {
		return fun.None[A]()
	}
	atomic.AddInt64(&(*xs.buf)[effectiveIndex].streamsLeft, -1)
	x := (*xs.buf)[effectiveIndex].value
	if (*xs.buf)[effectiveIndex].streamsLeft == 0 {
		*xs.deleted++
		*xs.buf = (*xs.buf)[1:]
	}
	xs.index++
	return fun.Some(x)
}

// ScatterCopy copies stream into N streams with all source elements.
func ScatterCopy[A any](xs Stream[A], n uint) []Stream[A] {
	buf := make([]bufEntry[A], 0)
	var mu sync.Mutex
	deleted := 0
	res := make([]Stream[A], 0, n)
	for i := uint(0); i < n; i++ {
		res = append(res, &scatterCopyImpl[A]{
			source:       xs,
			buf:          &buf,
			deleted:      &deleted,
			mu:           &mu,
			totalStreams: int64(n),
			index:        0,
		})
	}
	return res
}

// Scatter splits stream into N streams, each element from source stream goes to one of the result streams.
func Scatter[A any](xs Stream[A], n uint) []Stream[A] {
	ch := ToChannel(xs)
	res := make([]Stream[A], 0, n)
	for i := uint(0); i < n; i++ {
		outCh := make(chan A)
		go func() {
			for x := range ch {
				outCh <- x
			}
			close(outCh)
		}()
		res = append(res, FromChannel(outCh))
	}
	return res
}

// ScatterEvenly splits stream into N streams of source elements using round robin distribution
func ScatterEvenly[A any](xs Stream[A], n uint) []Stream[A] {
	ch := ToChannel(xs)

	chans := make([]chan A, 0, n)
	for i := uint(0); i < n; i++ {
		chans = append(chans, make(chan A))
	}
	go func() {
		for {
			for _, outCh := range chans {
				x, ok := <-ch
				if !ok {
					goto END
				}
				outCh <- x
			}
		}
	END:
		for _, outCh := range chans {
			close(outCh)
		}
	}()

	return CollectToSlice(Map(
		FromSlice(chans),
		func(c chan A) Stream[A] { return FromChannel(c) },
	))
}

// ScatterRoute routes elements from source stream into first matching predicate stream or last stream if non matched.
// Routes are functions from source element index and element to bool: does element match the route or not.
func ScatterRoute[A any](xs Stream[A], routes []func(uint, A) bool) []Stream[A] {
	ch := ToChannel(xs)

	n := len(routes)
	chans := make([]chan A, 0, n+1)
	for i := 0; i < n+1; i++ {
		chans = append(chans, make(chan A))
	}
	go func() {
		idx := uint(0)
		for {
			x, ok := <-ch
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
		func(c chan A) Stream[A] { return FromChannel(c) },
	))
}

// Gather gets elements from all input streams into single stream.
func Gather[A any](xss []Stream[A]) Stream[A] {
	ch := make(chan A)
	done := make(chan fun.Unit, len(xss))
	for _, xs := range xss {
		go func(xs Stream[A]) {
			chIn := ToChannel(xs)
			for x := range chIn {
				ch <- x
			}
			done <- fun.Unit1
		}(xs)
	}
	go func() {
		for range xss {
			<-done
		}
		close(ch)
	}()
	return FromChannel(ch)
}
