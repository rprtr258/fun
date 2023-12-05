package orderedmap_test

import (
	"cmp"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/fun"
	"github.com/rprtr258/fun/orderedmap"
)

func TestEmpty(t *testing.T) {
	tree := orderedmap.New[int, int](cmp.Less[int])
	assert.Equal(t, 0, tree.Size())
}

func TestSimple(t *testing.T) {
	tree := orderedmap.New[int, int](cmp.Less[int])
	for k, v := range map[int]int{0: 0, 1: 1, 2: 2, 3: 3, 4: 4, 5: 5} {
		tree.Put(k, v)
	}

	assert.Equal(t, 6, tree.Size())
	for i := 0; i < 6; i++ {
		got, ok := tree.Kth(i)
		assert.True(t, ok)
		assert.Equal(t, i, got)
	}

	min, ok := tree.Min()
	assert.True(t, ok)
	assert.Equal(t, 0, min)

	max, ok := tree.Max()
	assert.True(t, ok)
	assert.Equal(t, 5, max)
}

func TestCrossCheck(t *testing.T) {
	reference := make(map[int]int)
	tree := orderedmap.New[int, int](cmp.Less[int])

	const nops = 1000
	for i := 0; i < nops; i++ {
		switch rand.Intn(2) {
		case 0:
			key, val := rand.Intn(100), rand.Int()
			reference[key] = val
			tree.Put(key, val)
		case 1:
			var del int
			for k := range reference {
				del = k
				break
			}
			delete(reference, del)
			tree.Remove(del)
		}

		assert.Equal(t, len(reference), tree.Size())
		tree.Iter()(func(kv fun.Pair[int, int]) bool {
			assert.Contains(t, reference, kv.K)
			assert.Equal(t, reference[kv.K], kv.V)
			return true
		})
	}
}
