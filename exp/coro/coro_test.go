package coro

import "testing"

func TestRange(t *testing.T) {
	myrange := generator(func(yield func(int)) {
		t.Log("started")
		for i := 0; i < 10; i++ {
			t.Log(">", i)
			yield(i)
			t.Log("<", i)
		}
	})
	for i := 0; i < 10; i++ {
		j, _ := myrange.Resume(struct{}{})
		if j != i {
			t.Fatalf("expected %d, got %d", i, j)
		}
	}
}

type Status uint8

const (
	Success Status = iota
	NeedMoreInput
	BadInput
)

func TestParser(t *testing.T) {
	const input = `"hello world"`
	parser := New(func(yield func(Status) byte) {
		read := func() byte { return yield(NeedMoreInput) }
		if read() != '"' {
			yield(BadInput)
			return
		}
		var c byte
		for c != '"' {
			c = read()
			if c == '\\' {
				read()
			}
		}
		yield(Success)
	})

	for _, b := range []byte(input) {
		status, ok := parser.Resume(b)
		switch {
		case status == BadInput:
			t.Fatal("bad input")
		case status == Success:
			// done
			return
		case !ok:
			t.Fatal("ran out of input")
		}
	}
}

func counter() Coro[struct{}, int] {
	return New(func(yield func(int) struct{}) {
		for i := 2; ; i++ {
			yield(i)
		}
	})
}

func filter(p int, next func(struct{}) (int, bool)) Coro[struct{}, int] {
	return New(func(yield func(int) struct{}) {
		for {
			n, ok := next(struct{}{})
			if !ok {
				return
			}
			if n%p != 0 {
				yield(n)
			}
		}
	})
}

func TestPrimeSieve(t *testing.T) {
	gen := counter()
	t.Cleanup(gen.Cancel)

	for _, pp := range [...]int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29} {
		p, _ := gen.Resume(struct{}{})
		if p != pp {
			t.Fatalf("expected %d, got %d", pp, p)
		}
		gen = filter(p, gen.Resume)
	}
}

func TestCallUntilExhausted(t *testing.T) {
	gen := generator(func(yield func(string)) {
		yield("hello")
		yield("world")
		yield("done")
	})
	for i, step := range [...]struct {
		string
		bool
	}{
		{"hello", true},
		{"world", true},
		{"done", true},
		{"", false},
	} {
		s, ok := gen.Resume(struct{}{})
		if ok != step.bool || s != step.string {
			t.Fatalf(`failed on step %d:
s: expected %q, got %q
ok: expected %v, got %v`, i, step.string, s, step.bool, ok)
		}
	}
}

// TODO: make it work
func useDB(t *testing.T) Coro[struct{}, *int] {
	return generator(func(yield func(*int)) {
		db := new(int)
		*db = 1 // open db
		t.Log("open db")
		yield(db)
		t.Log("close db")
		*db = -1 // close db
	})
}

func TestUseDB(t *testing.T) {
	t.Skip()

	_ = useDB(t)
	// defer close()
	t.Fail()
	// db, ok := get()
	// t.Log("get db", *db, ok)
	// if !ok {
	// 	t.Fatal("failed")
	// }
	// t.Log("use db", *db)
	// if *db != 1 {
	// 	t.Fatalf("expected 1, got %d", *db)
	// }
}
