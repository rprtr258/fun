package fun

import "golang.org/x/exp/constraints"

func Min[T constraints.Ordered](xs ...T) T {
	if len(xs) == 0 {
		return Zero[T]()
	}

	min := xs[0]
	for i := 1; i < len(xs); i++ {
		if xs[i] < min {
			min = xs[i]
		}
	}
	return min
}

func Max[T constraints.Ordered](xs ...T) T {
	if len(xs) == 0 {
		return Zero[T]()
	}

	max := xs[0]
	for i := 1; i < len(xs); i++ {
		if xs[i] > max {
			max = xs[i]
		}
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
