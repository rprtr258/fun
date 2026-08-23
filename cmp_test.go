package fun_test

import (
	"fmt"

	"github.com/rprtr258/fun"
)

func ExampleMin() {
	fmt.Println(fun.Min(1, 2, 3))
	// 1
}

func ExampleMax() {
	fmt.Println(fun.Max(1, 2, 3))
	// 3
}

func ExampleClamp() {
	fmt.Println(fun.Clamp(99, 1, 10))
	// 10
}

func ExampleMinBy() {
	fmt.Println(fun.MinBy(func(s string) int {
		return len(s)
	}, "one", "two", "three"))
	// "one"
}

func ExampleMaxBy() {
	fmt.Println(fun.MaxBy(func(s string) int {
		return len(s)
	}, "one", "two", "three"))
	// "three"
}
