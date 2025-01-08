package fun

import (
	"fmt"
	"log"
)

// Pair is a data structure that has two values.
type Pair[K, V any] struct {
	K K
	V V
}

func Zero[T any]() T {
	var zero T
	return zero
}

// ToString converts the value to string.
func ToString[A any](a A) string {
	return fmt.Sprintf("%v", a)
}

// DebugP returns function that prints prefix with element and returns it.
// Useful for debug printing.
func DebugP[V any](prefix string) func(V) V {
	return func(v V) V {
		// TODO: print place
		log.Println(prefix, v)
		return v
	}
}

// Debug returns function that prints element and returns it.
// Useful for debug printing.
func Debug[V any](v V) V {
	// TODO: print place
	log.Println(v)
	return v
}

func Has[K comparable, V any](dict map[K]V, key K) bool {
	_, ok := dict[key]
	return ok
}

func Cond[R any](defaultValue R, cases ...func() (R, bool)) R {
	for _, case_ := range cases {
		if res, ok := case_(); ok {
			return res
		}
	}

	return defaultValue
}

func Ptr[T any](t T) *T {
	return &t
}

func Deref[T any](ptr *T) T {
	if ptr == nil {
		return Zero[T]()
	}
	return *ptr
}

func Pipe[T any](t T, endos ...func(T) T) T {
	for _, endo := range endos {
		t = endo(t)
	}
	return t
}
