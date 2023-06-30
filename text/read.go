// Package text provides some utilities to work with text files.
package text

// import (
// 	"bytes"
// 	"io"
// 	"log"

// 	s "github.com/rprtr258/go-flow/v2/stream"
// )

// const defaultChunkSize = 4 * 1024 // 4 KB

// // ReadByteChunks reads chunks of at most chunkSize byte size from the reader.
// func ReadByteChunks(reader io.Reader, chunkSize int) s.Stream[[]byte] {
// 	res := make(chan []byte)
// 	go func() {
// 		buf := make([]byte, chunkSize)
// 		for {
// 			n, err := reader.Read(buf)
// 			res <- buf[:n]
// 			if err != nil {
// 				if err == io.EOF {
// 					break
// 				}
// 				log.Printf("Err during ReadByteChunks: %s", err.Error())
// 				break
// 			}
// 			buf = make([]byte, chunkSize)
// 		}
// 		close(res)
// 	}()
// 	return res
// }

// // SplitBySeparator splits byte-chunk stream by the given separator.
// func SplitBySeparator(xs s.Stream[[]byte], sep byte) s.Stream[[]byte] {
// 	res := make(chan []byte)
// 	go func() {
// 		var curBuf []byte
// 		for chunk := range xs {
// 			curBuf = append(curBuf, chunk...)
// 			for {
// 				idx := bytes.IndexByte(curBuf, sep)
// 				if idx == -1 {
// 					break
// 				}
// 				res <- curBuf[:idx]
// 				curBuf = curBuf[idx+1:]
// 			}
// 		}
// 		res <- curBuf
// 		close(res)
// 	}()
// 	return res
// }

// // ReadLines reads text file line-by-line.
// func ReadLines(reader io.Reader) s.Stream[string] {
// 	chunks := ReadByteChunks(reader, defaultChunkSize)
// 	rows := SplitBySeparator(chunks, '\n')
// 	return s.Map(rows, func(x []byte) string { return string(x) })
// }
