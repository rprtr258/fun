package set

// Set is a collection of distinct elements.
type Set[T comparable] struct {
	m map[T]struct{}
}

// New creates new empty set with given capacity.
func New[T comparable](cap int) Set[T] {
	return Set[T]{m: make(map[T]struct{}, cap)}
}

func NewFrom[T comparable](slice []T) Set[T] {
	set := New[T](len(slice))
	set.Add(slice...)
	return set
}

func (s Set[T]) Size() int {
	return len(s.m)
}

func (s Set[T]) Clear() {
	for k := range s.m {
		delete(s.m, k)
	}
}

func (s Set[T]) Copy() Set[T] {
	res := New[T](s.Size())
	for a := range s.m {
		res.Add(a)
	}
	return res
}

func (s Set[T]) Iter() func(func(T) bool) bool {
	return func(yield func(T) bool) bool {
		for a := range s.m {
			if !yield(a) {
				return false
			}
		}
		return true
	}
}

func (s Set[T]) List() []T {
	res := make([]T, 0, s.Size())
	for a := range s.m {
		res = append(res, a)
	}
	return res
}

func (s Set[T]) Merge(as Set[T]) {
	for a := range as.m {
		s.Add(a)
	}
}

func (s Set[T]) RemoveSet(as Set[T]) {
	for a := range as.m {
		delete(s.m, a)
	}
}

func (s Set[T]) PopOk() (T, bool) {
	for a := range s.m {
		delete(s.m, a)
		return a, true
	}

	var zero T
	return zero, false
}

func (s Set[T]) Pop() T {
	res, _ := s.PopOk()
	return res
}

// Contains checks if element is in set.
func (s Set[T]) Contains(a T) bool {
	_, ok := s.m[a]
	return ok
}

func (s Set[T]) ContainsSubset(as Set[T]) bool {
	for a := range as.m {
		if !s.Contains(a) {
			return false
		}
	}
	return true
}

func (s Set[T]) ContainsAny(as ...T) bool {
	for _, a := range as {
		if s.Contains(a) {
			return true
		}
	}
	return false
}

func (s Set[T]) ContainsAll(as ...T) bool {
	for _, a := range as {
		if !s.Contains(a) {
			return false
		}
	}
	return true
}

func (s Set[T]) IsEqual(as Set[T]) bool {
	if s.Size() != as.Size() {
		return false
	}

	for a := range as.m {
		if !s.Contains(a) {
			return false
		}
	}
	return true
}

// Add adds element to the set. If it is already there, does nothing.
func (s Set[T]) Add(as ...T) {
	for _, a := range as {
		s.m[a] = struct{}{}
	}
}

// Remove removes element from set. If it is not there, does nothing.
func (s Set[T]) Remove(as ...T) {
	for _, a := range as {
		delete(s.m, a)
	}
}

// Subtract bss from as. Returns set with elements which are in as but not in
// any of bss.
func Subtract[T comparable](as Set[T], bss ...Set[T]) Set[T] {
	res := New[T](0)
OUTER:
	for a := range as.m {
		for _, bs := range bss {
			if bs.Contains(a) {
				continue OUTER
			}
		}
		res.Add(a)
	}
	return res
}

// Intersect finds sets intersection.
func Intersect[T comparable](as Set[T], bss ...Set[T]) Set[T] {
	res := New[T](0)
OUTER:
	for a := range as.m {
		for _, bs := range bss {
			if !bs.Contains(a) {
				continue OUTER
			}
		}
		res.Add(a)
	}
	return res
}

func Union[T comparable](as Set[T], bss ...Set[T]) Set[T] {
	res := as.Copy()
	for _, bs := range bss {
		for b := range bs.m {
			res.Add(b)
		}
	}
	return res
}

func SymmetricDifference[T comparable](as, bs Set[T]) Set[T] {
	res := New[T](0)
	for a := range as.m {
		if !bs.Contains(a) {
			res.Add(a)
		}
	}
	for b := range bs.m {
		if !as.Contains(b) {
			res.Add(b)
		}
	}
	return res
}
