package common

import (
	"encoding/hex"
	"fmt"
	"testing"
)

type st struct {
	a int
}

func test(s *[]*st) {
	*s = append(*s, &st{1}, &st{2})
}

func Test_getAddress(t *testing.T) {

	var blocks []*st

	test(&blocks)

	fmt.Printf("%v", blocks)

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

	//c1 := make(chan int32, 1000)

	//var i int32 = 0
	//go func() {
	//	for i = 0; i < 2000; {
	//		go func() {
	//			c1 <- atomic.AddInt32(&i, 1)
	//		}()
	//	}
	//}()
	//
	//for {
	//	select {
	//	case d := <-c1:
	//		println(d)
	//	}
	//}
}

func Test_getAddressFromPriv(t *testing.T) {
	priv := "ED77803C1C04AD646C3F0245B6D506EE6DF7A022187921F4D2ABCAF22012F72B"
	addr, err := GetAddressFromPrivKey(priv)
	if err != nil {
		t.Error(err)
	}
	did, err := GetDIDFromPrivKey(priv)
	if err != nil {
		t.Error(err)
	}
	println(addr, did)
}
