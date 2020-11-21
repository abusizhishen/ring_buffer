package ring_buffer

import "encoding/binary"

func IntToByte(num int) (b []byte) {
	buf := make([]byte, HEAD_SIZE)
	binary.PutVarint(buf, int64(num))
	return buf[:]
}

func ByteToInt(b []byte) (n int) {
	n6, _ := binary.Varint(b)
	return int(n6)
}
