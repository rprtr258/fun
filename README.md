[![Go Reference](https://pkg.go.dev/badge/github.com/rprtr258/fun.svg)](https://pkg.go.dev/github.com/rprtr258/fun)

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
type Result[T any] struct {Value T; Err error}
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
