// Package fun provides reusable general-purpose functions (Const, Swap, Curry) and
// data structures (Unit, Pair, Either).
package fun

import (
	"fmt"
	"log"
)

// ToString converts the value to string.
func ToString[A any](a A) string {
	return fmt.Sprintf("%v", a)
}

// DebugP returns function that prints prefix with element and returns it.
// Useful for debug printing.
func DebugP[V any](prefix string) func(V) V {
	return func(v V) V {
		// TODO: print place
		log.Println(prefix, v)
		return v
	}
}

// Debug returns function that prints element and returns it.
// Useful for debug printing.
func Debug[V any](v V) V {
	// TODO: print place
	log.Println(v)
	return v
}
