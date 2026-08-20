package json

import (
	"fmt"
	"testing"

	"github.com/rprtr258/assert"
)

func ExampleAndThen() {
	type Info struct{}

	var infoDecoderV4 Decoder[Info]
	var infoDecoderV3 Decoder[Info]

	infoHelp := func(version int) Decoder[Info] {
		switch version {
		case 3:
			return infoDecoderV3
		case 4:
			return infoDecoderV4
		default:
			return Fail[Info](fmt.Sprintf("Trying to decode info, but version %d is not supported.", version))
		}
	}

	info := AndThen(Int.Field("version"), infoHelp)
	_ = info
}

type User struct {
	ID    int
	Name  string
	Email string
}

var decoderUser = Map3(
	func(id int, name string, email string) User {
		return User{id, name, email}
	},
	Int.Field("id"),
	String.Field("name"),
	String.Field("email"),
)

func TestUser(t *testing.T) {
	t.Parallel()

	result, err := decoderUser.ParseString(`{"id": 123, "email": "sam@example.com", "name": "Sam"}`)
	assert.NoError(t, err)
	assert.Assert(t, result == User{123, "Sam", "sam@example.com"})
}

func TestUserList(t *testing.T) {
	t.Parallel()

	result, err := List(decoderUser).ParseString(`[{"id": 123, "email": "sam@example.com", "name": "Sam"}]`)
	assert.NoError(t, err)
	assert.Equal(t, result, []User{{123, "Sam", "sam@example.com"}})
}

func Example() {
	type Job struct {
		name      string
		id        int
		completed bool
	}

	point := Map3(
		func(name string, id int, completed bool) Job { return Job{name, id, completed} },
		String.Field("name"),
		Int.Field("id"),
		Bool.Field("completed"),
	)
	_ = point
}

func TestListNull(t *testing.T) {
	t.Parallel()

	result, err := List(Any).ParseString(`null`)
	assert.NoError(t, err)
	assert.Equal(t, result, []any(nil))
}
