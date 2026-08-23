package fun_test

import (
	"fmt"
	"testing"

	"github.com/rprtr258/assert"
	"github.com/rprtr258/fun"
)

func TestMap_noIndex(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		slice    []int
		f        func(int) int
		expected []int
	}{
		"example": {
			slice: []int{1, 2, 3},
			f: func(x int) int {
				return x + 1
			},
			expected: []int{2, 3, 4},
		},
		"empty slice": {
			slice: []int{},
			f: func(x int) int {
				return x + 1
			},
			expected: []int{},
		},
		"nil": {
			slice: nil,
			f: func(x int) int {
				return x + 1
			},
			expected: nil,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.expected, fun.Map[int](test.f, test.slice...))
		})
	}
}

func ExampleMapToSlice() {
	dict := map[int]string{
		0: "zero",
		1: "one",
		2: "two",
	}
	slice := fun.MapToSlice(dict, func(k int, v string) string {
		return fmt.Sprintf("%d=%s", k, v)
	})
	fun.SortBy(func(s string) string { return s }, slice...)
	fmt.Println(slice)
	// Output: [0=zero 1=one 2=two]
}

func ExampleMapFilterToSlice() {
	dict := map[int]string{
		0: "zero",
		1: "one",
		2: "two",
	}
	slice := fun.MapFilterToSlice(dict, func(k int, v string) (string, bool) {
		return fmt.Sprintf("%d=%s", k, v), k%2 == 0
	})
	fun.SortBy(func(s string) string { return s }, slice...)
	fmt.Println(slice)
	// Output: [0=zero 2=two]
}

func ExampleKeys() {
	dict := map[int]string{
		0: "zero",
		1: "one",
		2: "two",
	}
	keys := fun.Keys(dict)
	fun.SortBy(func(k int) int { return k }, keys...)
	fmt.Println(keys)
	// Output: [0 1 2]
}

func ExampleValues() {
	dict := map[int]string{
		0: "zero",
		1: "one",
		2: "two",
	}
	values := fun.Values(dict)
	fun.SortBy(func(s string) string { return s }, values...)
	fmt.Println(values)
	// Output: [one two zero]
}
