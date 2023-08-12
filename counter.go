package fun

// Counter is data structure to count how many elements are in there.
type Counter[A comparable] map[A]int

// NewCounter creates new empty Counter.
func NewCounter[A comparable]() Counter[A] {
	return make(Counter[A])
}

// CounterPlus makes new Counter with elements counts sums.
func CounterPlus[A comparable](a, b Counter[A]) Counter[A] {
	res := make(Counter[A], len(a))
	for k, v := range a {
		res[k] += v
	}
	for k, v := range b {
		res[k] += v
	}
	return res
}
