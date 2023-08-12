package fun_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/fun"
)

var (
	left  = fun.Left[int, int](1)
	right = fun.Right[int](2)
)

func TestIsLeftRight(t *testing.T) {
	assert.True(t, left.IsLeft)
	assert.False(t, right.IsLeft)
}
