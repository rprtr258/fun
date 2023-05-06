package fun

import "golang.org/x/exp/constraints"

func Min[T constraints.Ordered](x, y T) T {
	return If(x < y, x, y)
}

func Max[T constraints.Ordered](x, y T) T {
	return If(x > y, x, y)
}
