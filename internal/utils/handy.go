package utils

import "golang.org/x/exp/constraints"

func Clamp[T constraints.Ordered](v, x, y T) T {
	if v < x {
		return x
	} else if v > y {
		return y
	}
	return v
}
