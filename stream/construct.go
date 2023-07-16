package stream

// functions to make Stream from something that is not stream

import (
	"github.com/rprtr258/go-flow/v2/fun"
	"github.com/rprtr258/go-flow/v2/slice"
)

// Once returns a stream of one value.
func Once[A any](a A) Stream[A] {
	return FromMany(a)
}

// FromMany returns a stream with all the given values.
func FromMany[A any](as ...A) Stream[A] {
	return FromSlice(as)
}

// FromSlice constructs a stream from the slice.
func FromSlice[A any](xs []A) Stream[A] {
	i := 0
	var zero A
	return StreamFunc[A](func() (A, error) {
		if i == len(xs) {
			return zero, EOS
		}

		i++
		return xs[i-1], nil
	})
}

// FromSet constructs stream from set elements.
func FromSet[A comparable](xs fun.Set[A]) Stream[A] {
	slice := make([]A, 0, len(xs))
	for k := range xs {
		slice = append(slice, k)
	}
	return FromSlice(slice)
}

// Generate constructs an infinite stream of values using the production function.
func Generate[A any](x0 A, f func(A) A) Stream[A] {
	cur := x0
	return StreamFunc[A](func() (A, error) {
		res := cur
		cur = f(cur)
		return res, nil
	})
}

// NewStreamEmpty returns an empty stream.
func NewStreamEmpty[A any]() Stream[A] {
	return FromSlice[A](nil)
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
