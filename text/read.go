// Package text provides some utilities to work with text files.
package text

import (
	"io"
	"log"

	"github.com/rprtr258/goflow/fun"
	s "github.com/rprtr258/goflow/stream"
)

const defaultChunkSize = 4 * 1024 // 4 KB

type readByteChunksImpl struct {
	reader    io.Reader
	eof       bool
	chunkSize int
}

func (xs *readByteChunksImpl) Next() fun.Option[[]byte] {
	if xs.eof {
		return fun.None[[]byte]()
	}
	chunk := make([]byte, xs.chunkSize)
	cnt, err := xs.reader.Read(chunk)
	switch {
	case err == io.EOF || err == nil && cnt == 0:
		xs.eof = true
		return fun.None[[]byte]()
	case err == nil:
		return fun.Some(chunk[0:cnt])
	default:
		log.Println("Error reading chunk: ", err)
		return fun.None[[]byte]()
	}
}

// ReadByteChunks reads chunks from the reader.
func ReadByteChunks(reader io.Reader, chunkSize int) s.Stream[[]byte] {
	return &readByteChunksImpl{reader, false, chunkSize}
}

type splitByImpl struct {
	s.Stream[[]byte]
	curBuf    []byte
	separator byte
}

func (xs *splitByImpl) Next() fun.Option[[]byte] {
	x := xs.Stream.Next()
	switch {
	case len(xs.curBuf) == 0 && x.IsNone():
		return fun.None[[]byte]()
	case x.IsSome():
		xs.curBuf = append(xs.curBuf, x.Unwrap()...)
		i := 0
		for i < len(xs.curBuf) && xs.curBuf[i] != xs.separator {
			i++
			if i == len(xs.curBuf) {
				x := xs.Stream.Next()
				if x.IsNone() {
					return xs.end()
				}
				xs.curBuf = append(xs.curBuf, x.Unwrap()...)
			}
		}
		return xs.dump(i)
	default:
		i := 0
		for i < len(xs.curBuf) && xs.curBuf[i] != xs.separator {
			i++
			if i == len(xs.curBuf) {
				return xs.end()
			}
		}
		return xs.dump(i)
	}
}

// dump one piece and advance
func (xs *splitByImpl) dump(i int) fun.Option[[]byte] {
	res := xs.curBuf[:i]
	xs.curBuf = xs.curBuf[i+1:]
	return fun.Some(res)
}

// end stream and return everything that is left
func (xs *splitByImpl) end() fun.Option[[]byte] {
	res := xs.curBuf
	xs.curBuf = nil
	return fun.Some(res)
}

// SplitBySeparator splits byte-chunk stream by the given separator.
func SplitBySeparator(xs s.Stream[[]byte], sep byte) s.Stream[[]byte] {
	return &splitByImpl{xs, nil, sep}
}

// ReadLines reads text file line-by-line.
func ReadLines(reader io.Reader) s.Stream[string] {
	chunks := ReadByteChunks(reader, defaultChunkSize)
	rows := SplitBySeparator(chunks, '\n')
	return s.Map(rows, func(x []byte) string { return string(x) })
}
