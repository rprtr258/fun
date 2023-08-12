package fun

import "encoding/json"

type Option[T any] struct {
	value T
	valid bool
}

func (o Option[T]) MarshalJSON() ([]byte, error) {
	if !o.valid {
		return []byte("null"), nil
	}

	return json.Marshal(o.value)
}

func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.valid = false
		return nil
	}

	o.valid = true
	return json.Unmarshal(data, &o.value)
}

func Invalid[T any]() Option[T] {
	return Option[T]{}
}

func Valid[T any](t T) Option[T] {
	return Option[T]{
		value: t,
		valid: true,
	}
}

func Optional[T any](value T, valid bool) Option[T] {
	return Option[T]{
		valid: valid,
		value: value,
	}
}

func (o Option[T]) Valid() bool {
	return o.valid
}

func (o Option[T]) Unwrap() T {
	return o.value
}

func (o Option[T]) Unpack() (T, bool) {
	return o.value, o.valid
}

func (o Option[T]) Or(other Option[T]) Option[T] {
	return If(o.valid, o, other)
}

func (o Option[T]) OrDefault(value T) T {
	return If(o.valid, o.value, value)
}

func FromPtr[T any](ptr *T) Option[T] {
	if ptr == nil {
		return Invalid[T]()
	}

	return Valid(*ptr)
}

func (opt Option[T]) Ptr() *T {
	if !opt.valid {
		return nil
	}

	return &opt.value
}

func OptMap[T, R any](o Option[T], f func(T) R) Option[R] {
	if !o.valid {
		return Invalid[R]()
	}
	return Valid(f(o.value))
}

func OptFlatMap[T, R any](o Option[T], f func(T) Option[R]) Option[R] {
	if !o.valid {
		return Invalid[R]()
	}
	return f(o.value)
}
