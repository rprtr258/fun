package fun

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func _const[T any](value T) func() T {
	return func() T {
		return value
	}
}

func TestAll(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		got  string
		want string
	}{
		"true": {
			got:  If(true, "1").Else("2"),
			want: "1",
		},
		"false": {
			got:  If(false, "1").Else("2"),
			want: "2",
		},
		"ThenF ElseF true": {
			got:  IfF(true, _const("1")).ElseF(_const("2")),
			want: "1",
		},
		"ThenF ElseF false": {
			got:  IfF(false, _const("1")).ElseF(_const("2")),
			want: "2",
		},
		"ElseIf Else 1": {
			got: If(true, "1").
				ElseIf(true, "2").
				Else("3"),
			want: "1",
		},
		"ElseIf Else 2": {
			got: If(false, "1").
				ElseIf(true, "2").
				Else("3"),
			want: "2",
		},
		"ElseIf Else 3": {
			got: If(false, "1").
				ElseIf(false, "2").
				Else("3"),
			want: "3",
		},
		"ElseIf ElseIf 3": {
			got: If(false, "1").
				ElseIf(false, "2").
				ElseIf(true, "3").
				Else("4"),
			want: "3",
		},
		"ElseIf ElseIf 4": {
			got: If(false, "1").
				ElseIf(false, "2").
				ElseIf(false, "3").
				Else("4"),
			want: "4",
		},
	} {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.want, test.got)
		})
	}
}
