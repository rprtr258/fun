// Package text provides some utilities to work with text files.
package text

import (
	"bytes"
	"io"

	"github.com/rprtr258/fun"
	"github.com/rprtr258/fun/iter"
)

const defaultChunkSize = 4 * 1024 // 4 KB

// ReadByteChunks read using buffer of chunkSize size
func ReadByteChunks(r io.Reader, chunkSize int) iter.Seq[fun.Result[[]byte]] {
	return func(yield func(r fun.Result[[]byte]) bool) bool {
		b := make([]byte, chunkSize)
		for {
			n, err := r.Read(b)
			if !yield(fun.Result[[]byte]{append([]byte(nil), b[:n]...), err}) {
				return false
			}
			if err != nil {
				if err == io.EOF {
					return true
				}
			}
		}
	}
}

// SplitBySeparator splits byte-chunk stream by the given separator.
func SplitBySeparator(seq iter.Seq[[]byte], sep byte) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) bool {
		var curBuf []byte
		return seq(func(chunk []byte) bool {
			curBuf = append(curBuf, chunk...)
			for {
				idx := bytes.IndexByte(curBuf, sep)
				if idx == -1 {
					break
				}
				if !yield(curBuf[:idx]) {
					return false
				}
				curBuf = curBuf[idx+1:]
			}
			return true
		}) && yield(curBuf)
	}
}

// ReadLines reads text file line-by-line.
func ReadLines(reader io.Reader) iter.Seq[string] {
	chunks := ReadByteChunks(reader, defaultChunkSize)

	pull, stop := iter.Pull(chunks)
	defer stop()

	rows := SplitBySeparator(func(yield func([]byte) bool) bool {
		for r, ok := pull(); ok; r, ok = pull() {
			if r.V != nil || !yield(r.K) {
				return false
			}
		}
		return true
	}, '\n')

	return iter.Map(rows, func(x []byte) string { return string(x) })
}
