package ring_buffer

import (
	"errors"
	"log"
	"os"
	"time"
)
import (
	"encoding/binary"
)

const HEAD_SIZE = binary.MaxVarintLen32

var (
	ErrWaitingRead                  = errors.New("ring_buffer is waiting for read")
	ErrWaitingWrite                 = errors.New("ring_buffer waiting for write")
	ErrDataSizeBiggerThanBufferSize = errors.New("data size is bigger than bing_buffer size")
)

type RingBuffer struct {
	cap, head, tail int
	buf             []byte
	//sync.RWMutex
}

func NewRingBuffer(n int) *RingBuffer {
	return &RingBuffer{
		cap: n,
		buf: make([]byte, n),
	}
}

func (r *RingBuffer) Read() ([]byte, error) {
	//r.RLock()
	//defer r.RUnlock()

	//check head and tail to avoid read dirty data
	if r.head == r.tail {
		return nil, ErrWaitingWrite
	}
	n := r.readLength()
	if n == 0 {
		return nil, ErrWaitingWrite
	}

	if n < 0 {
		log.Printf("buf: %+v", r)
		log.Printf("read length:%d", n)
		time.Sleep(time.Second)
		os.Exit(0)
	}
	b := make([]byte, n)
	r.tail += HEAD_SIZE
	r.read(b)
	r.tail += n

	return b, nil
}

func (r *RingBuffer) Write(b []byte) (n int, err error) {
	if len(b) == 0 {
		return
	}

	if len(b) > r.cap-HEAD_SIZE {
		return 0, ErrDataSizeBiggerThanBufferSize
	}
	//r.Lock()
	//defer r.Unlock()

	if r.IsWriteIndexOverLoad(len(b)) {
		err = ErrWaitingRead
		return
	}

	r.writeLength(len(b))
	r.write(b)
	return
}

func (r *RingBuffer) write(b []byte) {
	start := r.head % r.cap
	if end := start + len(b); end <= r.cap {
		copy(r.buf[start:], b)
	} else {
		copy(r.buf[start:], b[:r.cap-start])
		copy(r.buf[:end-r.cap], b[r.cap-start:])
	}

	r.head += len(b)
	return
}

func (r *RingBuffer) read(b []byte) int {
	start := r.tail % r.cap

	if end := start + len(b); end <= r.cap {
		copy(b, r.buf[start:end])
	} else {
		copy(b, r.buf[start:])
		copy(b[r.cap-start:], r.buf[:end-r.cap])
	}

	return r.tail + len(b)
}

func (r *RingBuffer) readLength() (n int) {
	var b = make([]byte, HEAD_SIZE)
	r.read(b)
	return ByteToInt(b)
}

func (r *RingBuffer) writeLength(n int) {
	r.write(IntToByte(n))
	return
}

func (r *RingBuffer) IsWriteIndexOverLoad(n int) bool {
	return r.cap-(r.head-r.tail)-HEAD_SIZE < n
}
