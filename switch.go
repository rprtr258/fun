package fun

type switchCase[T comparable, R any] struct {
	predicate T
	result    R
	done      bool
}

// Switch is a pure functional switch/case/default statement.
func Switch[R any, T comparable](predicate T, defVal R) *switchCase[T, R] {
	return &switchCase[T, R]{
		predicate: predicate,
		result:    defVal,
		done:      false,
	}
}

// SwitchZero is a pure functional switch/case/default statement with default
// zero value.
func SwitchZero[R any, T comparable](predicate T) *switchCase[T, R] {
	return Switch(predicate, Zero[R]())
}

func (s *switchCase[T, R]) Case(result R, val T, vals ...T) *switchCase[T, R] {
	if !s.done && s.predicate == val {
		s.result = result
		s.done = true
		return s
	}
	for _, val := range vals {
		if !s.done && s.predicate == val {
			s.result = result
			s.done = true
			break
		}
	}

	return s
}

func (s *switchCase[T, R]) End() R {
	return s.result
}
