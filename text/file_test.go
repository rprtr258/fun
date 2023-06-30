package text

// import (
// 	"bytes"
// 	"io"
// 	"io/fs"
// 	"os"
// 	"testing"

// 	"github.com/stretchr/testify/assert"

// 	"github.com/rprtr258/go-flow/v2/fun"
// 	"github.com/rprtr258/go-flow/v2/result"
// 	s "github.com/rprtr258/go-flow/v2/stream"
// )

// const exampleText = `
// Line 2
// Line 30
// `

// func openFile(name string) result.Result[*os.File] {
// 	return result.FromGoResult(os.Open(name))
// }

// func TestTextStream(t *testing.T) {
// 	data := []byte(exampleText)
// 	r := bytes.NewReader(data)
// 	lines := ReadLines(r)
// 	linesLengths := s.Map(lines, func(s string) int { return len(s) })
// 	lensSlice := s.CollectToSlice(linesLengths)
// 	assert.Equal(t, []int{0, 6, 7, 0}, lensSlice)
// 	stream10_12 := s.FromMany(10, 11, 12)
// 	stream20_24 := s.Map(stream10_12, func(i int) int { return i * 2 })
// 	res := s.CollectToSlice(stream20_24)
// 	assert.Equal(t, []int{20, 22, 24}, res)
// }

// func TestTextStream2(t *testing.T) {
// 	chunks := ReadByteChunks(bytes.NewReader([]byte(exampleText)), 3)
// 	lines := SplitBySeparator(chunks, '\n')
// 	strings := s.Map(lines, func(x []byte) string { return string(x) })
// 	lens := s.Map(strings, func(s string) int { return len(s) })
// 	lensSlice := s.CollectToSlice(lens)
// 	assert.Equal(t, []int{0, 6, 7, 0}, lensSlice)
// 	stream10_12 := s.FromMany(10, 11, 12)
// 	stream20_24 := s.Map(stream10_12, func(i int) int { return i * 2 })
// 	res := s.CollectToSlice(stream20_24)
// 	assert.Equal(t, []int{20, 22, 24}, res)
// }

// func TestFile(t *testing.T) {
// 	path := t.TempDir() + "/hello.txt"
// 	content := "hello"
// 	assert.NoError(t, os.WriteFile(path, []byte(content), fs.ModePerm))
// 	contentResult := result.FlatMap(openFile(path), func(f *os.File) result.Result[string] {
// 		return result.Eval(func() (str string, err error) {
// 			var bytes []byte
// 			bytes, err = io.ReadAll(f)
// 			if err == nil {
// 				str = string(bytes)
// 			}
// 			return
// 		})
// 	})
// 	str := contentResult.Unwrap()
// 	assert.Equal(t, content, str)
// }

// func TestSplitBySeparator(t *testing.T) {
// 	lines := SplitBySeparator(ReadByteChunks(bytes.NewReader([]byte(exampleText)), defaultChunkSize), '\n')
// 	assert.Equal(t, [][]byte{{}, []byte("Line 2"), []byte("Line 30"), {}}, s.CollectToSlice(lines))
// }

// func TestTextStreamWrite2(t *testing.T) {
// 	lines := ReadLines(bytes.NewReader([]byte(exampleText)))
// 	assert.Equal(t, []string{"", "Line 2", "Line 30", ""}, s.CollectToSlice(lines))
// }

// func TestTextStreamWrite(t *testing.T) {
// 	linesStream := ReadLines(bytes.NewReader([]byte(exampleText)))
// 	lens := s.Map(linesStream, func(s string) int { return len(s) })
// 	lensAsString := s.Map(lens, fun.ToString[int])
// 	w := bytes.NewBuffer([]byte{})
// 	WriteLines(w, lensAsString)
// 	assert.Equal(t, `0
// 6
// 7
// 0`, w.String())
// }

// func TestTextStream3(t *testing.T) {
// 	chunks := ReadByteChunks(bytes.NewReader([]byte(`123
// 456
// `)), 3)
// 	lines := SplitBySeparator(chunks, '\n')
// 	strings := s.Map(lines, func(x []byte) string { return string(x) })
// 	stringsSlice := s.CollectToSlice(strings)
// 	assert.Equal(t, []string{"123", "456", ""}, stringsSlice)
// }

// func TestTextStream4(t *testing.T) {
// 	chunks := ReadByteChunks(bytes.NewReader([]byte(`123
// 456`)), 3)
// 	lines := SplitBySeparator(chunks, '\n')
// 	strings := s.Map(lines, func(x []byte) string { return string(x) })
// 	stringsSlice := s.CollectToSlice(strings)
// 	assert.Equal(t, []string{"123", "456"}, stringsSlice)
// }

// func TestWrite(t *testing.T) {
// 	buf := bytes.NewBuffer(make([]byte, 0, 1000))
// 	WriteByteChunks(buf, s.FromSlice([][]byte{[]byte("a"), []byte("bc")}))
// 	assert.Equal(t, "abc", buf.String())
// }
