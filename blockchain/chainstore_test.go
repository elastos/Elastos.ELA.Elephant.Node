package blockchain

import (
	"bytes"
	"github.com/elastos/Elastos.ELA.Elephant.Node/core/types"
	"github.com/elastos/Elastos.ELA/blockchain"
	"github.com/elastos/Elastos.ELA/common"
	"github.com/elastos/Elastos.ELA/common/config"
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

	genesisBlock := config.GenesisBlock(&common.Uint168{
		0x12, 0xc8, 0xa2, 0xe0, 0x67, 0x72, 0x27,
		0x14, 0x4d, 0xf8, 0x22, 0xb7, 0xd9, 0x24,
		0x6c, 0x58, 0xdf, 0x68, 0xeb, 0x11, 0xce,
	})
	chainStore, err := blockchain.NewChainStore("elastos/data/chain", genesisBlock)
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

func Test_GetCmcPrice(t *testing.T) {

	genesisBlock := config.GenesisBlock(&common.Uint168{
		0x12, 0xc8, 0xa2, 0xe0, 0x67, 0x72, 0x27,
		0x14, 0x4d, 0xf8, 0x22, 0xb7, 0xd9, 0x24,
		0x6c, 0x58, 0xdf, 0x68, 0xeb, 0x11, 0xce,
	})
	chainStore, err := blockchain.NewChainStore("elastos/data/chain", genesisBlock)
	if err != nil {
		t.Fatal(err)
	}
	defer chainStore.Close()
	chainStoreEx, err := NewChainStoreEx(chainStore, "elastos/data/ext")
	if err != nil {
		t.Fatal(err)
	}
	result := chainStoreEx.GetCmcPrice()
	t.Logf("%v", result)
}

func Test_TxCmcWithBareLeverDb(t *testing.T) {
	st, err := blockchain.NewLevelDB("elastos/data/ext")
	if err != nil {
		t.Fatal(err)
	}
	key := new(bytes.Buffer)
	key.WriteByte(byte(DataCmcPrefix))
	err = common.WriteVarString(key, "CMC")
	if err != nil {
		t.Fatal(err)
	}
	buf, err := st.Get(key.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	val := new(bytes.Buffer)
	val.Write(buf)
	cmcs := types.Cmcs{}
	cmcs.Deserialize(val)
	t.Logf("Name : %s , Price USD : %s", cmcs.C[0].Name, cmcs.C[0].Price_usd)
}
