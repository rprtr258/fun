package text

import (
	"bytes"
	"io"
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/go-flow/fun"
	"github.com/rprtr258/go-flow/result"
	i "github.com/rprtr258/go-flow/result"
	"github.com/rprtr258/go-flow/stream"
)

const exampleText = `
Line 2
Line 30
`

func openFile(name string) result.Result[*os.File] {
	return result.FromGoResult(os.Open(name))
}

func TestTextStream(t *testing.T) {
	data := []byte(exampleText)
	r := bytes.NewReader(data)
	strings := ReadLines(r)
	lens := stream.Map(strings, func(s string) int { return len(s) })
	lensSlice := stream.CollectToSlice(lens)
	assert.Equal(t, []int{0, 6, 7}, lensSlice)
	stream10_12 := stream.FromMany(10, 11, 12)
	stream20_24 := stream.Map(stream10_12, func(i int) int { return i * 2 })
	res := stream.CollectToSlice(stream20_24)
	assert.Equal(t, []int{20, 22, 24}, res)
}

func TestTextStream2(t *testing.T) {
	chunks := ReadByteChunks(bytes.NewReader([]byte(exampleText)), 3)
	rows := SplitBySeparator(chunks, '\n')
	strings := stream.Map(rows, func(x []byte) string { return string(x) })
	lens := stream.Map(strings, func(s string) int { return len(s) })
	lensSlice := stream.CollectToSlice(lens)
	assert.Equal(t, []int{0, 6, 7}, lensSlice)
	stream10_12 := stream.FromMany(10, 11, 12)
	stream20_24 := stream.Map(stream10_12, func(i int) int { return i * 2 })
	res := stream.CollectToSlice(stream20_24)
	assert.Equal(t, []int{20, 22, 24}, res)
}

func TestFile(t *testing.T) {
	path := t.TempDir() + "/hello.txt"
	content := "hello"
	assert.NoError(t, os.WriteFile(path, []byte(content), fs.ModePerm))
	contentResult := i.FlatMap(openFile(path), func(f *os.File) i.Result[string] {
		return i.Eval(func() (str string, err error) {
			var bytes []byte
			bytes, err = io.ReadAll(f)
			if err == nil {
				str = string(bytes)
			}
			return
		})
	})
	str := contentResult.Unwrap()
	assert.Equal(t, content, str)
}

func TestTextStreamWrite(t *testing.T) {
	linesStream := ReadLines(bytes.NewReader([]byte(exampleText)))
	lens := stream.Map(linesStream, func(s string) int { return len(s) })
	lensAsString := stream.Map(lens, fun.ToString[int])
	w := bytes.NewBuffer([]byte{})
	WriteLines(w, lensAsString)
	assert.Equal(t, `0
6
7`, w.String())
}

func TestTextStream3(t *testing.T) {
	chunks := ReadByteChunks(bytes.NewReader([]byte(`123
456
`)), 3)
	rows := SplitBySeparator(chunks, '\n')
	strings := stream.Map(rows, func(x []byte) string { return string(x) })
	lens := stream.Map(strings, func(s string) int { return len(s) })
	lensSlice := stream.CollectToSlice(lens)
	assert.Equal(t, []int{3, 3}, lensSlice)
}

func TestWrite(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, 1000))
	WriteByteChunks(buf, stream.FromSlice([][]byte{[]byte("a"), []byte("bc")}))
	assert.Equal(t, "abc", buf.String())
}
