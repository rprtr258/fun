package fun

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// Option is either value or nothing.
type Option[T any] struct {
	Value T
	Valid bool
}

// Returns Option with given value and validity.
func Optional[T any](value T, valid bool) Option[T] {
	return Option[T]{value, valid}
}

// Returns empty Option.
func Invalid[T any]() Option[T] { return Optional(*new(T), false) }

// Returns Option with given value.
func Valid[T any](t T) Option[T] { return Optional(t, true) }

// Returns Option with value from pointer.
func FromPtr[T any](ptr *T) Option[T] {
	if ptr == nil {
		return Invalid[T]()
	}

	return Valid(*ptr)
}

func FromMapGet[K comparable, V any](m map[K]V, key K) Option[V] {
	value, ok := m[key]
	return Optional(value, ok)
}

func FromSqlNull[T any](v sql.Null[T]) Option[T]           { return Optional(v.V, v.Valid) }
func FromSqlNullBool(v sql.NullBool) Option[bool]          { return Optional(v.Bool, v.Valid) }
func FromSqlNullString(v sql.NullString) Option[string]    { return Optional(v.String, v.Valid) }
func FromSqlNullByte(v sql.NullByte) Option[byte]          { return Optional(v.Byte, v.Valid) }
func FromSqlNullFloat64(v sql.NullFloat64) Option[float64] { return Optional(v.Float64, v.Valid) }
func FromSqlNullTime(v sql.NullTime) Option[time.Time]     { return Optional(v.Time, v.Valid) }
func FromSqlNullInt16(v sql.NullInt16) Option[int16]       { return Optional(v.Int16, v.Valid) }
func FromSqlNullInt32(v sql.NullInt32) Option[int32]       { return Optional(v.Int32, v.Valid) }
func FromSqlNullInt64(v sql.NullInt64) Option[int64]       { return Optional(v.Int64, v.Valid) }

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

func (o Option[T]) Unwrap() T {
	if !o.Valid {
		panic("tried to Unwrap() Invalid value")
	}
	return o.Value
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
func (o Option[T]) Ptr() *T {
	if !o.Valid {
		return nil
	}

	return &o.Value
}

// Returns new Option with transformed value.
func (o Option[T]) Map[R any](f func(T) R) Option[R] {
	if !o.Valid {
		return Invalid[R]()
	}
	return Valid(f(o.Value))
}

// Returns new Option with transformed optional value.
func (o Option[T]) FlatMap[R any](f func(T) Option[R]) Option[R] {
	if !o.Valid {
		return Invalid[R]()
	}
	return f(o.Value)
}

func (o Option[T]) All(yield func(T) bool) {
	if o.Valid {
		yield(o.Value)
	}
}
