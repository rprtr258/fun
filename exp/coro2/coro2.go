package main

import (
	"fmt"
	"io"
	"net"
	"syscall"
	"time"
)

func up[T any](h []T, less func(i, j T) bool, j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !less(h[j], h[i]) {
			break
		}
		h[i], h[j] = h[j], h[i]
		j = i
	}
}

// heapPush pushes the element x onto the heap.
// The complexity is O(log n) where n = h.Len().
func heapPush[T any](h *[]T, less func(i, j T) bool, x T) {
	*h = append(*h, x)
	up(*h, less, len(*h)-1)
}

func down[T any](h []T, less func(i, j T) bool, i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && less(h[j2], h[j1]) {
			j = j2 // = 2*i + 2  // right child
		}
		if !less(h[j], h[i]) {
			break
		}
		h[i], h[j] = h[j], h[i]
		i = j
	}
	return i > i0
}

// heapPop removes and returns the minimum element (according to Less) from the heap.
// The complexity is O(log n) where n = h.Len().
// Pop is equivalent to [Remove](h, 0).
func heapPop[T any](h *[]T, less func(i, j T) bool) T {
	n := len(*h) - 1
	(*h)[0], (*h)[n] = (*h)[n], (*h)[0]
	down(*h, less, 0, n)
	return Pop(h)
}

type Sleeper struct {
	ExpireAt time.Time
	F        func()
}

func Pop[T any](h *[]T) T {
	n := len(*h) - 1
	v := (*h)[n]
	*h = (*h)[:n]
	return v
}

type Scheduler struct {
	Ready                     []func()  // TODO: dequeue
	Sleeping                  []Sleeper // heap really
	WaitingRead, WaitingWrite map[int]func()
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		WaitingRead:  map[int]func(){},
		WaitingWrite: map[int]func(){},
	}
}

func less(i, j Sleeper) bool {
	return i.ExpireAt.Before(j.ExpireAt)
}

func (s *Scheduler) CallSoon(f func()) {
	heapPush(&s.Sleeping, less, Sleeper{time.Time{}, f})
}

func (s *Scheduler) CallLater(delay time.Duration, f func()) {
	deadline := time.Now().Add(delay)
	heapPush(&s.Sleeping, less, Sleeper{deadline, f})
}

func (s *Scheduler) CallRead(fd int, f func()) {
	s.WaitingRead[fd] = f
}

func (s *Scheduler) CallWrite(fd int, f func()) {
	s.WaitingWrite[fd] = f
}

func (s *Scheduler) Run() {
	for len(s.Ready) > 0 || len(s.Sleeping) > 0 || len(s.WaitingRead) > 0 || len(s.WaitingWrite) > 0 {
		if len(s.Ready) == 0 {
			// nothing to run now
			timeout := time.Duration(0) // wait forever
			if len(s.Sleeping) > 0 {
				timeout = time.Until(s.Sleeping[0].ExpireAt) // take min deadline
			}

			// wait for either read ready, or write ready, or sleep enough
			epollFd, _ := syscall.EpollCreate1(0)
			events := make([]syscall.EpollEvent, 0, len(s.WaitingRead)+len(s.WaitingWrite))
			for fd := range s.WaitingRead {
				events = append(events, syscall.EpollEvent{Fd: int32(fd), Events: syscall.EPOLLIN})
			}
			for fd := range s.WaitingWrite {
				events = append(events, syscall.EpollEvent{Fd: int32(fd), Events: syscall.EPOLLOUT})
			}
			n, _ := syscall.EpollWait(epollFd, events, int(timeout/time.Millisecond))
			canRead, canWrite := []int{}, []int{}
			for i := 0; i < n; i++ {
				fd := int(events[i].Fd)
				if events[i].Events&syscall.EPOLLIN != 0 {
					canRead = append(canRead, fd)
				}
				if events[i].Events&syscall.EPOLLOUT != 0 {
					canWrite = append(canWrite, fd)
				}
			}
			for _, fd := range canRead {
				s.Ready = append(s.Ready, s.WaitingRead[fd])
			}
			for _, fd := range canWrite {
				s.Ready = append(s.Ready, s.WaitingWrite[fd])
			}

			now := time.Now()
			for len(s.Sleeping) > 0 {
				if now.After(s.Sleeping[0].ExpireAt) {
					sr := heapPop(&s.Sleeping, less)
					s.Ready = append(s.Ready, sr.F) // if timeout
				} // TODO: if not sleeping, push sleep tasks back
			}
		}

		f := s.Ready[0]
		s.Ready = s.Ready[1:]
		f()
	}
}

func Countdown(s *Scheduler, n int) {
	if n > 0 {
		fmt.Println("Down", n)
		time.Sleep(time.Second)
		s.CallLater(4*time.Second, func() { Countdown(s, n-1) })
	}
}

func Countup(s *Scheduler, n int) {
	var _run func(int)
	_run = func(x int) {
		if x < n {
			fmt.Println("Up", x)
			s.CallLater(time.Second, func() { _run(x + 1) })
		}
	}
	_run(0)
}

func Producer(q chan<- int, count int) {
	for n := 0; n < count; n++ {
		fmt.Println("Producing", n)
		q <- n
		time.Sleep(time.Second)
	}
	fmt.Println("Producer done")
	close(q)
}

func Consumer(q <-chan int) {
	for n := range q {
		fmt.Println("Consuming", n)
	}
	fmt.Println("Producer done")
}

type AsyncQueue[T any] struct {
	s       *Scheduler
	items   []T      // TODO: dequeue
	waiting []func() // TODO: dequeue
	closed  bool
}

func NewAsyncQueue[T any](s *Scheduler) *AsyncQueue[T] {
	return &AsyncQueue[T]{s: s}
}

func (q *AsyncQueue[T]) Close() {
	q.closed = true
	if len(q.waiting) > 0 {
		for _, f := range q.waiting {
			q.s.CallSoon(f) // avoid recursion
		}
	}
}

func (q *AsyncQueue[T]) Put(v T) {
	if q.closed {
		panic("closed")
	}

	q.items = append(q.items, v)
	if len(q.waiting) > 0 {
		f := q.waiting[0]
		q.waiting = q.waiting[1:]
		q.s.CallSoon(f) // avoid recursion
	}
}

func (q *AsyncQueue[T]) Get(callback func(T, bool)) {
	if len(q.items) > 0 {
		v := q.items[0]
		q.items = q.items[1:]
		callback(v, true)
	} else {
		if q.closed {
			callback(*new(T), false)
			return
		}
		q.waiting = append(q.waiting, func() { q.Get(callback) })
	}
}

func ProducerCoro(s *Scheduler, q *AsyncQueue[int], count int) {
	var _run func(int)
	_run = func(n int) {
		if n < count {
			fmt.Println("Producing", n)
			q.Put(n)
			s.CallLater(time.Second, func() { _run(n + 1) })
			return
		}
		fmt.Println("Producer done")
		q.Close()
	}
	_run(0)
}

func ConsumerCoro(s *Scheduler, q *AsyncQueue[int]) {
	var _run func(int, bool)
	_run = func(n int, ok bool) {
		if !ok {
			fmt.Println("Consumer done")
			return
		}
		fmt.Println("Consuming", n)
		s.CallSoon(func() { ConsumerCoro(s, q) })
	}
	q.Get(_run)
}

func Generator(yield func(int)) {
	for i := 0; i < 10; i++ {
		yield(i)
	}
}

func EchoHandler(s *Scheduler, conn net.Conn) {
	fd := any(conn).(struct { // TODO: KAL
		fd *struct {
			_    uint64
			_, _ uint32
			fd   int
		}
	}).fd.fd

	b := make([]byte, 1024)
	for {
		s.CallRead(fd, func() {
			n, err := conn.Read(b) // nonblocking
			if err == io.EOF {
				return
			}

			s.CallWrite(fd, func() {
				_, _ = conn.Write(b[:n])
			})
		})
	}
}

func TcpServer(s *Scheduler, addr string) {
	sock, _ := net.Listen("tcp", addr)
	for {
		client, _ := sock.Accept()
		s.CallSoon(func() {
			EchoHandler(s, client)
		})
	}
}

func mainTcpServer(s *Scheduler) {
	s.CallSoon(func() { TcpServer(s, ":8080") })
}

func mainCount(s *Scheduler) {
	s.CallSoon(func() { Countdown(s, 5) })
	s.CallSoon(func() { Countup(s, 5) })
}

func mainConsumerProducer(s *Scheduler) {
	q := NewAsyncQueue[int](s)
	s.CallSoon(func() { ProducerCoro(s, q, 10) })
	s.CallSoon(func() { ConsumerCoro(s, q) })
}

func main() {
	s := NewScheduler()
	// mainCount(s)
	mainConsumerProducer(s)
	s.Run()
}
