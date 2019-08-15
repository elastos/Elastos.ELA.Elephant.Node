package common

import (
	"encoding/hex"
	"sync/atomic"
	"testing"
)

func Test_getAddress(t *testing.T) {
	b, _ := hex.DecodeString("02776f2ad3fc822caa976bf0a83eb33cf4047518c9b6d411603be4a864b24acb4b")
	s, _ := GetProgramHash(b)
	addr, _ := s.ToAddress()
	println(addr)

	str, _ := ReverseHexString("a3d0eaa466df74983b5d7c543de6904f4c9418ead5ffd6d25814234a96db37b0")
	println(str)
	c := make(chan bool, 1)
	println("chan ", len(c))
	c <- true
	println("chan ", len(c))
	v := <-c
	println("chan ", len(c), v)
	defer println(1)
	defer println(2)

	c1 := make(chan int32, 1000)

	var i int32 = 0
	go func() {
		for i = 0; i < 2000; {
			go func() {
				c1 <- atomic.AddInt32(&i, 1)
			}()
		}
	}()

	for {
		select {
		case d := <-c1:
			println(d)
		}
	}
}
