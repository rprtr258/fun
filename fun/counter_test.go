package fun

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounterPlusEmpty(t *testing.T) {
	c1 := NewCounter[int]()
	c2 := NewCounter[int]()
	c2[2] = 1
	got := CounterPlus(c1, c2)
	assert.Equal(t, Counter[int](map[int]uint{2: 1}), got)
}

func TestCounterPlus(t *testing.T) {
	c1 := NewCounter[int]()
	c1[1] = 1
	c1[2] = 2
	c2 := NewCounter[int]()
	c2[2] = 1
	got := CounterPlus(c1, c2)
	assert.Equal(t, Counter[int](map[int]uint{1: 1, 2: 3}), got)
}
