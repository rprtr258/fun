package json

import "fmt"

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

func example2() {
	type User struct {
		ID    int
		Name  string
		Email string
	}
	newUser := func(id int) func(name string) func(email string) User {
		return func(name string) func(email string) User {
			return func(email string) User {
				return User{id, name, email}
			}
		}
	}

	userDecoder :=
		Required("email", String,
			Required("name", String,
				Required("id", Int,
					Success(newUser))))

	var result User
	if err := userDecoder([]byte(`{"id": 123, "email": "sam@example.com", "name": "Sam"}`), &result); err != nil {
		panic(err)
	}
	fmt.Println(result)
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
