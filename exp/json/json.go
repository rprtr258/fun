// See for inspiration
// https://package.elm-lang.org/packages/elm/json/latest/Json-Decode
// https://package.elm-lang.org/packages/NoRedInk/elm-json-decode-pipeline/latest/Json-Decode-Pipeline
package json

import (
	"encoding/json"
	"fmt"
)

type Decoder[T any] func([]byte) (T, error)

func decoder[T any](b []byte) (T, error) {
	var x T
	err := json.Unmarshal(b, &x)
	return x, err
}

var Int Decoder[int] = decoder[int]
var String Decoder[string] = decoder[string]
var Bool Decoder[bool] = decoder[bool]
var Float Decoder[float64] = decoder[float64]

type Maybe[T any] struct {
	Value T
	Valid bool
}

func nullable[T any](decoder Decoder[T]) Decoder[Maybe[T]] {
	return func(b []byte) (Maybe[T], error) {
		if string(b) == "null" {
			return Maybe[T]{}, nil
		}

		value, err := decoder(b)
		if err != nil {
			return Maybe[T]{}, nil
		}

		return Maybe[T]{value, true}, nil
	}
}

func dict[T any](decoder Decoder[T]) Decoder[map[string]T] {
	return func(b []byte) (map[string]T, error) {
		var m map[string]any
		if err := json.Unmarshal(b, &m); err != nil {
			return nil, err
		}

		res := make(map[string]T, len(m))
		for k, v := range m {
			// vahui
			bb, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}

			vv, err := decoder(bb)
			if err != nil {
				return nil, err
			}

			res[k] = vv
		}
		return res, nil
	}
}

func List[T any](decoder Decoder[T]) Decoder[[]T] {
	return func(b []byte) ([]T, error) {
		var x []any
		if err := json.Unmarshal(b, &x); err != nil {
			return nil, err
		}

		// TODO: vahui
		res := make([]T, len(x))
		for i, value := range x {
			bb, err := json.Marshal(value)
			if err != nil {
				return nil, err
			}
			res[i], err = decoder(bb)
			if err != nil {
				return nil, err
			}
		}
		return res, nil
	}
}

// TODO: expressible as required
func Field[T any](name string, decoder Decoder[T]) Decoder[T] {
	return required(name, decoder, succeed(func(t T) T { return t }))
}

func oneOf[T any](decoders []Decoder[T]) Decoder[T] {
	return func(b []byte) (T, error) {
		errors := make([]error, 0, len(decoders))
		for _, decoder := range decoders {
			if res, err := decoder(b); err == nil {
				return res, nil
			} else {
				errors = append(errors, err)
			}
		}
		return *new(T), fmt.Errorf("all variants failed: %v", errors)
	}
}

func Map[T, R any](decoder Decoder[T], f func(T) R) Decoder[R] {
	return func(b []byte) (R, error) {
		value, err := decoder(b)
		if err != nil {
			return *new(R), nil
		}

		return f(value), nil
	}
}

func map2[A, B, T any](
	combine func(A, B) T,
	da Decoder[A],
	db Decoder[B],
) Decoder[T] {
	return func(b []byte) (T, error) {
		aa, err := da(b)
		if err != nil {
			return *new(T), err
		}
		bb, err := db(b)
		if err != nil {
			return *new(T), err
		}

		return combine(aa, bb), nil
	}
}

func map3[A, B, C, T any](
	combine func(A, B, C) T,
	da Decoder[A],
	db Decoder[B],
	dc Decoder[C],
) Decoder[T] {
	return func(b []byte) (T, error) {
		aa, err := da(b)
		if err != nil {
			return *new(T), err
		}
		bb, err := db(b)
		if err != nil {
			return *new(T), err
		}
		cc, err := dc(b)
		if err != nil {
			return *new(T), err
		}

		return combine(aa, bb, cc), nil
	}
}

func andThen[A, B any](da Decoder[A], f func(A) Decoder[B]) Decoder[B] {
	return func(b []byte) (B, error) {
		a, err := da(b)
		if err != nil {
			return *new(B), err
		}
		return f(a)(b)
	}
}

func exampleAndThen() {
	type Info struct{}

	var infoDecoderV4 Decoder[Info]
	var infoDecoderV3 Decoder[Info]

	infoHelp := func(version int) Decoder[Info] {
		switch version {
		case 3:
			return infoDecoderV3
		case 4:
			return infoDecoderV4
		default:
			return fail[Info](fmt.Sprintf("Trying to decode info, but version %d is not supported.", version))
		}
	}

	info := andThen(Field("version", Int), infoHelp)
	_ = info
}

func succeed[T any](x T) Decoder[T] {
	return func([]byte) (T, error) {
		return x, nil
	}
}

func null[T any](value T) Decoder[T] {
	return func(b []byte) (T, error) {
		if string(b) == "null" {
			return value, nil
		}
		return *new(T), fmt.Errorf("not null")
	}
}

func fail[T any](msg string) Decoder[T] {
	return func([]byte) (T, error) {
		return *new(T), fmt.Errorf("%s", msg)
	}
}

// Decode a required field.
func required[A, B any](name string, da Decoder[A], df Decoder[func(A) B]) Decoder[B] {
	return func(b []byte) (B, error) {
		a, err := Field(name, da)(b)
		if err != nil {
			return *new(B), err
		}
		f, err := df(b)
		if err != nil {
			return *new(B), err
		}
		return f(a), nil
	}
}

func at[T any](names []string, decoder Decoder[T]) Decoder[T] {
	res := decoder
	for i := len(names) - 1; i >= 0; i-- {
		res = Field(names[i], res)
	}
	return res
}

func index[T any](i int, decoder Decoder[T]) Decoder[T] {
	return func(b []byte) (T, error) {
		var x []any
		if err := json.Unmarshal(b, &x); err != nil {
			return *new(T), err
		}

		if len(x) <= i {
			return *new(T), fmt.Errorf("no such index %d", i)
		}

		value := x[i]
		// TODO: vahui
		bb, err := json.Marshal(value)
		if err != nil {
			return *new(T), err
		}

		return decoder(bb)
	}
}

func optional[A, B any](name string, da Decoder[A], fallback A, df Decoder[func(A) B]) Decoder[B] {
	return func(b []byte) (B, error) {
		var x map[string]any
		if err := json.Unmarshal(b, &x); err != nil {
			return *new(B), err
		}

		value, ok := x[name]
		if !ok {
			// TODO: vahui
			bb, err := json.Marshal(fallback)
			if err != nil {
				return *new(B), err
			}
			var xx any
			if err := json.Unmarshal(bb, &xx); err != nil {
				return *new(B), err
			}
			value = xx
		}

		// TODO: vahui
		bb, err := json.Marshal(value)
		if err != nil {
			return *new(B), err
		}
		a, err := da(bb)
		if err != nil {
			return *new(B), err
		}
		f, err := df(bb)
		if err != nil {
			return *new(B), err
		}
		return f(a), nil
	}
}

func DecodeString[T any](s string, decoder Decoder[T]) (T, error) {
	return decoder([]byte(s))
}

func example2() {
	type User struct {
		ID    int
		Name  string
		Email string
	}
	newUser := func(id int) func(name string) func(email string) User {
		return func(name string) func(email string) User {
			return func(email string) User {
				return User{id, name, email}
			}
		}
	}

	userDecoder :=
		required("email", String,
			required("name", String,
				required("id", Int,
					succeed(newUser))))

	result, err := userDecoder([]byte(`{"id": 123, "email": "sam@example.com", "name": "Sam"}`))
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

func example() {
	type Job struct {
		name      string
		id        int
		completed bool
	}

	var point Decoder[Job] = map3(
		func(name string, id int, completed bool) Job { return Job{name, id, completed} },
		Field("name", String),
		Field("id", Int),
		Field("completed", Bool),
	)
	_ = point
}
