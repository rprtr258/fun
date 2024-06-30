package coro

import (
	"io"
	"os"
	"syscall"
	"unsafe"
)

type syscallSubmission struct {
	opcode uint8  // type of operation for this sqe
	fd     int32  // file descriptor to do IO on
	addr   uint64 // pointer to buffer or iovecs
	len    uint32 // buffer size or number of iovecs
	off    uint64 // offset into file
}

type syscallCompletion struct {
	// user_data uint64 // data submission passed back
	res   int32 // result code for this event
	flags uint32
}

type perfectCoro = Coro[syscallCompletion, syscallSubmission]

type system struct {
	yield func(syscallSubmission) syscallCompletion
}

func (s system) read(fd int, p []byte) (int, error) {
	res := s.yield(syscallSubmission{
		opcode: syscall.SYS_READ,
		fd:     int32(fd),
		addr:   uint64(uintptr(unsafe.Pointer(&p[0]))),
		len:    uint32(len(p)),
	})
	return int(res.res), syscall.Errno(res.flags)
}

func (s system) write(fd int, p []byte) (int, error) {
	res := s.yield(syscallSubmission{
		opcode: syscall.SYS_WRITE,
		fd:     int32(fd),
		addr:   uint64(uintptr(unsafe.Pointer(&p[0]))),
		len:    uint32(len(p)),
	})
	return int(res.res), syscall.Errno(res.flags)
}

func newPerfect(f func(system)) perfectCoro {
	return New(func(yield func(syscallSubmission) syscallCompletion) {
		sys := system{func(ss syscallSubmission) syscallCompletion {
			return yield(ss)
		}}
		f(sys)
	})
}

func run(c perfectCoro) error {
	response := syscallCompletion{}
	for {
		req, ok := c.Resume(response)
		if !ok {
			return nil
		}
		r1, r2, errno := syscall.Syscall6(
			uintptr(req.opcode),
			uintptr(req.fd),
			uintptr(req.addr),
			uintptr(req.len),
			uintptr(req.off), 0, 0,
		)
		_ = errno
		response = syscallCompletion{
			res:   int32(r1),
			flags: uint32(r2),
			// user_data: uint64(???),
		}
	}
}

func ExampleCopyFilePerfect(from, to string) error {
	fin, _ := os.Open(from)
	defer fin.Close()

	fout, _ := os.Create(to)
	defer fout.Close()

	return run(newPerfect(func(s system) {
		b := make([]byte, 1024)
		fdin := int(fin.Fd())
		fdout := int(fout.Fd())
		for {
			if _, err := s.read(fdin, b); err != nil {
				if err != io.EOF {
					panic(err.Error())
				}
				break
			}
			_, _ = s.write(fdout, b)
		}
	}))
}
