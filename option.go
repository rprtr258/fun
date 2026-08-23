package fun

import (
	"encoding/json"
	"fmt"
)

// Option is either value or nothing.
type Option[T any] struct {
	Value T
	Valid bool
}

func (o Option[T]) String() string {
	if !o.Valid {
		return "None"
	}

	return fmt.Sprintf("Some(%v)", o.Value)
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

// Returns empty Option.
func Invalid[T any]() Option[T] {
	return Option[T]{}
}

// Returns Option with given value.
func Valid[T any](t T) Option[T] {
	return Option[T]{
		Value: t,
		Valid: true,
	}
}

// Returns Option with given value and validity.
func Optional[T any](value T, valid bool) Option[T] {
	return Option[T]{
		Value: value,
		Valid: valid,
	}
}

// Returns Option with value from pointer.
func FromPtr[T any](ptr *T) Option[T] {
	if ptr == nil {
		return Invalid[T]()
	}

	return Valid(*ptr)
}

// Returns value and validity.
func (o Option[T]) Unpack() (T, bool) {
	return o.Value, o.Valid
}

// Returns first valid Option.
func (o Option[T]) Or(other Option[T]) Option[T] {
	return IF(o.Valid, o, other)
}

// Returns value if Option is valid, otherwise returns default value.
func (o Option[T]) OrDefault(value T) T {
	return IF(o.Valid, o.Value, value)
}

// Returns pointer to value if Option is valid, otherwise returns nil.
func (opt Option[T]) Ptr() *T {
	if !opt.Valid {
		return nil
	}

	return &opt.Value
}

// Returns new Option with transformed value.
func OptMap[I, O any](o Option[I], f func(I) O) Option[O] {
	if !o.Valid {
		return Invalid[O]()
	}
	return Valid(f(o.Value))
}

// Returns new Option with transformed optional value.
func OptFlatMap[I, O any](o Option[I], f func(I) Option[O]) Option[O] {
	if !o.Valid {
		return Invalid[O]()
	}
	return f(o.Value)
}
