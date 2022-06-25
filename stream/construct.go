package stream

// functions to make Stream from something that is not stream

import (
	"github.com/rprtr258/goflow/fun"
	"github.com/rprtr258/goflow/slice"
)

// Once returns a stream of one value.
func Once[A any](a A) Stream[A] {
	return FromSlice([]A{a})
}

// FromMany returns a stream with all the given values.
func FromMany[A any](as ...A) Stream[A] {
	return FromSlice(as)
}

type sliceImpl[A any] struct {
	data []A
	i    int
}

func (xs *sliceImpl[A]) Next() fun.Option[A] {
	if xs.i == len(xs.data) {
		return fun.None[A]()
	}
	xs.i++
	return fun.Some(xs.data[xs.i-1])
}

// FromSlice constructs a stream from the slice.
func FromSlice[A any](xs []A) Stream[A] {
	return &sliceImpl[A]{xs, 0}
}

// FromSet constructs stream from set elements.
func FromSet[A comparable](xs fun.Set[A]) Stream[A] {
	slice := make([]A, 0, len(xs))
	for k := range xs {
		slice = append(slice, k)
	}
	return FromSlice(slice)
}

type generateImpl[A any] struct {
	x A
	f func(A) A
}

func (xs *generateImpl[A]) Next() fun.Option[A] {
	res := xs.x
	xs.x = xs.f(xs.x)
	return fun.Some(res)
}

// Generate constructs an infinite stream of values using the production function.
func Generate[A any](x0 A, f func(A) A) Stream[A] {
	return &generateImpl[A]{x0, f}
}

type emptyImpl[A any] struct{}

func (xs *emptyImpl[A]) Next() fun.Option[A] {
	return fun.None[A]()
}

// NewStreamEmpty returns an empty stream.
func NewStreamEmpty[A any]() Stream[A] {
	return &emptyImpl[A]{}
}

// FromMap constructs Stream of key/value pairs from given map.
func FromMap[A comparable, B any](kv map[A]B) Stream[fun.Pair[A, B]] {
	return FromSlice(slice.FromMap(kv))
}

// Range makes stream of numers starting with start, step equal to step and going up to end, but not including end.
func Range[N fun.RealNumber](start, end, step N) Stream[N] {
	return TakeWhile(
		Generate(start, func(x N) N { return x + step }),
		func(x N) bool { return x < end },
	)
}
