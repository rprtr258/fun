// Package result provides functions to work with Result-s
package result

import (
	"fmt"

	"github.com/rprtr258/go-flow/fun"
)

// Result represents a calculation that will yield a value of type A once executed.
// The calculation might as well fail.
// It is designed to not panic ever.
type Result[A any] fun.Either[A, error]

// ChangeErr if result is error to a given error.
func (r Result[A]) ChangeErr(err error) Result[A] {
	return Fold(r, Success[A], fun.Const[Result[A], error](Err[A](err)))
}

// TryRecover tries to recover result from error using a function that might fail.
func TryRecover[A any](ma Result[A], f func(error) Result[A]) Result[A] {
	return Fold(ma, Success[A], f)
}

// Map maps underlying success value if present, or just saves error otherwise.
func Map[A, B any](mx Result[A], f func(A) B) Result[B] {
	return fun.Fold(fun.Either[A, error](mx), fun.Compose(f, Success[B]), Err[B])
}

// FlatMap applies function to underlying successful value, or just saves error otherwise.
func FlatMap[A, B any](mx Result[A], f func(A) Result[B]) Result[B] {
	return fun.Fold(fun.Either[A, error](mx), f, Err[B])
}

// FlatMap2 uses FlatMap two times.
func FlatMap2[A, B, C any](mx Result[A], f1 func(A) Result[B], f2 func(B) Result[C]) Result[C] {
	return FlatMap(FlatMap(mx, f1), f2)
}

// FlatMap3 uses FlatMap three times.
func FlatMap3[A, B, C, D any](mx Result[A], f1 func(A) Result[B], f2 func(B) Result[C], f3 func(C) Result[D]) Result[D] {
	return FlatMap(FlatMap2(mx, f1, f2), f3)
}

// WrapErrf wraps an error with additional context information
func WrapErrf[A any](io Result[A], format string, args ...interface{}) Result[A] {
	return TryRecover(io, func(err error) Result[A] {
		argv := make([]any, 0)
		argv = append(argv, args...)
		argv = append(argv, err)
		return Err[A](fmt.Errorf(format+": %w", argv...))
	})
}
