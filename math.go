package fun

import "golang.org/x/exp/constraints"

func Min[T constraints.Ordered](x, y T) T {
	return If(x < y, x, y)
}

func Max[T constraints.Ordered](x, y T) T {
	return If(x > y, x, y)
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
