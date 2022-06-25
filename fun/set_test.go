package fun_test

import (
	"testing"

	"github.com/rprtr258/goflow/fun"
	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	set := fun.NewSet[int]()
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
	set1 := fun.NewSet[int]()
	set1.Add(1)
	set1.Add(2)
	set2 := fun.NewSet[int]()
	set2.Add(1)
	set2.Add(3)
	setIntersection := fun.Intersect(set1, set2)
	assert.True(t, setIntersection.Contains(1))
	assert.False(t, setIntersection.Contains(2))
	assert.False(t, setIntersection.Contains(3))
}
