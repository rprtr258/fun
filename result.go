package fun

import (
	"encoding/json"
	"fmt"
)

// Result represents a calculation that will yield a value of type A once executed.
// The calculation might as well fail.
// It is designed to not panic ever.
type Result[T any] struct {
	Value T
	Err   error
}

func Ok[T any](value T) Result[T] {
	return Result[T]{value, nil}
}

func Err[T any](err error) Result[T] {
	return Result[T]{*new(T), err}
}

func Try[T any](f func() (T, error)) Result[T] {
	res, err := f()
	return Result[T]{res, err}
}

func (r Result[T]) String() string {
	if r.Err != nil {
		return fmt.Sprintf("Err(%s)", r.Err.Error())
	}

	return fmt.Sprintf("Ok(%v)", r.Value)
}

func (r Result[T]) MarshalJSON() ([]byte, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	return json.Marshal(r.Value)
}

func (r *Result[T]) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &r.Value)
	r.Err = err
	return nil
}

func (r Result[T]) Unwrap() T {
	if r.Err != nil {
		panic(r.Err)
	}
	return r.Value
}

func (r Result[T]) Unpack() (T, error) {
	return r.Value, r.Err
}

func (r Result[T]) Or(other Result[T]) Result[T] {
	return IF(r.Err == nil, r, other)
}

func (r Result[T]) OrDefault(value T) T {
	return IF(r.Err == nil, r.Value, value)
}

func (r Result[T]) Get() Option[T] {
	return Optional(r.Value, r.Err == nil)
}

func (r Result[T]) Map[R any](f func(T) R) Result[R] {
	if r.Err != nil {
		return Result[R]{*new(R), r.Err}
	}
	return Result[R]{f(r.Value), nil}
}

func (r Result[T]) FlatMap[R any](f func(T) Result[R]) Result[R] {
	if r.Err != nil {
		return Result[R]{*new(R), r.Err}
	}
	return f(r.Value)
}

func (r Result[T]) ReifyErr(f func(error) error) Result[T] {
	if r.Err != nil {
		return Result[T]{r.Value, f(r.Err)}
	}
	return r
}
