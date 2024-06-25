package iter

import (
	"errors"
	"fmt"
	"io"
	"os"
)

type Seq[T any] func(bool) (T, bool)
type Seq2[K, V any] func(bool) (K, V, bool)

var ErrEnd = errors.New("THE END")

type ESeq[T any] func(bool) (T, error)
type ESeq2[K, V any] func(bool) (K, V, error)

func Map[T, R any](f func(T) R, s Seq[T]) Seq[R] {
	return func(b bool) (R, bool) {
		x, ok := s(b)
		if !ok {
			return *new(R), false
		}
		return f(x), ok
	}
}

func Backward[T any](s []T) Seq2[int, T] {
	i := len(s) - 1
	return func(onBreak bool) (int, T, bool) {
		if onBreak || i < 0 {
			// cleanup
			return 0, *new(T), false
		}
		idx, elem := i, s[i]
		i--
		return idx, elem, true
	}
}

func exampleBackward() {
	s := []int{1, 2, 3}
	for it := Backward(s); ; {
		_, el, ok := it(false)
		if !ok {
			break // it(true) does not need to be called because the `false` was called
		}

		fmt.Print(el, " ")
	}
}

func Reader(r io.ReadCloser) func(bool) ([]byte, error) {
	var err error
	b := make([]byte, 4*1024)
	return func(onBreak bool) ([]byte, error) {
		if onBreak || err != nil {
			_ = r.Close()
			return nil, err
		}
		n, err := r.Read(b)
		if err == io.EOF {
			err = ErrEnd
		}
		return b[:n], err
	}
}

func example() {
	f, _ := os.Open("aboba.txt")
	for it := Reader(f); ; {
		b, err := it(false)
		if err != nil {
			break // it(true) does not need to be called because the `false` was called
		}

		fmt.Print(string(b))
	}
}
