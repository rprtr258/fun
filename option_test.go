package fun_test

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/rprtr258/assert"

	"github.com/rprtr258/fun"
)

func TestOptionJSONMarshal(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		opt  fun.Option[int]
		want string
	}{
		"valid": {
			opt:  fun.Valid(1),
			want: "1",
		},
		"invalid": {
			opt:  fun.Invalid[int](),
			want: "null",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := json.Marshal(test.opt)
			assert.NoError(t, err)
			assert.Equal(t, test.want, string(got))
		})
	}
}

func TestOptionJSONUnmarshal(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		opt  string
		want fun.Option[int]
	}{
		"valid": {
			opt:  "1",
			want: fun.Valid(1),
		},
		"invalid": {
			opt:  "null",
			want: fun.Invalid[int](),
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var opt fun.Option[int]
			err := json.Unmarshal([]byte(test.opt), &opt)
			assert.NoError(t, err)
			assert.Equal(t, test.want, opt)
		})
	}
}

func ExampleInvalid() {
	fmt.Println(fun.Invalid[int]())
	// Output: None
}

func ExampleValid() {
	fmt.Println(fun.Valid(1))
	// Output: Some(1)
}

func ExampleOptional() {
	fmt.Println(fun.Optional(1, true))
	// Output: Some(1)
}

func ExampleFromPtr_valid() {
	x := 1
	fmt.Println(fun.FromPtr(&x))
	// Output: Some(1)
}

func ExampleFromPtr_invalid() {
	fmt.Println(fun.FromPtr[int](nil))
	// Output: None
}

func ExampleOption_Unpack() {
	fmt.Println(fun.Valid(1).Unpack())
	// Output: 1 true
}

func ExampleOption_Or() {
	fmt.Println(fun.Valid(1).Or(fun.Invalid[int]()))
	// Output: Some(1)
}

func ExampleOption_OrDefault() {
	fmt.Println(fun.Valid(1).OrDefault(0))
	// Output: 1
}

func ExampleOption_Ptr() {
	fmt.Println(*fun.Valid(1).Ptr())
	// Output: 1
}

func ExampleOption_Map() {
	v := fun.Valid(67)
	fmt.Println(v.Map(strconv.Itoa).OrDefault("wtf"))
	// Output: 67
}

func ExampleOption_FlatMap() {
	fmt.Println(fun.
		Valid(67).
		FlatMap(func(i int) fun.Option[string] {
			if i == 0 {
				return fun.Invalid[string]()
			}
			return fun.Valid(strconv.Itoa(i))
		}).
		OrDefault("wtf"))
	// Output: 67
}
