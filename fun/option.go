package fun

// Option is either value or nothing.
type Option[A any] struct{ ptr *A }

// IsNone checks if option does not contain value.
func (x Option[A]) IsNone() bool {
	return x.ptr == nil
}

// IsSome checks if option does contain value.
func (x Option[A]) IsSome() bool {
	return x.ptr != nil
}

// Unwrap gets value from option if present, SIGSEGV otherwise.
func (x Option[A]) Unwrap() A {
	return *x.ptr
}

// Consume consumes value and calls appropriate function.
func (x Option[A]) Consume(fSome func(A), fNone func()) {
	if x.IsNone() {
		fNone()
	} else {
		fSome(*x.ptr)
	}
}

// None constructs option value with nothing.
func None[A any]() Option[A] {
	return Option[A]{nil}
}

// Some constructs option value with value.
func Some[A any](a A) Option[A] {
	return Option[A]{&a}
}

// FoldOption makes value from option from either value or nothing paths.
func FoldOption[A, B any](x Option[A], fLeft func(A) B, fRight func() B) B {
	if x.IsNone() {
		return fRight()
	}
	return fLeft(*x.ptr)
}

// Map applies function to value if present.
func Map[A, B any](x Option[A], f func(A) B) Option[B] {
	return FoldOption(x, Compose(f, Some[B]), None[B])
}

// FlatMap applies function to value if present.
func FlatMap[A, B any](mx Option[A], f func(A) Option[B]) Option[B] {
	return FoldOption(mx, f, func() Option[B] { return None[B]() })
}
