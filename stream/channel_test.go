package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendingDataThroughChannel(t *testing.T) {
	t.Parallel()
	ch := make(chan int)
	results := CollectToSlice(FromPairOfChannels(nats10(), ch, ch))
	assert.ElementsMatch(t, results, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
}

func TestStreamConversion(t *testing.T) {
	t.Parallel()
	p := Once(10)
	input := make(chan int)
	output := make(chan int)
	go func() {
		for x := range input {
			output <- x
		}
		close(output)
	}()
	out := FromPairOfChannels(Map(p, mul2), input, output)
	assert.Equal(t, 20, <-out)
}
