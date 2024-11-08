package fun

import (
	"testing"

	"github.com/rprtr258/assert"
)

func TestToString(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "1", ToString(1))
}
