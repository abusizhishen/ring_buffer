package ring_buffer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntToByte(t *testing.T) {
	var n = 11
	var byt = IntToByte(n)
	m := ByteToInt(byt)
	assert.Equal(t, n, m)
}
