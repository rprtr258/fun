package fun

// Result represents a calculation that will yield a value of type A once executed.
// The calculation might as well fail.
// It is designed to not panic ever.
type Result[A any] Pair[A, error]
