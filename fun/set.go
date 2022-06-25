package fun

// Set is a collection of distinct elements.
type Set[A comparable] map[A]Unit

// NewSet creates new empty set.
func NewSet[A comparable]() Set[A] {
	return make(Set[A])
}

// Contains checks if element is in set.
func (s *Set[A]) Contains(a A) bool {
	_, ok := (*s)[a]
	return ok
}

// Add adds element to the set. If it is already there, does nothing.
func (s *Set[A]) Add(a A) {
	(*s)[a] = Unit1
}

// Remove removes element from set. If it is not there, does nothing.
func (s *Set[A]) Remove(a A) {
	delete(*s, a)
}

// Intersect finds set intersection.
func Intersect[A comparable](as, bs Set[A]) Set[A] {
	res := NewSet[A]()
	for a := range as {
		if bs.Contains(a) {
			res.Add(a)
		}
	}
	return res
}
