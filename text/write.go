package text

import (
	"io"
	"log"

	s "github.com/rprtr258/go-flow/v2/stream"
)

var endline = "\n"

// WriteByteChunks writes byte chunks to writer.
func WriteByteChunks(writer io.Writer, xs s.Stream[[]byte]) {
	s.ForEach(
		xs,
		func(chunk []byte) {
			cnt, err := writer.Write(chunk)
			if err != nil {
				log.Printf("error writing to %v: %v\n", writer, err)
			} else if cnt != len(chunk) {
				log.Printf("only %d out of %d bytes were written\n", cnt, len(chunk))
			}
		},
	)
}

// MapStringToBytes converts stream of strings to stream of byte chunks.
func MapStringToBytes(stm s.Stream[string]) s.Stream[[]byte] {
	return s.Map(stm, func(s string) []byte { return []byte(s) })
}

// WriteLines creates a sink that receives strings and saves them to writer.
// It adds \n after each line.
func WriteLines(writer io.Writer, xs s.Stream[string]) {
	s.ForEach(
		s.Intersperse(xs, endline),
		func(chunk string) {
			s := []byte(chunk)
			_, err := writer.Write(s)
			if err != nil {
				log.Printf("error writing to %v: %v\n", writer, err)
			}
		},
	)
}
