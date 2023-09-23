package fun

// Set is a collection of distinct elements.
type Set[A comparable] map[A]Unit

// NewSet creates new empty set.
func NewSet[A comparable]() Set[A] {
	return make(Set[A])
}

func (s Set[A]) Iter() func(func(A) bool) bool {
	return func(yield func(A) bool) bool {
		for a := range s {
			if !yield(a) {
				return false
			}
		}
		return true
	}
}

// Contains checks if element is in set.
func (s *Set[A]) Contains(a A) bool {
	_, ok := (*s)[a]
	return ok
}

// Add adds element to the set. If it is already there, does nothing.
func (s *Set[A]) Add(as ...A) {
	for _, a := range as {
		(*s)[a] = Unit1
	}
}

// Remove removes element from set. If it is not there, does nothing.
func (s *Set[A]) Remove(as ...A) {
	for _, a := range as {
		delete(*s, a)
	}
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

func SliceToSet[T comparable](slice []T) Set[T] {
	set := make(Set[T], len(slice))
	for _, elem := range slice {
		set[elem] = Unit{}
	}
	return set
}
