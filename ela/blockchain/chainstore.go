package blockchain

import (
	. "github.com/elastos/Elastos.ELA/blockchain"
	"github.com/elastos/Elastos.ELA/common/log"
	. "github.com/elastos/Elastos.ELA/core/types"
	"sync"
	"time"
)

type ChainStoreExtend struct {
	IChainStore
	taskChEx chan interface{}
	quitEx   chan chan bool
	mu       sync.Mutex
}

func (c ChainStoreExtend) AddTask(task interface{}) {
	c.taskChEx <- task
}

func NewChainStoreEx(store IChainStore) ChainStoreExtend {
	c := ChainStoreExtend{
		IChainStore: store,
		taskChEx:    make(chan interface{}, TaskChanCap),
		quitEx:      make(chan chan bool, 1),
	}
	DefaultChainStoreEx = c
	go c.loop()
	return c
}

func (c ChainStoreExtend) SaveHistory(block *Block) error {
	//TODO Finish store history
	log.Info("Simulate handle save history")
	return nil
}

func (c ChainStoreExtend) CloseEx() {
	closed := make(chan bool)
	c.quitEx <- closed
	<-closed
}

func (c ChainStoreExtend) loop() {
	for {
		select {
		case t := <-c.taskChEx:
			now := time.Now()
			switch kind := t.(type) {
			case *Block:
				c.SaveHistory(kind)
				tcall := float64(time.Now().Sub(now)) / float64(time.Second)
				log.Debugf("handle SaveHistory time cost: %g num transactions:%d", tcall, len(kind.Transactions))
			}
		case closed := <-c.quitEx:
			closed <- true
			return
		}
	}
}
