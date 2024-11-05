// package coro3
package main

import (
	"fmt"
	"iter"
	_ "unsafe"
)

type coro struct{}

//go:linkname newcoro runtime.newcoro
func newcoro(func(*coro)) *coro

//go:linkname coroswitch runtime.coroswitch
func coroswitch(*coro)

type Seq[T any] func() (T, bool)

func (s Seq[T]) Skip() bool {
	_, ok := s()
	return ok
}

func (s Seq[T]) Take2() (T, T, bool) {
	x1, ok1 := s()
	x2, ok2 := s()
	return x1, x2, ok1 && ok2
}

func (s Seq[T]) Take(n int) ([]T, bool) {
	res := make([]T, n)
	for i := 0; i < n; i++ {
		x, ok := s()
		if !ok {
			return res, false
		}
		res[i] = x
	}
	return res, true
}

// TODO: cancel?
func Pull[T any](seq iter.Seq[T]) func() (T, bool) {
	var v T
	var ok bool
	c := newcoro(func(c *coro) {
		for x := range seq {
			v, ok = x, true
			coroswitch(c)
		}
		ok = false
	})
	return func() (T, bool) {
		coroswitch(c)
		if !ok {
			return *new(T), false
		}
		return v, ok
	}
}

func chain[T any](xss ...Seq[T]) Seq[T] {
	return Pull(func(yield func(T) bool) {
		for _, xs := range xss {
			for {
				if x, ok := xs(); !ok || !yield(x) {
					return
				}
			}
		}
	})
}

func parseQuoted(next Seq[byte]) bool {
	if c, ok := next(); !ok || c != '"' {
		return false
	}

	for {
		if c, ok := next(); !ok || c == '"' {
			return c == '"'
		} else if c == '\\' && !next.Skip() {
			return false
		}
	}
}

var base64 = func() [128]uint8 {
	res := [128]uint8{}
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	for i := 0; i < len(alphabet); i++ {
		res[alphabet[i]] = uint8(i)
	}
	return res
}()

func base64Decode(next Seq[byte]) Seq[byte] {
	var acc uint16
	var n = 0
	return func() (byte, bool) {
		for n < 8 {
			c, ok := next()
			if !ok {
				return 0, false
			}

			n += 6
			acc = (acc << 6) | uint16(base64[c])
		}
		mask := uint16(0xFF << (n - 8))
		c := acc & mask
		acc = acc &^ mask
		return byte(c), true
	}
}

type Tree[T any] struct {
	l, r  *Tree[T]
	value T
}

func (t *Tree[T]) All() func() (T, bool) {
	return Pull(func(yield func(T) bool) {
		var dfs func(t *Tree[T]) bool
		dfs = func(t *Tree[T]) bool {
			if t == nil {
				return true
			}

			return dfs(t.l) &&
				yield(t.value) &&
				dfs(t.r)
		}
		dfs(t)
	})
}

func slice[T any](xs []T) Seq[T] {
	i := -1
	return func() (T, bool) {
		if i+1 == len(xs) {
			return *new(T), false
		}
		i++
		return xs[i], true
	}
}

func main() {
	for _, s := range []string{
		// "true" cases
		`"ABOBA"`,
		`"AB\nBA"`,
		`"AB\"BA"`,
		`"AB"BA"`, // NOTE: prefix is string, so true
		// "false" cases
		`ABOBA"`,
		`"ABOBA`,
	} {
		fmt.Println(s, parseQuoted(slice([]byte(s))))
	}

	tree := &Tree[int]{
		l: &Tree[int]{
			l: &Tree[int]{
				value: 1,
			},
			value: 2,
		},
		value: 3,
		r: &Tree[int]{
			value: 4,
			r: &Tree[int]{
				value: 5,
			},
		},
	}

	fmt.Print("DFS: ")
	next := tree.All()
	for {
		t, ok := next()
		if !ok {
			break
		}
		if t > 3 {
			break
		}
		if t%2 == 0 {
			continue
		}
		fmt.Print(t, " ")
	}
	fmt.Println()
}
