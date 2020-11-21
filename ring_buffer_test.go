package ring_buffer

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func TestRingBuffer_Read(t *testing.T) {
	var r = NewRingBuffer(51)
	go write(r)
	read(r)
}

func write(r *RingBuffer) {
	var ss = []string{
		"i love you",
		"are you ok",
		"hello world",
		"1234567890", //10
		"9876",       //4
		"abcdef",     //6
		"ლ(′◉❥◉｀ლ)",
		"helo world",
	}
	for {
		for _, s := range ss {
			_, err := r.Write([]byte(s))
			if err != nil {
				fmt.Println(err)
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(40)))
			}
		}
	}
}

func read(r *RingBuffer) {
	for {
		byt, err := r.Read()
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(50)))
			continue
		}
		fmt.Println("read content : ", string(byt))
	}
}
