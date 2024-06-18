// based on https://research.swtch.com/coro
// NOTE: f does not return value, since in original implementation it is counted as finished value
// along with zero values afterwards, which does not click with me
// Input argument is also removed since it is useless.
package coro

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var ErrCanceled = errors.New("coroutine canceled")

type msg[T any] struct {
	panic any
	val   T
}

type Coro[In, Out any] struct {
	running *bool
	cin     chan msg[In]
	cout    chan msg[Out]
}

func New[In, Out any](
	f func(yield func(Out) In),
) Coro[In, Out] {
	cin := make(chan msg[In])
	cout := make(chan msg[Out], 1)
	running := true
	go func() {
		defer func() {
			running = false
			cout <- msg[Out]{panic: recover()}
		}()
		f(func(out Out) In {
			cout <- msg[Out]{val: out}
			m := <-cin
			if m.panic != nil {
				panic(m.panic)
			}
			return m.val
		})
	}()
	return Coro[In, Out]{
		running: &running,
		cin:     cin,
		cout:    cout,
	}
}

func (c Coro[In, Out]) resume(in In) (out Out, ok bool) {
	if !*c.running {
		return
	}

	isrunning := *c.running
	c.cin <- msg[In]{val: in}
	m := <-c.cout
	if m.panic != nil {
		panic(m.panic)
	}
	return m.val, isrunning
}

func (c Coro[In, Out]) cancel() {
	e := fmt.Errorf("%w", ErrCanceled) // unique wrapper
	c.cin <- msg[In]{panic: e}
	m := <-c.cout
	if m.panic != nil && m.panic != e {
		panic(m.panic)
	}
}

// generator same as New but simplied for case where In is struct{}.
// That is, f is function not accepting values and just producing ones itself.
func generator[T any](
	f func(yield func(T)),
) Coro[struct{}, T] {
	return New(func(yield func(T) struct{}) {
		f(func(o T) {
			yield(o)
		})
	})
}

type GeneratorSet[T any] struct {
	coros []Coro[struct{}, T]
	done  []bool
}

func (s *GeneratorSet[T]) Add(c Coro[struct{}, T]) *bool {
	s.coros = append(s.coros, c)
	s.done = append(s.done, false)
	return &s.done[len(s.done)-1]
}

func (s *GeneratorSet[T]) Resume() (T, bool) {
	selectCases := make([]reflect.SelectCase, len(s.coros))
	for i, c := range s.coros {
		if !*c.running {
			s.done[i] = true
		} else {
			selectCases = append(selectCases, reflect.SelectCase{
				Dir:  reflect.SelectSend,
				Chan: reflect.ValueOf(c.cin),
				Send: reflect.ValueOf(struct{}{}),
			})
		}
	}
	for {
		// chosen, recv, recvOk := reflect.Select(selectCases)
	}
}

type GeneratorSet2[T any] struct {
	coro []Coro[struct{}, T]
	done []bool
}

func newSet2[T any](c1, c2 Coro[struct{}, T]) *GeneratorSet2[T] {
	return &GeneratorSet2[T]{
		[]Coro[struct{}, T]{c1, c2},
		[]bool{false, false},
	}
}

func (s *GeneratorSet2[T]) First() (T, bool) {
	for i, c := range s.coro {
		if !*c.running {
			s.done[i] = true
			return *new(T), false
		}
	}

	selectCases := make([]reflect.SelectCase, len(s.coro))
	for i, c := range s.coro {
		selectCases[i] = reflect.SelectCase{
			Dir:  reflect.SelectSend,
			Chan: reflect.ValueOf(c.cin),
			Send: reflect.ValueOf(msg[struct{}]{val: struct{}{}}),
		}
	}

	couts := make([]reflect.Value, len(s.coro))
	for i := range couts {
		couts[i] = reflect.ValueOf(s.coro[i].cout)
	}

	for {
		chosen, recv, recvOk := reflect.Select(selectCases)
		if selectCases[chosen].Dir == reflect.SelectRecv {
			s.done[chosen] = true
			m := recv.Interface().(msg[T])
			if m.panic != nil {
				panic(m.panic)
			}
			// TODO: what happens with other coroutines?
			return m.val, recvOk // TODO: do we really want to return recvOk here?
		}
		selectCases[chosen] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: couts[chosen],
		}
	}
}

type result[T any] struct {
	t T
	bool
}

func (s *GeneratorSet2[T]) All() []result[T] {
	res := make([]result[T], len(s.coro))
	var wg sync.WaitGroup
	wg.Add(len(s.coro))
	for i, c := range s.coro {
		go func(i int, c Coro[struct{}, T]) {
			t, ok := c.resume(struct{}{})
			res[i] = result[T]{t, ok}
		}(i, c)
	}
	wg.Done()
	return res
}
