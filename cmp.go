package fun

import "cmp"

// Min returns the minimum of the given values
func Min[T cmp.Ordered](xs ...T) T {
	if len(xs) == 0 {
		return Zero[T]()
	}

	res := xs[0]
	for _, x := range xs[1:] {
		res = min(res, x)
	}
	return res
}

// Max returns the maximum of the given values
func Max[T cmp.Ordered](xs ...T) T {
	if len(xs) == 0 {
		return Zero[T]()
	}

	res := xs[0]
	for _, x := range xs[1:] {
		res = max(res, x)
	}
	return res
}

// Clamp returns x clamped between low and high
func Clamp[T cmp.Ordered](x, low, high T) T {
	return max(low, min(x, high))
}

// MinBy returns the minimum of the given values using the given order function
func MinBy[T any, R cmp.Ordered](order func(T) R, xs ...T) T {
	if len(xs) == 0 {
		return Zero[T]()
	}

	res := xs[0]
	fres := order(res)
	for _, x := range xs[1:] {
		if fx := order(x); fres > fx {
			res, fres = x, fx
		}
	}
	return res
}

// MaxBy returns the maximum of the given values using the given order function
func MaxBy[T any, R cmp.Ordered](order func(T) R, xs ...T) T {
	if len(xs) == 0 {
		return Zero[T]()
	}

	res := xs[0]
	fres := order(res)
	for _, x := range xs[1:] {
		if fx := order(x); fres < fx {
			res, fres = x, fx
		}
	}
	return res
}
