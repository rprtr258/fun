# Iterator and functional utilities
[![Go Reference](https://pkg.go.dev/badge/github.com/rprtr258/go-flow.svg)](https://pkg.go.dev/github.com/rprtr258/go-flow)
[![Go Report Card](https://goreportcard.com/badge/github.com/rprtr258/go-flow)](https://goreportcard.com/report/github.com/rprtr258/go-flow)
![Go](https://github.com/rprtr258/go-flow/workflows/Test/badge.svg?branch=main)
![CodeQL](https://github.com/rprtr258/go-flow/workflows/CodeQL/badge.svg?branch=main)
![Coverage](https://img.shields.io/badge/Coverage-37.1%25-yellow)

The design is inspired by rust [iterators](https://doc.rust-lang.org/std/iter/trait.Iterator.html) and [Result](https://doc.rust-lang.org/std/result/enum.Result.html).

## Result processing
```go
func OpenFile(name string) fun.Result[*os.File] {
	f, err := os.Open(name)
	return fun.Result{f, err, err == nil}
}

func UnmarshalJson[J any](body []byte) fun.Result[J] {
	var j J
	err := json.Unmarshal(body, &j)
  return fun.Result{j, err, err == nil}
}

// LookupEnv gets environment variable
func LookupEnv(varName string) result.Result[string] {
	env, ok := os.LookupEnv(varName)
	return fun.Result{env, fmt.Errorf("env var %q is not defined", varName), ok}
}

// LookupIntEnv gets environment variable and parses it to int
func LookupIntEnv(varName string) result.Result[int] {
	env, ok := os.LookupEnv(varName)
	if !ok {
		return fun.Result{"", fmt.Errorf("env var %q is not defined", varName), false}
	}

	res, err := strconv.Atoi(env)
	return fun.Result{res, err, err == nil}
}

func Process(x int) (int, error) (
	var result fun.Result = // some processing using fun.Result
	return result.Left, result.Right
}
```
## Iter processing
`Iter[V]` is iterator of values of type `V`. They can be finite or infinite. What can be done with iterators:
1. Create iterator either from slice, map, etc. or you can make iterator using `func(yield func(T) bool) bool` function signature.
2. Process iterator.
3. Destroy iterator to single value.

In such way you can build different pipelines. For example there is sample pipeline for counting users' memberships:

![sample flow](doc/flow.png)

How to implement it using iterators:
```go
// Count how many groups users belong to. Groups are:
//   - some community
//   - user's friends
//   - post likers
func CountMemberships(communityID, userID, postID uint) fun.Counter[User] {
	communityMembers := getCommunityMembers(communityID)
	userFriends := getFriends(userID)
	postLikers := getLikers(postID)
	groupsSeqs := iter.Gather([]iter.Seq[iter.Seq[User]]{
		communityMembers,
		userFriends,
		postLikers,
	})
	counters := iter.Map(chans, iter.CollectCounter[User])
	resCounter := iter.Reduce(counters, fun.NewCounter[User](), fun.CounterPlus[User])
	return resCounter
}
```
