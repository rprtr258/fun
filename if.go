package fun

func IF[T any](predicate bool, ifTrue, ifFalse T) T {
	if predicate {
		return ifTrue
	}
	return ifFalse
}

type ifElse[T any] struct {
	value     func() T
	predicate bool
}

func If[T any](predicate bool, value T) ifElse[T] { //nolint:revive // don't export type
	return ifElse[T]{
		predicate: predicate,
		value:     func() T { return value },
	}
}

func IfF[T any](predicate bool, value func() T) ifElse[T] { //nolint:revive // don't export type
	return ifElse[T]{
		predicate: predicate,
		value:     value,
	}
}

func (i ifElse[T]) ElseIf(predicate bool, value T) ifElse[T] {
	if i.predicate {
		return i
	}

	return ifElse[T]{
		predicate: predicate,
		value:     func() T { return value },
	}
}

func (i ifElse[T]) ElseIfF(predicate bool, value func() T) ifElse[T] {
	if i.predicate {
		return i
	}

	return ifElse[T]{
		predicate: predicate,
		value:     value,
	}
}

func (i ifElse[T]) Else(value T) T {
	if i.predicate {
		return i.value()
	}
	return value
}

func (i ifElse[T]) ElseF(value func() T) T {
	if i.predicate {
		return i.value()
	}
	return value()
}

func (i ifElse[T]) ElseDeref(value *T) T {
	if i.predicate {
		return i.value()
	}
	return *value
}
