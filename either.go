package fun

// Either is either A value or B value.
type Either[A, B any] struct {
	Left   A
	Right  B
	IsLeft bool
}

// Left constructs Either that is left.
func Left[A, B any](a A) Either[A, B] {
	return Either[A, B]{a, Zero[B](), true}
}

// Right constructs Either that is right.
func Right[A, B any](b B) Either[A, B] {
	return Either[A, B]{Zero[A](), b, false}
}
