package blockchain

import (
	. "github.com/elastos/Elastos.ELA/blockchain"
	. "github.com/elastos/Elastos.ELA/core/types"
)

var DefaultChainStoreEx IChainStoreExtend

type IChainStoreExtend interface {
	IChainStore
	persistHistory(block *Block) error
	CloseEx()
	AddTask(task interface{})
}
