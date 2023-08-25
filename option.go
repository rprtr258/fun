package fun

import "encoding/json"

// Option is either value or nothing.
type Option[T any] struct {
	Value T
	Valid bool
}

func (o Option[T]) MarshalJSON() ([]byte, error) {
	if !o.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(o.Value)
}

func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.Valid = false
		return nil
	}

	o.Valid = true
	return json.Unmarshal(data, &o.Value)
}

func Invalid[T any]() Option[T] {
	return Option[T]{}
}

func Valid[T any](t T) Option[T] {
	return Option[T]{
		Value: t,
		Valid: true,
	}
}

func Optional[T any](value T, valid bool) Option[T] {
	return Option[T]{
		Value: value,
		Valid: valid,
	}
}

func (o Option[T]) Unpack() (T, bool) {
	return o.Value, o.Valid
}

func (o Option[T]) Or(other Option[T]) Option[T] {
	return IF(o.Valid, o, other)
}

func (o Option[T]) OrDefault(value T) T {
	return IF(o.Valid, o.Value, value)
}

func FromPtr[T any](ptr *T) Option[T] {
	if ptr == nil {
		return Invalid[T]()
	}

	return Valid(*ptr)
}

func (opt Option[T]) Ptr() *T {
	if !opt.Valid {
		return nil
	}

	return &opt.Value
}

func OptMap[I, O any](o Option[I], f func(I) O) Option[O] {
	if !o.Valid {
		return Invalid[O]()
	}
	return Valid(f(o.Value))
}

func OptFlatMap[I, O any](o Option[I], f func(I) Option[O]) Option[O] {
	if !o.Valid {
		return Invalid[O]()
	}
	return f(o.Value)
}
