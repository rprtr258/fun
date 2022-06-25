package fun

// Pair is a data structure that has two values.
type Pair[A, B any] struct {
	Left  A
	Right B
}

// NewPair constructs the pair.
func NewPair[A, B any](a A, b B) Pair[A, B] {
	return Pair[A, B]{Left: a, Right: b}
}
