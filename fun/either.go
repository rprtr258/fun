package fun

// Either is either A value or B value.
type Either[A, B any] struct {
	left  *A
	right *B
}

// IsLeft checks whether the provided Either is left or not.
func (x Either[A, B]) IsLeft() bool {
	return Fold(x, Const[bool, A](true), Const[bool, B](false))
}

// IsRight checks whether the provided Either is right or not.
func (x Either[A, B]) IsRight() bool {
	return Fold(x, Const[bool, A](false), Const[bool, B](true))
}

// Consume consumes either value and calls according callback.
func (x Either[A, B]) Consume(fLeft func(A), fRight func(B)) {
	switch {
	case x.left != nil:
		fLeft(*x.left)
	default:
		fRight(*x.right)
	}
}

// Fold pattern matches Either with two given pattern match handlers.
func Fold[A, B, C any](x Either[A, B], fLeft func(A) C, fRight func(B) C) C {
	switch {
	case x.left != nil:
		return fLeft(*x.left)
	default:
		return fRight(*x.right)
	}
}

// Left constructs Either that is left.
func Left[A, B any](a A) Either[A, B] {
	return Either[A, B]{&a, nil}
}

// Right constructs Either that is right.
func Right[A, B any](b B) Either[A, B] {
	return Either[A, B]{nil, &b}
}
