// See for inspiration
// https://package.elm-lang.org/packages/elm/json/latest/Json-Decode
// https://package.elm-lang.org/packages/NoRedInk/elm-json-decode-pipeline/latest/Json-Decode-Pipeline
package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type Decoder[T any] func([]byte, *T) error

func primitiveDecoder[T any](b []byte, res *T) error {
	return json.Unmarshal(b, res)
}

var Int Decoder[int] = primitiveDecoder[int]
var String Decoder[string] = primitiveDecoder[string]
var Bool Decoder[bool] = primitiveDecoder[bool]
var Float Decoder[float64] = primitiveDecoder[float64]
var Time Decoder[time.Time] = primitiveDecoder[time.Time]

type Maybe[T any] struct {
	Value T
	Valid bool
}

func nullable[T any](decoder Decoder[T]) Decoder[Maybe[T]] {
	return func(b []byte, res *Maybe[T]) error {
		if string(b) == "null" {
			return nil
		}

		err := decoder(b, &res.Value)
		if err != nil {
			return nil
		}

		res.Valid = true
		return nil
	}
}

func dict[T any](decoder Decoder[T]) Decoder[map[string]T] {
	return func(b []byte, res *map[string]T) error {
		var m map[string]any
		if err := json.Unmarshal(b, &m); err != nil {
			return err
		}

		for k, v := range m {
			// vahui
			bb, err := json.Marshal(v)
			if err != nil {
				return err
			}

			var vv T
			if err := decoder(bb, &vv); err != nil {
				return err
			}

			(*res)[k] = vv
		}
		return nil
	}
}

func List[T any](decoder Decoder[T]) Decoder[[]T] {
	return func(b []byte, res *[]T) error {
		var x []any
		if err := json.Unmarshal(b, &x); err != nil {
			return err
		}

		for _, value := range x {
			// TODO: vahui
			bb, err := json.Marshal(value)
			if err != nil {
				return err
			}
			var t T
			err = decoder(bb, &t)
			if err != nil {
				return err
			}
			*res = append(*res, t)
		}
		return nil
	}
}

// TODO: expressible as required
func Field[T any](name string, decoder Decoder[T]) Decoder[T] {
	return Required(name, decoder, Success(func(t T) T { return t }))
}

func OneOf[T any](decoders []Decoder[T]) Decoder[T] {
	return func(b []byte, res *T) error {
		errors := make([]error, 0, len(decoders))
		for _, decoder := range decoders {
			var t T
			if err := decoder(b, &t); err == nil {
				*res = t
				return nil
			} else {
				errors = append(errors, err)
			}
		}
		return fmt.Errorf("all variants failed: %v", errors)
	}
}

func Map[T, R any](decoder Decoder[T], f func(T) R) Decoder[R] {
	return func(b []byte, res *R) error {
		var t T
		err := decoder(b, &t)
		if err != nil {
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
	return func(b []byte, res *T) error {
		var aa A
		if err := da(b, &aa); err != nil {
			return err
		}

		var bb B
		if err := db(b, &bb); err != nil {
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
	return func(b []byte, res *T) error {
		var aa A
		if err := da(b, &aa); err != nil {
			return err
		}
		var bb B
		if err := db(b, &bb); err != nil {
			return err
		}
		var cc C
		if err := dc(b, &cc); err != nil {
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
	return func(b []byte, res *T) error {
		var destA A
		if err := da(b, &destA); err != nil {
			return err
		}
		var destB B
		if err := db(b, &destB); err != nil {
			return err
		}
		var destC C
		if err := dc(b, &destC); err != nil {
			return err
		}
		var destD D
		if err := dd(b, &destD); err != nil {
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
	return func(b []byte, res *T) error {
		var destA A
		if err := da(b, &destA); err != nil {
			return err
		}
		var destB B
		if err := db(b, &destB); err != nil {
			return err
		}
		var destC C
		if err := dc(b, &destC); err != nil {
			return err
		}
		var destD D
		if err := dd(b, &destD); err != nil {
			return err
		}
		var destE E
		if err := de(b, &destE); err != nil {
			return err
		}

		*res = combine(destA, destB, destC, destD, destE)
		return nil
	}
}

func AndThen[A, B any](da Decoder[A], f func(A) Decoder[B]) Decoder[B] {
	return func(b []byte, res *B) error {
		var a A
		if err := da(b, &a); err != nil {
			return err
		}
		return f(a)(b, res)
	}
}

func Success[T any](x T) Decoder[T] {
	return func(_ []byte, res *T) error {
		*res = x
		return nil
	}
}

func Null[T any](value T) Decoder[T] {
	return func(b []byte, res *T) error {
		if string(b) == "null" {
			*res = value
			return nil
		}
		return fmt.Errorf("not null")
	}
}

func Fail[T any](msg string) Decoder[T] {
	return func([]byte, *T) error {
		return fmt.Errorf("%s", msg)
	}
}

// Decode a Required field.
func Required[A, B any](name string, da Decoder[A], df Decoder[func(A) B]) Decoder[B] {
	return func(b []byte, res *B) error {
		dec := json.NewDecoder(bytes.NewReader(b))
		if tok, err := dec.Token(); err != nil {
			return err
		} else if tok != json.Delim('{') {
			return fmt.Errorf("not a dict")
		}
		for {
			if tok, err := dec.Token(); err != nil {
				return err
			} else if tok == json.Delim('}') {
				return fmt.Errorf("key not found")
			} else if s, ok := tok.(string); ok {
				if s != name {
					dec.Token()
					continue
				}
			}
			break
		}

		var a A
		if err := dec.Decode(&a); err != nil {
			return err
		}
		var f func(A) B
		if err := df(b, &f); err != nil {
			return err
		}
		*res = f(a)
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
	return func(b []byte, res *T) error {
		var x []any
		if err := json.Unmarshal(b, &x); err != nil {
			return err
		}

		if len(x) <= i {
			return fmt.Errorf("no such index %d", i)
		}

		value := x[i]
		// TODO: vahui
		bb, err := json.Marshal(value)
		if err != nil {
			return err
		}

		return decoder(bb, res)
	}
}

func Optional[A any](
	name string,
	da Decoder[A],
	fallback A,
) Decoder[A] {
	return func(b []byte, res *A) error {
		var x map[string]any
		if err := json.Unmarshal(b, &x); err != nil {
			return err
		}

		value, ok := x[name]
		if !ok {
			// TODO: vahui
			bb, err := json.Marshal(fallback)
			if err != nil {
				return err
			}
			var xx any
			if err := json.Unmarshal(bb, &xx); err != nil {
				return err
			}
			value = xx
		}

		// TODO: vahui
		bb, err := json.Marshal(value)
		if err != nil {
			return err
		}
		var a A
		if err := da(bb, &a); err != nil {
			return err
		}
		*res = a
		return nil
	}
}

func DecodeString[T any](s string, decoder Decoder[T]) (T, error) {
	var t T
	err := decoder([]byte(s), &t)
	return t, err
}
