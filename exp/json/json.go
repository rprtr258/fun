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

func Field[T any](name string, decoder Decoder[T]) Decoder[T] {
	return func(b []byte) (T, error) {
		var x map[string]any
		if err := json.Unmarshal(b, &x); err != nil {
			return *new(T), err
		}

		value, ok := x[name]
		if !ok {
			return *new(T), fmt.Errorf("field %q not found", name)
		}

		// TODO: vahui
		bb, err := json.Marshal(value)
		if err != nil {
			return *new(T), err
		}
		return decoder(bb)
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

func succeed[T any](x T) Decoder[T] {
	return func([]byte) (T, error) {
		return x, nil
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

func requiredAt[A, B any](names []string, da Decoder[A], df Decoder[func(A) B]) Decoder[B] {
	if len(names) == 0 {
		panic("names must not be empty")
	}
	if len(names) == 1 {
		return required(names[0], da, df)
	}
	return func(b []byte) (B, error) {
		var m map[string]any
		if err := json.Unmarshal(b, &m); err != nil {
			return *new(B), err
		}

		v := m[names[0]]
		for _, name := range names[1:] {
			v = v.(map[string]any)[name]
		}

		// vahui
		bb, err := json.Marshal(v)
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
