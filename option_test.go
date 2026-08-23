package fun_test

import (
	"encoding/json"
	"fmt"
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

func ExampleOptMap() {
	fmt.Println(fun.OptMap(fun.Valid(1), func(x int) string {
		return fmt.Sprintf("%d", x)
	}))
	// Output: Some(1)
}

func ExampleOptFlatMap() {
	fmt.Println(fun.OptFlatMap(fun.Valid(1), func(x int) fun.Option[string] {
		return fun.Valid(fmt.Sprintf("%d", x))
	}))
	// Output: Some(1)
}
