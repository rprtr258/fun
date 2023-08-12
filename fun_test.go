package fun

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func concat(a string, b string) string {
	return a + b
}

func TestFun(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "hello", ConstUnit("hello")(Unit1))
	assert.Equal(t, "hello", Identity("hello"))
	concatc := Curry(concat)
	assert.Equal(t, "ab", concatc("a")("b"))
	assert.Equal(t, "ba", Swap(concatc)("a")("b"))
}

func TestPair(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "a", NewPair("a", "b").Left)
	assert.Equal(t, "b", NewPair("a", "b").Right)
}

func TestEither(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "left", Fold(Left[string, string]("left"), Identity[string], Const[string, string]("other")))
	assert.Equal(t, "other", Fold(Right[string]("right"), Identity[string], Const[string, string]("other")))
}

func TestToString(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "1", ToString(1))
}

func TestCompose(t *testing.T) {
	inc := func(x int) int { return x + 1 }
	double := func(x int) int { return x * 2 }

	f := Compose(inc, double)
	assert.Equal(t, 6, f(2))

	g := Compose(double, inc)
	assert.Equal(t, 5, g(2))
}
