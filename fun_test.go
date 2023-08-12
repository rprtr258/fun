package fun

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToString(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "1", ToString(1))
}
