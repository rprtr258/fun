// See for inspiration
// https://package.elm-lang.org/packages/elm/json/latest/Json-Decode
// https://package.elm-lang.org/packages/NoRedInk/elm-json-decode-pipeline/latest/Json-Decode-Pipeline
package json

import (
	"encoding/json"
	"fmt"
	"time"
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
var String Decoder[string] = primitiveDecoder[string]
var Bool Decoder[bool] = primitiveDecoder[bool]
var Float Decoder[float64] = primitiveDecoder[float64]
var Time Decoder[time.Time] = func(a any, t *time.Time) error {
	x, ok := a.(string)
	if !ok {
		return fmt.Errorf("not a string")
	}
	var err error
	*t, err = time.Parse(time.RFC3339, x)
	return err
}

type Maybe[T any] struct {
	Value T
	Valid bool
}

func Nullable[T any](decoder Decoder[T]) Decoder[Maybe[T]] {
	return func(v any, res *Maybe[T]) error {
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
		vl, ok := v.([]any)
		if !ok {
			return fmt.Errorf("not a list")
		}

		*res = make([]T, len(vl))
		for i, v := range vl {
			var t T
			if err := decoder(v, &t); err != nil {
				return err
			}
			(*res)[i] = t
		}
		return nil
	}
}

// TODO: expressible as required
func Field[T any](name string, decoder Decoder[T]) Decoder[T] {
	return Required(name, decoder)
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

func Map[T, R any](decoder Decoder[T], f func(T) R) Decoder[R] {
	return func(v any, res *R) error {
		var t T
		if err := decoder(v, &t); err != nil {
			return nil
		}

		*res = f(t)
		return nil
	}
}

func Map2[A, B, T any](
	combine func(A, B) T,
	da Decoder[A],
	db Decoder[B],
) Decoder[T] {
	return func(v any, res *T) error {
		var aa A
		if err := da(v, &aa); err != nil {
			return err
		}

		var bb B
		if err := db(v, &bb); err != nil {
			return err
		}

		*res = combine(aa, bb)
		return nil
	}
}

func Map3[A, B, C, T any](
	combine func(A, B, C) T,
	da Decoder[A],
	db Decoder[B],
	dc Decoder[C],
) Decoder[T] {
	return func(v any, res *T) error {
		var aa A
		if err := da(v, &aa); err != nil {
			return err
		}
		var bb B
		if err := db(v, &bb); err != nil {
			return err
		}
		var cc C
		if err := dc(v, &cc); err != nil {
			return err
		}

		*res = combine(aa, bb, cc)
		return nil
	}
}

func Map4[A, B, C, D, T any](
	combine func(A, B, C, D) T,
	da Decoder[A],
	db Decoder[B],
	dc Decoder[C],
	dd Decoder[D],
) Decoder[T] {
	return func(v any, res *T) error {
		var destA A
		if err := da(v, &destA); err != nil {
			return err
		}
		var destB B
		if err := db(v, &destB); err != nil {
			return err
		}
		var destC C
		if err := dc(v, &destC); err != nil {
			return err
		}
		var destD D
		if err := dd(v, &destD); err != nil {
			return err
		}

		*res = combine(destA, destB, destC, destD)
		return nil
	}
}

func Map5[A, B, C, D, E, T any](
	combine func(A, B, C, D, E) T,
	da Decoder[A],
	db Decoder[B],
	dc Decoder[C],
	dd Decoder[D],
	de Decoder[E],
) Decoder[T] {
	return func(v any, res *T) error {
		var destA A
		if err := da(v, &destA); err != nil {
			return err
		}
		var destB B
		if err := db(v, &destB); err != nil {
			return err
		}
		var destC C
		if err := dc(v, &destC); err != nil {
			return err
		}
		var destD D
		if err := dd(v, &destD); err != nil {
			return err
		}
		var destE E
		if err := de(v, &destE); err != nil {
			return err
		}

		*res = combine(destA, destB, destC, destD, destE)
		return nil
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
func Required[A any](name string, da Decoder[A]) Decoder[A] {
	return func(v any, res *A) error {
		vm, ok := v.(map[string]any)
		if !ok {
			return fmt.Errorf("not a dict")
		}
		v, ok = vm[name]
		if !ok {
			return fmt.Errorf("key %q not found", name)
		}

		if err := da(v, res); err != nil {
			return err
		}
		return nil
	}
}

func At[T any](names []string, decoder Decoder[T]) Decoder[T] {
	res := decoder
	for i := len(names) - 1; i >= 0; i-- {
		res = Field(names[i], res)
	}
	return res
}

func Index[T any](i int, decoder Decoder[T]) Decoder[T] {
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

func Optional[A any](
	name string,
	da Decoder[A],
	fallback A,
) Decoder[A] {
	return func(v any, res *A) error {
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
