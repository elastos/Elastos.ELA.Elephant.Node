package common

import (
	"encoding/hex"
	"testing"
)

func Test_getAddress(t *testing.T) {
	b, _ := hex.DecodeString("02776f2ad3fc822caa976bf0a83eb33cf4047518c9b6d411603be4a864b24acb4b")
	s, _ := GetProgramHash(b)
	addr, _ := s.ToAddress()
	println(addr)
}
