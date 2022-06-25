// Package fun provides reusable general-purpose functions (Const, Swap, Curry) and
// data structures (Unit, Pair, Either).
package fun

import (
	"fmt"
	"log"
)

// Const creates a function that will ignore it's input and return the specified value.
func Const[A, B any](a A) func(B) A {
	return func(B) A {
		return a
	}
}

// ConstUnit creates a function that will ignore it's Unit input and return the specified value.
func ConstUnit[B any](b B) func(Unit) B {
	return Const[B, Unit](b)
}

// Swap returns a curried function with swapped order of arguments.
func Swap[A, B, C any](f func(a A) func(b B) C) func(b B) func(a A) C {
	return func(b B) func(a A) C {
		return func(a A) C {
			return f(a)(b)
		}
	}
}

// Curry takes a function that has two arguments and returns a function with two argument lists.
func Curry[A, B, C any](f func(a A, b B) C) func(a A) func(b B) C {
	return func(a A) func(b B) C {
		return func(b B) C {
			return f(a, b)
		}
	}
}

// Identity function returns the given value unchanged.
func Identity[A any](a A) A {
	return a
}

// ToString converts the value to string.
func ToString[A any](a A) string {
	return fmt.Sprintf("%v", a)
}

// Compose executes the given functions in sequence.
func Compose[A, B, C any](f func(A) B, g func(B) C) func(A) C {
	return func(a A) C {
		return g(f(a))
	}
}

// Debug returns function that prints element and returns it. Useful for debug printing.
func Debug[A any](prefix string) func(A) A {
	return func(a A) A {
		log.Println(prefix, a)
		return a
	}
}
