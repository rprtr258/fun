package text

import (
	"io"
	"log"

	"github.com/rprtr258/go-flow/stream"
)

var endline = "\n"

// WriteByteChunks writes byte chunks to writer.
func WriteByteChunks(writer io.Writer, xs stream.Stream[[]byte]) {
	stream.ForEach(
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
func MapStringToBytes(stm stream.Stream[string]) stream.Stream[[]byte] {
	return stream.Map(stm, func(s string) []byte { return []byte(s) })
}

// WriteLines creates a sink that receives strings and saves them to writer.
// It adds \n after each line.
func WriteLines(writer io.Writer, xs stream.Stream[string]) {
	stream.ForEach(
		stream.Intersperse(xs, endline),
		func(chunk string) {
			s := []byte(chunk)
			cnt, err := writer.Write(s)
			if err != nil {
				log.Printf("error writing to %v: %v\n", writer, err)
			} else if cnt != len(s) {
				log.Printf("only %d out of %d bytes were written\n", cnt, len(s))
			}
		},
	)
}
