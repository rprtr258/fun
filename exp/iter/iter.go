package iter

import (
	"errors"
	"io"
)

type (
	Seq[T any]     func(bool) (T, bool)
	Seq2[K, V any] func(bool) (K, V, bool)
)

var ErrEnd = errors.New("THE END")

type (
	ESeq[T any]     func(bool) (T, error)
	ESeq2[K, V any] func(bool) (K, V, error)
)

func (s Seq[T]) Map[R any](f func(T) R) Seq[R] {
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
