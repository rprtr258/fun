package stream

import (
	"github.com/rprtr258/go-flow/fun"
)

// ToChannel sends all stream elements to the given channel.
// When stream is completed, channel is closed.
func ToChannel[A any](xs Stream[A]) <-chan A {
	if chStream, ok := xs.(chanStream[A]); ok {
		return chStream.channel()
	}

	ch := make(chan A)
	go func() {
		ForEach(xs, func(a A) {
			ch <- a
		})
		close(ch)
	}()
	return ch
}

type fromChannelImpl[A any] <-chan A

func (xs fromChannelImpl[A]) Next() fun.Option[A] {
	x, ok := <-xs
	if !ok {
		return fun.None[A]()
	}
	return fun.Some(x)
}

func (xs fromChannelImpl[A]) channel() <-chan A {
	return xs
}

// FromChannel constructs a stream that reads from the given channel.
// When channel is closed, the stream is also closed.
func FromChannel[A any](ch <-chan A) Stream[A] {
	return fromChannelImpl[A](ch)
}

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
	return FromChannel(out)
}
