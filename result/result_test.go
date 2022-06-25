package result

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResult(t *testing.T) {
	t.Parallel()
	io10 := Success(10)
	io20 := Map(io10, func(i int) int { return i * 2 })
	res := FlatMap(io10, func(i int) Result[int] {
		return Map(
			io20,
			func(j int) int {
				return i + j
			},
		)
	})
	assert.Equal(t, res.Unwrap(), 30)
}

func TestErr(t *testing.T) {
	t.Parallel()
	var ptr *string
	ptrio := Success(ptr)
	uptr := FlatMap(ptrio, Dereference[string])
	assert.Equal(t, ErrorNilDeref, uptr.UnwrapErr())
	res := WrapErrf(uptr, "my message %d", 10)
	assert.Equal(t, "my message 10: nil pointer dereference", res.UnwrapErr().Error())
}
