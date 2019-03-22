package blockchain

import (
	"github.com/elastos/Elastos.ELA.Elephant.Node/ela/core/types"
	. "github.com/elastos/Elastos.ELA/blockchain"
	. "github.com/elastos/Elastos.ELA/core/types"
)

var DefaultChainStoreEx IChainStoreExtend

type IChainStoreExtend interface {
	IChainStore
	persistTxHistory(block *Block) error
	CloseEx()
	AddTask(task interface{})
	GetTxHistory(addr string) types.TransactionHistorySorter
	GetTxHistoryByPage(addr string, pageNum, pageSize uint32) types.TransactionHistorySorter
}
