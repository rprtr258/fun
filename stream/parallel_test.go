package stream

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	t.Parallel()
	sourceStream := Take(nats(), 100)
	pool := Scatter(sourceStream, 10)
	for i, st := range pool {
		pool[i] = Map(st, func(id int) int {
			// long processing task
			time.Sleep(10 * time.Millisecond)
			return id
		})
	}
	got := Gather(pool)
	start := time.Now()
	assert.Equal(t, 100, Count(got))
	assert.WithinDuration(t, time.Now(), start, 150*time.Millisecond)
}
