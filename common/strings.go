package common

import "encoding/hex"

func ReverseHexString(s string) (string,error){
	b , err := hex.DecodeString(s)
	if err != nil {
		return "",err
	}
	for i:=0; i<len(b)/2; i++ {
		b[i],b[len(b)-i-1] = b[len(b)-i-1],b[i]
	}
	return hex.EncodeToString(b) , nil
}