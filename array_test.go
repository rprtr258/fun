package fun_test

import (
	"fmt"
	"testing"

	"github.com/rprtr258/assert"
	"github.com/rprtr258/fun"
)

func TestChunk(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			got := fun.Chunk(test.chunkSize, test.slice...)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestFlatten(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		slices [][]int
		want   []int
	}{
		"multiple": {
			slices: [][]int{{1, 2}, {3}, {4, 5}},
			want:   []int{1, 2, 3, 4, 5},
		},
		"single": {
			slices: [][]int{{1, 2, 3}},
			want:   []int{1, 2, 3},
		},
		"with empty": {
			slices: [][]int{{}, {1}, {}, {2}},
			want:   []int{1, 2},
		},
		"all empty": {
			slices: [][]int{{}, {}},
			want:   []int{},
		},
		"no args": {
			slices: nil,
			want:   []int{},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := fun.Flatten(test.slices...)
			assert.Equal(t, test.want, got)
		})
	}
}

func ExampleCopy() {
	fmt.Println(fun.Copy(1, 2, 3))
	// Output: [1 2 3]
}
