package iter

import (
	"fmt"
	"os"
)

func ExampleBackward() {
	s := []int{1, 2, 3}
	for it := Backward(s); ; {
		_, el, ok := it(false)
		if !ok {
			break // it(true) does not need to be called because the `false` was called
		}

		fmt.Print(el, " ")
	}
}

func Example() {
	f, _ := os.Open("aboba.txt")
	for it := Reader(f); ; {
		b, err := it(false)
		if err != nil {
			break // it(true) does not need to be called because the `false` was called
		}

		fmt.Print(string(b))
	}
}
