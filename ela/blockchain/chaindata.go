package blockchain

import (
	"bytes"
	"github.com/elastos/Elastos.ELA.Elephant.Node/ela/core/types"
	"github.com/elastos/Elastos.ELA/common"
	"github.com/elastos/tmp/Elastos.ELA/log"
	"os"
)

func (c ChainStoreExtend) begin() {
	c.NewBatch()
}

func (c ChainStoreExtend) commit() {
	c.BatchCommit()
}

func (c ChainStoreExtend) rollback() {
	c.rollback()
}

// key: DataEntryPrefix + height + address
// value: serialized history
func (c ChainStoreExtend) persistTransactionHistory(txhs []types.TransactionHistory) error {
	c.begin()
	for _, txh := range txhs {
		err := c.doPersistTransactionHistory(txh)
		if err != nil {
			c.rollback()
			log.Fatal("Error persist transaction history")
			os.Exit(-1)
		}
	}
	c.commit()
	return nil
}

func (c ChainStoreExtend) doPersistTransactionHistory(history types.TransactionHistory) error {
	key := new(bytes.Buffer)
	key.WriteByte(byte(DataTxHistoryPrefix))
	err := common.WriteVarString(key, history.Address)
	if err != nil {
		return err
	}
	err = common.WriteUint64(key, history.Height)
	if err != nil {
		return err
	}

	value := new(bytes.Buffer)
	history.Serialize(value)
	c.BatchPut(key.Bytes(), value.Bytes())
	return nil
}
