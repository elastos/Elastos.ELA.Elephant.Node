package blockchain

import (
	"bytes"
	"github.com/elastos/Elastos.ELA.Elephant.Node/ela/core/types"
	"github.com/elastos/Elastos.ELA/blockchain"
	"github.com/elastos/Elastos.ELA/common"
	"testing"
)

func Test_Txhistory(t *testing.T) {
	st, err := blockchain.NewLevelDB("elastos/data/ext")
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
	data, err := st.Get(key.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	txh := types.TransactionHistory{}
	value := new(bytes.Buffer)
	value.Write(data)
	txh.Deserialize(value)
}

func Test_TxHistoryIterator(t *testing.T) {
	st, err := blockchain.NewLevelDB("elastos/data/ext")
	if err != nil {
		t.Fatal(err)
	}
	key := new(bytes.Buffer)
	key.WriteByte(byte(DataTxHistoryPrefix))
	err = common.WriteVarString(key, "EN8WSL4Wt1gM3YjTcHgG7ckBiadtcaNgx4")
	if err != nil {
		t.Fatal(err)
	}
	it := st.NewIterator([]byte(key.Bytes()))
	for it.Next() {
		val := new(bytes.Buffer)
		val.Write(it.Value())
		txh := types.TransactionHistory{}
		txh.Deserialize(val)
		t.Logf("txid : %s, height = %d , value = %d , type = %s", txh.Txid, txh.Height, txh.Value, txh.Type)
	}
	it.Release()
}

func Test_GetTxHistory(t *testing.T) {

	chainStore, err := blockchain.NewChainStore("elastos/data/chain")
	if err != nil {
		t.Fatal(err)
	}
	defer chainStore.Close()
	chainStoreEx, err := NewChainStoreEx(chainStore, "elastos/data/ext")
	if err != nil {
		t.Fatal(err)
	}
	result := chainStoreEx.GetTxHistory("EN8WSL4Wt1gM3YjTcHgG7ckBiadtcaNgx4")
	for _, v := range result {
		t.Log(v)
	}
}
