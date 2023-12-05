package fun_test

import (
	"testing"

	"github.com/rprtr258/fun"
	"github.com/stretchr/testify/assert"
)

func TestChunk(t *testing.T) {
	for name, test := range map[string]struct {
		slice     []int
		chunkSize int
		want      [][]int
	}{
		"even": {
			slice:     []int{1, 2, 3, 4, 5, 6},
			chunkSize: 2,
			want:      [][]int{{1, 2}, {3, 4}, {5, 6}},
		},
		"uneven": {
			slice:     []int{1, 2, 3, 4, 5},
			chunkSize: 2,
			want:      [][]int{{1, 2}, {3, 4}, {5}},
		},
	} {
		t.Run(name, func(t *testing.T) {
			got := fun.Chunk(test.chunkSize, test.slice...)
			assert.Equal(t, test.want, got)
		})
	}
}
