package fun

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap_noIndex(t *testing.T) {
	for name, test := range map[string]struct {
		slice    []int
		f        func(int) int
		expected []int
	}{
		"example": {
			slice: []int{1, 2, 3},
			f: func(x int) int {
				return x + 1
			},
			expected: []int{2, 3, 4},
		},
		"empty slice": {
			slice: []int{},
			f: func(x int) int {
				return x + 1
			},
			expected: []int{},
		},
		"nil": {
			slice: nil,
			f: func(x int) int {
				return x + 1
			},
			expected: nil,
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, Map[int](test.f, test.slice...))
		})
	}
}
