package blockchain

import (
	"bytes"
	"encoding/hex"
	types2 "github.com/elastos/Elastos.ELA.Elephant.Node/id/types"
	"github.com/elastos/Elastos.ELA.SideChain/database"
	"github.com/elastos/Elastos.ELA/common"
	"testing"
)

func Test_getDidProperty(t *testing.T) {

	db, err := database.NewLevelDB("elastos_did/data/chain")
	if err != nil {
		t.Fatal(err)
	}
	key := new(bytes.Buffer)
	key.WriteByte(byte(DID_PropertyPrefix))
	did, _ := hex.DecodeString("0354E3AE040052CD61A38DF72E189EDCBBB8BA81599DED25019E286F2E013A3726")
	common.WriteVarBytes(key, did)
	common.WriteVarBytes(key, []byte("clark"))
	it := db.NewIterator([]byte(key.Bytes()))
	for it.Next() {
		val := new(bytes.Buffer)
		val.Write(it.Value())
		dp := types2.DidProperty{}
		dpd, _ := dp.Deserialize(val)
		t.Logf("txid : %s, height = %d , Property_key = %s , Property_value = %s", dpd.Txid, dpd.Height, dpd.Property_key, dpd.Property_value)
	}
	it.Release()

}
