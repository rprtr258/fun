package result

// functions to construct something that is not Result from Result

import (
	"github.com/rprtr258/go-flow/v2/fun"
)

// Panic is just a panic function that can be used without calling
func Panic[A, B any](a A) B {
	panic(a)
}

// Unwrap returns value if present, panics otherwise.
func (r Result[A]) Unwrap() A {
	return Fold(r, fun.Identity[A], Panic[error, A])
}

// UnwrapOr get result value or provided value if result errored.
func (r Result[A]) UnwrapOr(defaultValue A) A {
	return Fold(r, fun.Identity[A], fun.Const[A, error](defaultValue))
}

// UnwrapErr returns error if present, panics otherwise.
func (r Result[A]) UnwrapErr() error {
	return Fold(r, Panic[A, error], fun.Identity[error])
}

// IsSuccess checks if result is successful value.
func (r Result[A]) IsSuccess() bool {
	return fun.Either[A, error](r).IsLeft()
}

// IsErr checks if result is error value.
func (r Result[A]) IsErr() bool {
	return fun.Either[A, error](r).IsRight()
}

// Consume consumes value or error and executes according callback.
func (r Result[A]) Consume(fSuccess func(A), fFail func(error)) {
	fun.Either[A, error](r).Consume(fSuccess, fFail)
}

// GoResult makes Go-style result or error.
func (r Result[A]) GoResult() (A, error) {
	if r.IsErr() {
		var a A
		return a, r.UnwrapErr()
	}
	return r.Unwrap(), nil
}

// Fold constructs value by either success or fail paths.
func Fold[A, B any](mx Result[A], fSuccess func(A) B, fFail func(error) B) B {
	return fun.Fold(fun.Either[A, error](mx), fSuccess, fFail)
}
