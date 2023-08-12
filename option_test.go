package fun_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/fun"
)

func TestOptionJSONMarshal(t *testing.T) {
	for name, test := range map[string]struct {
		opt  fun.Option[int]
		want string
	}{
		"valid": {
			opt:  fun.Valid(1),
			want: "1",
		},
		"invalid": {
			opt:  fun.Invalid[int](),
			want: "null",
		},
	} {
		t.Run(name, func(t *testing.T) {
			got, err := json.Marshal(test.opt)
			assert.NoError(t, err)
			assert.Equal(t, test.want, string(got))
		})
	}
}

func TestOptionJSONUnmarshal(t *testing.T) {
	for name, test := range map[string]struct {
		opt  string
		want fun.Option[int]
	}{
		"valid": {
			opt:  "1",
			want: fun.Valid(1),
		},
		"invalid": {
			opt:  "null",
			want: fun.Invalid[int](),
		},
	} {
		t.Run(name, func(t *testing.T) {
			var opt fun.Option[int]
			err := json.Unmarshal([]byte(test.opt), &opt)
			assert.NoError(t, err)
			assert.Equal(t, test.want, opt)
		})
	}
}
