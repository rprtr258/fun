package result

// functions to construct Result from something that is not Result

import (
	"fmt"

	"github.com/rprtr258/go-flow/fun"
)

// ErrorNilDeref is an error on nil pointer dereference
var ErrorNilDeref = fmt.Errorf("nil pointer dereference")

// Success constructs a successfule Result[A].
func Success[A any](a A) Result[A] {
	return Result[A](fun.Left[A, error](a))
}

// Err constructs failed Result[A].
func Err[A any](err error) Result[A] {
	return Result[A](fun.Right[A](err))
}

// FromGoResult constructs Result from Go result/error pair.
func FromGoResult[A any](a A, err error) Result[A] {
	if err != nil {
		return Err[A](err)
	}
	return Success(a)
}

// FromMaybe constructs Result from result/isValid pair
func FromMaybe[A any](a A, valid bool) Result[A] {
	if !valid {
		return Err[A](fmt.Errorf("invalid value"))
	}
	return Success(a)
}

// Eval constructs Result from a function that might fail.
// If there is panic in the function, it's recovered from and represented as an error.
func Eval[A any](f func() (A, error)) Result[A] {
	var (
		a   A
		err error
	)
	defer recoverToError(&err)
	a, err = f()
	return FromGoResult(a, err)
}

// ToKleisli converts go function returning result or error to function returning Result.
func ToKleisli[A, B any](f func(A) (B, error)) func(A) Result[B] {
	return func(a A) Result[B] {
		return FromGoResult(f(a))
	}
}

// Dereference retrieves the value by pointer. Fails if pointer is nil.
func Dereference[A any](ptra *A) Result[A] {
	if ptra == nil {
		return Err[A](ErrorNilDeref)
	}
	return Success(*ptra)
}

// recoverToError recovers and places the recovered error into the given variable.
func recoverToError(err *error) {
	if err2 := recover(); err2 != nil {
		if err != nil {
			*err = fmt.Errorf("recovered (err: %w) from: %v", *err, err2)
		} else {
			*err = fmt.Errorf("recovered from: %v", err2)
		}
	}
}
