package json

import (
	"fmt"
	"testing"

	"github.com/rprtr258/assert"
)

func exampleAndThen() {
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

	info := AndThen(Field("version", Int), infoHelp)
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
	Required("id", Int),
	Required("name", String),
	Required("email", String),
)

func TestUser(t *testing.T) {
	result, err := decoderUser.ParseString(`{"id": 123, "email": "sam@example.com", "name": "Sam"}`)
	assert.NoError(t, err)
	assert.Assert(t, result == User{123, "Sam", "sam@example.com"})
}

func TestUserList(t *testing.T) {
	result, err := List(decoderUser).ParseString(`[{"id": 123, "email": "sam@example.com", "name": "Sam"}]`)
	assert.NoError(t, err)
	assert.Equal(t, result, []User{{123, "Sam", "sam@example.com"}})
}

func example() {
	type Job struct {
		name      string
		id        int
		completed bool
	}

	var point Decoder[Job] = Map3(
		func(name string, id int, completed bool) Job { return Job{name, id, completed} },
		Field("name", String),
		Field("id", Int),
		Field("completed", Bool),
	)
	_ = point
}

func TestListNull(t *testing.T) {
	result, err := List(Any).ParseString(`null`)
	assert.NoError(t, err)
	assert.Equal(t, result, []any(nil))
}
