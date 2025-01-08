# Iterator and functional utilities

The design is inspired by [samber/lo](https://github.com/samber/lo) and [iterator proposal](https://github.com/golang/go/issues/61897). This library does not deal with channel/pipes/concurrency as that is beyond the scope of this project.

## Root package
Root package `github.com/rprtr258/fun` provides common slice and functional utilities.

### Core types

```go
// Pair is a data structure that has two values.
type Pair[K, V any] struct {K K; V V}

// Option is either value or nothing.
type Option[T any] struct {Value T; Valid bool}

// Result is either value or error.
type Result[T any] Pair[T, error]
```

### Core constraints
```go
// RealNumber is a generic number interface that covers all Go real number types.
type RealNumber interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}

// Number is a generic number interface that covers all Go number types.
type Number interface {
	RealNumber | complex64 | complex128
}
```

### Design decisions
Declarations like
```go
func Map[R, T any, F interface {
	func(T) R | func(T, int) R
}](f F, slice ...T) []R
```
exists for the reason that we want both `func(elem)` and `func(elem, index)` functions work. With such declaration Go cannot infer type `R`, so we have to specify it explicitly on usage: `fun.Map[string](fn, slice...)`

Another moment is that slice arguments are variadic. That allows user not to construct slice in some cases like `fun.Contains(status, "OK", "Success")` instead of `fun.Contains(status, []string{"OK", "Success"})`.

### slice functions
#### Map
Applies function to all elements and returns slice with results.

```go
fun.Map(func(x int64, _ int) string {
	return strconv.FormatInt(x, 10)
}, 0, 1, 2)
// []string{"0", "1", "2"}
```

#### Filter
Filters slice elements using given predicate.

```go
fun.Filter(func(x int64, _ int) bool {
	return x%2 == 0
}, 0, 1, 2)
// []int64{0, 2}
```

#### FilterMap
Transform each element, leaving only those for which true is returned.

```go
fun.FilterMap(func(x int64, _ int) (string, bool) {
	return strconv.FormatInt(x, 10), x%2 == 0
}, 0, 1, 2)
// []string{"0", "2"}
```

#### MapDict
Like `Map` but uses dictionary instead of transform function.

```go
dict := map[int]string{
	0: "zero",
	1: "one",
	2: "two",
}
fun.MapDict(dict, 0, 1, 2)
// []string{"zero", "one", "two"}
```

#### MapErr
Like `Map` but returns first error got from transform.

```go
fun.MapErr(func(x int64, _ int) (string, error) {
	if x%2 == 0 {
		return strconv.FormatInt(x, 10), nil
	}
	return "", errors.New("odd")
}, 0, 1, 2)
// []string{"0"}, errors.New("odd")
```

#### MapToSlice
Transforms map to slice using transform function. Order is not guaranteed.

```go
dict := map[int]string{
	0: "zero",
	1: "one",
	2: "two",
}
fun.MapToSlice(dict, func(k int, v string) string {
	return fmt.Sprintf("%d: %s", k, v)
})
// []string{"0: zero", "1: one", "2: two"}
```

#### MapFilterToSlice
Transforms map to slice using transform function and returns only those for which true is returned. Order is not guaranteed.

```go
dict := map[int]string{
	0: "zero",
	1: "one",
	2: "two",
}
fun.MapFilterToSlice(dict, func(k int, v string) (string, bool) {
	return fmt.Sprintf("%d: %s", k, v), k%2 == 0
})
// []string{"0: zero", "2: two"}
```

#### Keys
Returns keys of map. Order is not guaranteed.

```go
dict := map[int]string{
	0: "zero",
	1: "one",
	2: "two",
}
fun.Keys(dict)
// []int{0, 1, 2}
```

#### Values
Returns values of map. Order is not guaranteed.

```go
dict := map[int]string{
	0: "zero",
	1: "one",
	2: "two",
}
fun.Values(dict)
// []string{"zero", "one", "two"}
```

#### FindKeyBy
Returns the key of the first element predicate returns truthy for.

```go
dict := map[int]string{
	0: "zero",
	1: "one",
	2: "two",
}
fun.FindKeyBy(dict, func(k int, v string) bool {
	return v == "zero"
})
// 0, true
```

#### Uniq
Returns unique values of slice. In other words, removes duplicates.

```go
fun.Uniq(1, 2, 3, 1, 2)
// []int{1, 2, 3}
```

#### Index
Returns first found element by predicate along with it's index.

```go
fun.Index(func(s string, _ int) bool {
	return strings.HasPrefix(s, "o")
}, "zero", "one", "two")
// "one", 1, true
```

#### Contains
Returns true if an element is present in a collection.

```go
fun.Contains("zero", "zero", "one", "two")
// true
```

#### SliceToMap
Returns a map containing key-value pairs provided by transform function applied to elements of the given slice.

```go
fun.SliceToMap(func(x int, _ int) (int, int) {
	return x, x * 10
}, 0, 1, 2)
// map[int]int{0: 0, 1: 10, 2: 20}
```

#### FromMap
Returns slice of key/value pairs from map.

```go
dict := map[int]string{
	0: "zero",
	1: "one",
	2: "two",
}
fun.FromMap(dict)
// []fun.Pair[int, string]{0: "zero", 1: "one", 2: "two"}
```

#### Copy
Returns copy of slice.

```go
fun.Copy(1, 2, 3)
// []int{1, 2, 3}
```

#### ReverseInplace
Reverses slice in place.

```go
xs := []int{1, 2, 3}
fun.ReverseInplace(xs)
// xs becomes []int{3, 2, 1}
```

#### Subslice
Returns slice from start to end without panicking on out of bounds.

```go
xs := []int{1, 2, 3, 4, 5}
fun.Subslice(1, 4, xs...)
// []int{2, 3, 4}
```

#### Chunk
Divides slice into chunks of size chunkSize.

```go
xs := []int{1, 2, 3, 4, 5}
fun.Chunk(2, xs...)
// [][]int{{1, 2}, {3, 4}, {5}}
```

#### ConcatMap
Like `Map` but concatenates results.

```go
fun.ConcatMap(func(x int) []int {
	return []int{x, x + 10, x + 100}
}, 0, 1, 2)
// []int{0, 10, 100, 1, 11, 101, 2, 12, 102}
```

#### All
Returns true if all elements satisfy the condition.

```go
fun.All(func(x int) bool {
	return x%2 == 0
}, 0, 2, 4)
// true
```

#### Any
Returns true if any (at least one) element satisfies the condition.

```go
fun.Any(func(x int) bool {
	return x%2 == 0
}, 0, 1, 2)
// true
```

#### SortBy
Sorts slice in place by given function.

```go
xs := []int{1, 2, 3, 4, 5}
fun.SortBy(func(x int) int {
	return -x
}, xs)
// xs becomes []int{5, 4, 3, 2, 1}
```

#### GroupBy
Groups elements by key.

```go
fun.GroupBy(func(x int) int {
	return x % 2
}, 0, 1, 2, 3, 4)
// map[int][]int{0: {0, 2, 4}, 1: {1, 3}}
```

### cmp
Utilities utilizing values comparison.

#### Min
Returns the minimum of the given values.

```go
fun.Min(1, 2, 3)
// 1
```

#### Max
Returns the maximum of the given values.

```go
fun.Max(1, 2, 3)
// 3
```

#### Clamp(x, low, high)
Returns x clamped between low and high.

```go
fun.Clamp(99, 1, 10)
// 10
```

#### MinBy
Returns first minimum of given values using given order function.

```go
fun.MinBy(func(s string) int {
	return len(s)
}, "one", "two", "three")
// "one"
```

#### MaxBy
Returns first maximum of given values using given order function.

```go
fun.MaxBy(func(s string) int {
	return len(s)
}, "one", "two", "three")
// "three"
```

### Working with Option type

#### Invalid
Returns empty Option.

```go
fun.Invalid[int]()
// Option[int]{}
```

#### Valid
Returns Option with given value.

```go
fun.Valid(1)
// Option[int]{Value: 1, Valid: true}
```

#### Optional
Returns Option with given value and validity.

```go
fun.Optional(1, true)
// Option[int]{Value: 1, Valid: true}
```

#### FromPtr
Returns Option with value from pointer.

```go
x := 1
fun.FromPtr(&x)
// Option[int]{Value: 1, Valid: true}
fun.FromPtr[int](nil)
// Option[int]{}
```

#### Option.Unpack
Returns value and validity.

```go
fun.Valid(1).Unpack()
// (1, true)
```

#### Option.Or
Returns first valid Option.

```go
fun.Valid(1).Or(fun.Invalid[int]())
// Option[int]{Value: 1, Valid: true}
```

#### Option.OrDefault
Returns value if Option is valid, otherwise returns default value.

```go
fun.Valid(1).OrDefault(0)
// 1
```

#### Option.Ptr
Returns pointer to value if Option is valid, otherwise returns nil.

```go
fun.Valid(1).Ptr()
// &[]int{1}[0]
```

#### OptMap
Returns new Option with transformed value.

```go
fun.Valid(1).OptMap(func(x int) string {
	return fmt.Sprintf("%d", x)
})
// Option[string]{Value: "1", Valid: true}
```

#### OptFlatMap
Returns new Option with transformed optional value.

```go
fun.Valid(1).OptFlatMap(func(x int) Option[string] {
	return fun.Valid(fmt.Sprintf("%d", x))
})
// Option[string]{Value: "1", Valid: true}
```

### fp

#### Zero
Returns zero value of given type.

```go
fun.Zero[int]()
// 0
```

#### Debug
Prints value and returns it. Useful for debug printing.

```go
fun.Debug(2+2)*2
// prints 4
```

#### Has
Returns true if map has such key.

```go
dict := map[int]string{
	0: "zero",
	1: "one",
	2: "two",
}
fun.Has(dict, 2)
// true
```

#### Cond
Returns first value for which true is returned.

```go
fun.Cond(
	1,
	func() (int, bool) { return 2, true },
	func() (int, bool) { return 3, false },
)
// 2
```

#### Ptr
Returns pointer to value.

```go
fun.Ptr(1)
// &[]int{1}[0]
```

#### Deref
Returns value from pointer. If pointer is nil returns zero value.

```go
fun.Deref[int](nil) // 0
fun.Deref[int](new(int)) // 0
x := 1
fun.Deref[int](&x) // 1
```

#### Pipe
Returns value after applying endomorphisms. Endomorphism is just function from type to itself.

```go
fun.Pipe(
	"hello  ",
	strings.TrimSpace,
	strings.NewReplacer("l", "|").Replace,
	strings.ToUpper,
)
// "HE||O"
```

#### If
There are multiple variations of `if` statement usable as expression.

#### IF
Simple ternary function.

```go
fun.IF(true, 1, 0)
// 1
```

#### If, IfF
Returns value from branch for which predicate is true. `F` suffix can be used to get values not evaluated immediately.

```go
fun.If(true, 1).Else(0)
// 1
fun.If(false, 1).ElseF(func() int { return 0 })
// 0
fun.If(false, 1).ElseIf(true, 2).Else(3)
// 2
fun.IfF(false, func() int { return 1 }).Else(0)
// 0

fun.If(true, db.Get(0)).Else(db.Get(1))
// db.Get(0) result, db.Get called two times
fun.IfF(true, func() Thing { return db.Get(0) }).ElseF(func() Thing { return db.Get(1) })
// db.Get(0) result, db.Get called once
```

### Switch
`switch` usable as expression.

```go
fun.Switch("one", -1).
	Case("zero", 0).
	Case("one", 1).
	Case("two", 2).
	End()
// 1
```

## Iter

`github.com/rprtr258/fun/iter` introduces iterator primitives for which `iter.Seq[T]` is basic.

```go
type Seq[V any] func(yield func(V) bool)
```

Which is a function which accepts function to `yield` values from iteration. `yield` must return `false` when iteration must stop (analogous to `break`).

Example iterator yielding numbers from 1 to `n`, including `n`:

```go
func Range(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := range n {
			if !yield(i) {
				return
			}
		}
	}
}
```

## Set

`github.com/rprtr258/fun/set` introduces `Set[T]` primitive for collections of unique `comparable` values.

## Ordered map

`github.com/rprtr258/fun/orderedmap` introduces `OrderedMap[K, V]` data structure which acts like hashmap but also allows to iterate over keys in sorted order. Internally, binary search tree is used.
