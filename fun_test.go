package fun_test

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/rprtr258/assert"
	"github.com/rprtr258/fun"
)

func TestToString(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "1", fun.ToString(1))
}

func ExampleEntries() {
	dict := map[int]string{
		0: "zero",
		1: "one",
		2: "two",
	}
	entries := fun.Entries(dict)
	fun.SortBy(func(kv fun.Pair[int, string]) int { return kv.K }, entries...)
	fmt.Println(entries)
	// Output: [{0 zero} {1 one} {2 two}]
}

func ExampleSliceToMap() {
	fmt.Println(fun.SliceToMap[int, int](func(x int, _ int) (int, int) {
		return x, x * 10
	}, 0, 1, 2))
	// Output: map[0:0 1:10 2:20]
}

func ExampleMap() {
	fmt.Println(fun.Map[string, int64](func(x int64, _ int) string {
		return strconv.FormatInt(x, 10)
	}, 0, 1, 2))
	// Output: [0 1 2]
}

func ExampleFilter() {
	fmt.Println(fun.Filter[int64](func(x int64, _ int) bool {
		return x%2 == 0
	}, 0, 1, 2))
	// Output: [0 2]
}

func ExampleFilterMap() {
	fmt.Println(fun.FilterMap[string, int64](func(x int64, _ int) (string, bool) {
		return strconv.FormatInt(x, 10), x%2 == 0
	}, 0, 1, 2))
	// Output: [0 2]
}

func ExampleMapDict() {
	dict := map[int]string{
		0: "zero",
		1: "one",
		2: "two",
	}
	fmt.Println(fun.MapDict(dict, 0, 1, 2))
	// Output: [zero one two]
}

func ExampleMapErr() {
	fmt.Println(fun.MapErr[string, int64](func(x int64, _ int) (string, error) {
		if x%2 == 0 {
			return strconv.FormatInt(x, 10), nil
		}
		return "", errors.New("odd")
	}, 0, 1, 2))
	// Output: [] odd
}

func ExampleFindKeyBy() {
	dict := map[int]string{
		0: "zero",
		1: "one",
		2: "two",
	}
	fmt.Println(fun.FindKeyBy(dict, func(k int, v string) bool {
		return v == "zero"
	}))
	// Output: 0 true
}

func ExampleUniq() {
	fmt.Println(fun.Uniq(1, 2, 3, 1, 2))
	// Output: [1 2 3]
}

func ExampleIndex() {
	fmt.Println(fun.Index(func(s string, _ int) bool {
		return strings.HasPrefix(s, "o")
	}, "zero", "one", "two"))
	// Output: one 1 true
}

func ExampleContains() {
	value := "zero"
	fmt.Println(fun.Contains(value, "zero", "one", "two"))
	// Output: true
}

func ExampleReverseInplace() {
	xs := []int{1, 2, 3}
	fun.ReverseInplace(xs)
	fmt.Println(xs)
	// Output: [3 2 1]
}

func ExampleSubslice() {
	xs := []int{1, 2, 3, 4, 5}
	fmt.Println(fun.Subslice(1, 4, xs...))
	// Output: [2 3 4]
}

func ExampleChunk() {
	xs := []int{1, 2, 3, 4, 5}
	fmt.Println(fun.Chunk(2, xs...))
	// Output: [[1 2] [3 4] [5]]
}

func ExampleConcatMap() {
	fmt.Println(fun.ConcatMap(func(x int) []int {
		return []int{x, x + 10, x + 100}
	}, 0, 1, 2))
	// Output: [0 10 100 1 11 101 2 12 102]
}

func ExampleAll() {
	fmt.Println(fun.All(func(x int) bool {
		return x%2 == 0
	}, 0, 2, 4))
	// Output: true
}

func ExampleAny() {
	fmt.Println(fun.Any(func(x int) bool {
		return x%2 == 0
	}, 0, 1, 2))
	// Output: true
}

func ExampleSortBy() {
	xs := []int{1, 2, 3, 4, 5}
	fun.SortBy(func(x int) int {
		return -x
	}, xs...)
	fmt.Println(xs)
	// Output: [5 4 3 2 1]
}

func ExampleGroupBy() {
	fmt.Println(fun.GroupBy(func(x int) int {
		return x % 2
	}, 0, 1, 2, 3, 4))
	// Output: map[0:[0 2 4] 1:[1 3]]
}

func ExampleZero() {
	fmt.Println(fun.Zero[int]())
	// Output: 0
}

func ExampleDebug() {
	fmt.Println(fun.Debug(2+2) * 2)
	// Output: 8
}

func ExampleHas() {
	dict := map[int]string{
		0: "zero",
		1: "one",
		2: "two",
	}
	fmt.Println(fun.Has(dict, 2))
	// Output: true
}

func ExampleCond() {
	fmt.Println(fun.Cond(
		1,
		func() (int, bool) { return 2, true },
		func() (int, bool) { return 3, false },
	))
	// 2
}

func ExamplePtr() {
	fmt.Println(*fun.Ptr(1))
	// Output: 1
}

func ExampleDeref() {
	fmt.Println(fun.Deref[int](nil))
	fmt.Println(fun.Deref(new(int)))
	x := 1
	fmt.Println(fun.Deref(&x))
	// Output:
	// 0
	// 0
	// 1
}

func ExamplePipe() {
	fmt.Println(fun.Pipe(
		"hello  ",
		strings.TrimSpace,
		strings.NewReplacer("l", "|").Replace,
		strings.ToUpper,
	))
	// Output: HE||O
}

func ExampleIF() {
	fmt.Println(fun.IF(true, 1, 0))
	// Output: 1
}

func ExampleIf() {
	fun.If(true, 1).Else(0)
	fun.If(false, 1).ElseF(func() int { return 0 })
	fun.If(false, 1).ElseIf(true, 2).Else(3)
	// fun.If(true, db.Get(0)).Else(db.Get(1))
	// Output:
	// 1
	// 0
	// 2
	// db.Get(0) result, db.Get called two times
}

func ExampleIfF() {
	fun.IfF(false, func() int { return 1 }).Else(0)

	// fun.IfF(true, func() Thing { return db.Get(0) }).ElseF(func() Thing { return db.Get(1) })
	// Output:
	// 0
	// db.Get(0) result, db.Get called once
}

func ExampleSwitch() {
	fmt.Println(fun.Switch("one", -1).
		Case(0, "zero").
		Case(1, "one").
		Case(2, "two").
		End())
	// Output: 1
}
