package common

import (
	"encoding/hex"
	"testing"
)

func Test_GenDid(t *testing.T) {
	str := "02ECF46B0DE8435DD4E4A93341763F3DDBF12C106C0BE00363B114EFE90F5D2F58"
	b, err := hex.DecodeString(str)
	if err != nil {
		println(err.Error())
	}
	did, err := GenDid(b)
	if err != nil {
		t.Fatal(err)
	}
	println(did)
}
