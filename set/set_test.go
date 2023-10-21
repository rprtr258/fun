package set_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/fun/set"
)

func TestSet(t *testing.T) {
	set := set.New[int](0)
	assert.False(t, set.Contains(1))
	assert.False(t, set.Contains(2))
	set.Add(1)
	assert.True(t, set.Contains(1))
	assert.False(t, set.Contains(2))
	set.Add(2)
	assert.True(t, set.Contains(1))
	assert.True(t, set.Contains(2))
	set.Remove(2)
	assert.True(t, set.Contains(1))
	assert.False(t, set.Contains(2))
}

func TestSetIntersect(t *testing.T) {
	set1 := set.New[int](0)
	set1.Add(1)
	set1.Add(2)
	set2 := set.New[int](0)
	set2.Add(1)
	set2.Add(3)
	setIntersection := set.Intersect(set1, set2)
	assert.True(t, setIntersection.Contains(1))
	assert.False(t, setIntersection.Contains(2))
	assert.False(t, setIntersection.Contains(3))
}
