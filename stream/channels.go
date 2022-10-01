package stream

// FromPairOfChannels - takes two channels that are being used to
// talk to some external process and convert them into a single pipe.
// It first starts a separate go routine that will continuously run
// the input stream and send all it's contents to the `input` channel.
// The current thread is left with reading from the output channel.
func FromPairOfChannels[A, B any](xs Stream[A], in chan<- A, out <-chan B) Stream[B] {
	go func() {
		ForEach(xs, func(a A) { in <- a })
		close(in)
	}()
	return out
}
