package fun

import "cmp"

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

func Clamp[T cmp.Ordered](x, low, high T) T {
	return max(low, min(x, high))
}
