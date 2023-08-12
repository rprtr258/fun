package fun

import "golang.org/x/exp/constraints"

func Min[T constraints.Ordered](xs ...T) T {
	if len(xs) == 0 {
		return Zero[T]()
	}

	min := xs[0]
	for i := 1; i < len(xs); i++ {
		min = If(xs[i] < min, xs[i], min)
	}
	return min
}

func Max[T constraints.Ordered](xs ...T) T {
	if len(xs) == 0 {
		return Zero[T]()
	}

	max := xs[0]
	for i := 1; i < len(xs); i++ {
		max = If(xs[i] > max, xs[i], max)
	}
	return max
}

func Clamp[T constraints.Ordered](x, min, max T) T {
	if x > max {
		return max
	}
	if x < min {
		return min
	}
	return x
}
