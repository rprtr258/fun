package fun_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/goflow/fun"
)

var (
	left  = fun.Left[int, int](1)
	right = fun.Right[int](2)
)

func TestIsLeftRight(t *testing.T) {
	assert.True(t, left.IsLeft(), right.IsRight())
	assert.False(t, left.IsRight(), right.IsLeft())
}

func TestConsume(t *testing.T) {
	var leftInt int
	left.Consume(func(x int) { leftInt = x }, func(x int) { t.Fail() })
	assert.Equal(t, leftInt, 1)

	var rightInt int
	right.Consume(func(x int) { t.Fail() }, func(x int) { rightInt = x })
	assert.Equal(t, rightInt, 2)
}
