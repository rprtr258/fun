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

var userDecoder = Required("email", String,
	Required("name", String,
		Required("id", Int,
			Success(func(id int) func(name string) func(email string) User {
				return func(name string) func(email string) User {
					return func(email string) User {
						return User{id, name, email}
					}
				}
			}))))

func TestUser(t *testing.T) {
	var result User
	err := userDecoder([]byte(`{"id": 123, "email": "sam@example.com", "name": "Sam"}`), &result)
	assert.NoError(t, err)
	assert.Assert(t, result == User{123, "Sam", "sam@example.com"})
}

func TestUser2(t *testing.T) {
	decoder := Map3(
		func(id int, name string, email string) User {
			return User{id, name, email}
		},
		Field("id", Int),
		Field("name", String),
		Field("email", String),
	)

	var result User
	err := decoder([]byte(`{"id": 123, "email": "sam@example.com", "name": "Sam"}`), &result)
	assert.NoError(t, err)
	assert.Assert(t, result == User{123, "Sam", "sam@example.com"})
}

func TestUserList(t *testing.T) {
	var result []User
	err := List(userDecoder)([]byte(`[{"id": 123, "email": "sam@example.com", "name": "Sam"}]`), &result)
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
