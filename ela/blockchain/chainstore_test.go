package blockchain

import (
	"bytes"
	"github.com/elastos/Elastos.ELA.Elephant.Node/ela/core/types"
	"github.com/elastos/Elastos.ELA/common"
	"github.com/elastos/Elastos.ELA/blockchain"
	"testing"
)

func Test_Txhistory(t *testing.T){
	st, err := blockchain.NewLevelDB("/Users/clark/workspace/golang/src/github.com/elastos/Elastos.ELA.Elephant.Node/elastos/data/ext")
	if err != nil {
		t.Fatal(err)
	}
	key := new(bytes.Buffer)
	key.WriteByte(byte(DataTxHistoryPrefix))
	err = common.WriteUint64(key, 243738)
	if err != nil {
		t.Fatal(err)
	}
	err = common.WriteVarString(key, "EY7x8rCoK3aeHBD31xUFYFzFMQFVgYbG5N")
	if err != nil {
		t.Fatal(err)
	}
	data , err := st.Get(key.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	txh := types.TransactionHistory{}
	value := new(bytes.Buffer)
	value.Write(data)
	txh.Deserialize(value)
}
