package fun

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounterPlusEmpty(t *testing.T) {
	c1 := map[int]int{}
	c2 := map[int]int{}
	c2[2] = 1
	got := CounterPlus(c1, c2)
	assert.Equal(t, map[int]int{2: 1}, got)
}

func TestCounterPlus(t *testing.T) {
	c1 := map[int]int{1: 1, 2: 2}
	c2 := map[int]int{2: 1}
	got := CounterPlus(c1, c2)
	assert.Equal(t, map[int]int{1: 1, 2: 3}, got)
}
