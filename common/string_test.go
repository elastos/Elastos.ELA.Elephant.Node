package common

import (
	"testing"
)

func Test_reverseByte(t *testing.T){
	s , _ := ReverseHexString("7e6f5a46740a84998e070d0c6560ce2efa92979c84bb8fb2cae3caaab9c808d5")
	println(s)
}