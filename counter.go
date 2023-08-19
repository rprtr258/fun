package fun

// CounterPlus makes new Counter with elements counts sums.
func CounterPlus[A comparable](a, b map[A]int) map[A]int {
	res := make(map[A]int, len(a))
	for k, v := range a {
		res[k] += v
	}
	for k, v := range b {
		res[k] += v
	}
	return res
}
