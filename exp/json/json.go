// See for inspiration
// https://package.elm-lang.org/packages/elm/json/latest/Json-Decode
// https://package.elm-lang.org/packages/NoRedInk/elm-json-decode-pipeline/latest/Json-Decode-Pipeline
package json

import (
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/rprtr258/fun"
)

type Decoder[T any] func(any, *T) error

func (decoder Decoder[T]) ParseBytes(b []byte) (T, error) {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return *new(T), err
	}
	var t T
	if err := decoder(v, &t); err != nil {
		return t, err
	}
	return t, nil
}

func (decoder Decoder[T]) ParseString(s string) (T, error) {
	return decoder.ParseBytes([]byte(s))
}

func primitiveDecoder[T any](v any, res *T) error {
	x, ok := v.(T)
	if !ok {
		return fmt.Errorf("not a %T", x)
	}
	*res = x
	return nil
}

var Int Decoder[int] = func(v any, i *int) error {
	var f float64
	if err := primitiveDecoder(v, &f); err != nil {
		return err
	}
	if f != float64(int(f)) {
		return fmt.Errorf("not an int")
	}
	*i = int(f)
	return nil
}

var (
	String Decoder[string]    = primitiveDecoder[string]
	Bool   Decoder[bool]      = primitiveDecoder[bool]
	Float  Decoder[float64]   = primitiveDecoder[float64]
	Time   Decoder[time.Time] = func(a any, t *time.Time) error {
		x, ok := a.(string)
		if !ok {
			return fmt.Errorf("not a string")
		}
		var err error
		*t, err = time.Parse(time.RFC3339, x)
		return err
	}
)

func Any(v any, dest *any) error {
	*dest = v
	return nil
}

func Nullable[T any](decoder Decoder[T]) Decoder[fun.Option[T]] {
	return func(v any, res *fun.Option[T]) error {
		if v == nil {
			return nil
		}

		if err := decoder(v, &res.Value); err != nil {
			return nil
		}

		res.Valid = true
		return nil
	}
}

func Dict[T any](decoder Decoder[T]) Decoder[map[string]T] {
	return func(v any, res *map[string]T) error {
		vmap, ok := v.(map[string]any)
		if !ok {
			return fmt.Errorf("not a dict")
		}

		*res = make(map[string]T, len(vmap))
		for k, v := range vmap {
			var t T
			if err := decoder(v, &t); err != nil {
				return err
			}

			(*res)[k] = t
		}
		return nil
	}
}

func List[T any](decoder Decoder[T]) Decoder[[]T] {
	return func(v any, res *[]T) error {
		switch v := v.(type) {
		case nil: // parse null to nil slice
			*res = nil
		case []any:
			*res = make([]T, len(v))
			for i, v := range v {
				if err := decoder(v, &(*res)[i]); err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("not a list")
		}
		return nil
	}
}

func OneOf[T any](decoders ...Decoder[T]) Decoder[T] {
	return func(v any, res *T) error {
		errors := make([]error, len(decoders))
		for i, decoder := range decoders {
			var t T
			if err := decoder(v, &t); err == nil {
				*res = t
				return nil
			} else {
				errors[i] = err
			}
		}
		return fmt.Errorf("all variants failed: %v", errors)
	}
}

func AndThen[A, B any](da Decoder[A], f func(A) Decoder[B]) Decoder[B] {
	return func(v any, res *B) error {
		var a A
		if err := da(v, &a); err != nil {
			return err
		}
		return f(a)(v, res)
	}
}

func Success[T any](x T) Decoder[T] {
	return func(_ any, res *T) error {
		*res = x
		return nil
	}
}

func Null[T any](value T) Decoder[T] {
	return func(v any, res *T) error {
		if v != nil {
			return fmt.Errorf("not null")
		}
		*res = value
		return nil
	}
}

func Fail[T any](msg string) Decoder[T] {
	return func(any, *T) error {
		return fmt.Errorf("%s", msg)
	}
}

// Decode a Required field.
func (decoder Decoder[T]) Field(name string) Decoder[T] {
	return func(v any, res *T) error {
		vm, ok := v.(map[string]any)
		if !ok {
			return fmt.Errorf("not a dict")
		}
		v, ok = vm[name]
		if !ok {
			return fmt.Errorf("key %q not found", name)
		}

		if err := decoder(v, res); err != nil {
			return err
		}
		return nil
	}
}

func (decoder Decoder[T]) At(names []string) Decoder[T] {
	res := decoder
	for _, name := range slices.Backward(names) {
		res = res.Field(name)
	}
	return res
}

func (decoder Decoder[T]) Index(i int) Decoder[T] {
	return func(v any, res *T) error {
		vl, ok := v.([]any)
		if !ok {
			return fmt.Errorf("not a list")
		}

		if i < 0 || len(vl) <= i {
			return fmt.Errorf("no such index %d", i)
		}

		return decoder(vl[i], res)
	}
}

func (da Decoder[T]) Optional(name string, fallback T) Decoder[T] {
	return func(v any, res *T) error {
		x, ok := v.(map[string]any)
		if !ok {
			return fmt.Errorf("not a dict")
		}
		v, ok = x[name]
		if !ok {
			*res = fallback
			return nil
		}

		if err := da(v, res); err != nil {
			return err
		}
		return nil
	}
}

func Option[T any](
	name string,
	da Decoder[T],
) Decoder[fun.Option[T]] {
	return func(v any, res *fun.Option[T]) error {
		x, ok := v.(map[string]any)
		if !ok {
			return nil
		}

		v, ok = x[name]
		if !ok {
			return nil
		}

		if err := da(v, &res.Value); err != nil {
			return err
		}
		return nil
	}
}

func (d Decoder[T]) Validate(check func(T) error) Decoder[T] {
	return func(v any, res *T) error {
		if err := d(v, res); err != nil {
			return err
		}
		return check(*res)
	}
}

func Std[T any]() Decoder[T] {
	return func(v any, res *T) error {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}

		return json.Unmarshal(b, res)
	}
}
