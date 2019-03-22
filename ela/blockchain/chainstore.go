package blockchain

import (
	"bytes"
	"github.com/elastos/Elastos.ELA.Elephant.Node/common"
	"github.com/elastos/Elastos.ELA.Elephant.Node/ela/core/types"
	. "github.com/elastos/Elastos.ELA/blockchain"
	"github.com/elastos/Elastos.ELA/common/log"
	. "github.com/elastos/Elastos.ELA/core/types"
	"github.com/elastos/Elastos.ELA/core/types/payload"
	common2 "github.com/xiaomingfuckeasylife/Elastos.ELA/common"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	INCOME      string = "income"
	SPEND       string = "spend"
	MINING_ADDR string = "0000000000000000000000000000000000"
	ELA         uint64 = 100000000
)

var txTypeEnum = map[TransactionType]string{
	CoinBase:                "CoinBase",
	RegisterAsset:           "RegisterAsset",
	TransferAsset:           "TransferAsset",
	Record:                  "Record",
	Deploy:                  "Deploy",
	SideChainPow:            "SideChainPow",
	RechargeToSideChain:     "RechargeToSideChain",
	WithdrawFromSideChain:   "WithdrawFromSideChain",
	TransferCrossChainAsset: "TransferCrossChainAsset",
}

type ChainStoreExtend struct {
	IChainStore
	IStore
	taskChEx chan interface{}
	quitEx   chan chan bool
	mu       sync.Mutex
}

func (c ChainStoreExtend) AddTask(task interface{}) {
	c.taskChEx <- task
}

func NewChainStoreEx(chainstore IChainStore, filePath string) (ChainStoreExtend, error) {
	st, err := NewLevelDB(filePath)
	if err != nil {
		return ChainStoreExtend{}, err
	}
	c := ChainStoreExtend{
		IChainStore: chainstore,
		IStore:      st,
		taskChEx:    make(chan interface{}, TaskChanCap),
		quitEx:      make(chan chan bool, 1),
	}
	DefaultChainStoreEx = c
	go c.loop()
	return c, nil
}

func (c ChainStoreExtend) Close() {

}

func (c ChainStoreExtend) persistTxHistory(block *Block) error {
	txs := block.Transactions
	txhs := make([]types.TransactionHistory, 0)
	for i := 0; i < len(txs); i++ {
		tx := txs[i]
		if tx.TxType == CoinBase {
			vouts := txs[i].Outputs
			var to []string
			hold := make(map[string]types.TransactionHistory)
			for _, vout := range vouts {
				address, _ := vout.ProgramHash.ToAddress()
				if !common.Contains(address, to) {
					to = append(to, address)
					txh := types.TransactionHistory{}
					txh.Value = uint64(vout.Value)
					txh.Address = address
					txh.Inputs = []string{MINING_ADDR}
					txh.TxType = txTypeEnum[tx.TxType]
					txh.Txid, _ = common.ReverseHexString(tx.Hash().String())
					txh.Height = uint64(block.Height)
					txh.CreateTime = uint64(block.Header.Timestamp)
					txh.Type = INCOME
					txh.Fee = 0
					txhs = append(txhs, txh)
				} else {
					txh := hold[address]
					txh.Value += uint64(vout.Value)
				}
			}
			for _, txh := range txhs {
				txh.Outputs = to
			}
		} else {
			isCrossTx := false
			if tx.TxType == TransferCrossChainAsset {
				isCrossTx = true
			}
			spend := make(map[string]int64)
			var totalInput int64 = 0
			var from []string
			var to []string
			for _, input := range tx.Inputs {
				txid := input.Previous.TxID
				index := input.Previous.Index
				//txResp, err := get("http://" + config.Conf.Ela.Host + TransactionDetail + vintxid)
				referTx, _, err := c.GetTransaction(txid)
				if err != nil {
					return err
				}
				address, _ := referTx.Outputs[index].ProgramHash.ToAddress()
				totalInput += int64(referTx.Outputs[index].Value)
				v, ok := spend[address]
				if ok {
					spend[address] = v + int64(referTx.Outputs[index].Value)
				} else {
					spend[address] = int64(referTx.Outputs[index].Value)
				}
				if !common.Contains(address, from) {
					from = append(from, address)
				}
			}
			receive := make(map[string]int64)
			var totalOutput int64 = 0
			for _, output := range tx.Outputs {
				address, _ := output.ProgramHash.ToAddress()
				var valueCross int64
				if isCrossTx == true && (address == MINING_ADDR || strings.Index(address, "X") == 0) {
					switch pl := tx.Payload.(type) {
					case *payload.PayloadTransferCrossChainAsset:
						valueCross = int64(pl.CrossChainAmounts[0])
					}
				}
				if valueCross != 0 {
					totalOutput += valueCross
				} else {
					totalOutput += int64(output.Value)
				}
				v, ok := receive[address]
				if ok {
					receive[address] = v + int64(output.Value)
				} else {
					receive[address] = int64(output.Value)
				}
				if !common.Contains(address, to) {
					to = append(to, address)
				}
			}
			fee := totalInput - totalOutput
			for k, r := range receive {
				transferType := INCOME
				s, ok := spend[k]
				var value int64
				if ok {
					if s > r {
						value = s - r
						transferType = SPEND
					} else {
						value = r - s
					}
					delete(spend, k)
				} else {
					value = r
				}
				var realFee uint64 = uint64(fee)
				if transferType == INCOME {
					realFee = 0
				}
				txh := types.TransactionHistory{}
				txh.Value = uint64(value)
				txh.Address = k
				txh.Inputs = from
				txh.TxType = txTypeEnum[tx.TxType]
				txh.Txid, _ = common.ReverseHexString(tx.Hash().String())
				txh.Height = uint64(block.Height)
				txh.CreateTime = uint64(block.Header.Timestamp)
				txh.Type = transferType
				txh.Fee = realFee
				txh.Outputs = to
				txhs = append(txhs, txh)
			}

			for k, r := range spend {
				txh := types.TransactionHistory{}
				txh.Value = uint64(r)
				txh.Address = k
				txh.Inputs = from
				txh.TxType = txTypeEnum[tx.TxType]
				txh.Txid, _ = common.ReverseHexString(tx.Hash().String())
				txh.Height = uint64(block.Height)
				txh.CreateTime = uint64(block.Header.Timestamp)
				txh.Type = SPEND
				txh.Fee = uint64(fee)
				txh.Outputs = to
				txhs = append(txhs, txh)
			}
		}
	}
	c.persistTransactionHistory(txhs)
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
				c.persistTxHistory(kind)
				tcall := float64(time.Now().Sub(now)) / float64(time.Second)
				log.Debugf("handle SaveHistory time cost: %g num transactions:%d", tcall, len(kind.Transactions))
			}
		case closed := <-c.quitEx:
			closed <- true
			return
		}
	}
}

func (c ChainStoreExtend) GetTxHistory(addr string) types.TransactionHistorySorter {
	key := new(bytes.Buffer)
	key.WriteByte(byte(DataTxHistoryPrefix))
	common2.WriteVarString(key, addr)

	iter := c.NewIterator(key.Bytes())
	defer iter.Release()
	var txhs types.TransactionHistorySorter
	for iter.Next() {
		val := new(bytes.Buffer)
		val.Write(iter.Value())
		txh := types.TransactionHistory{}
		txh.Deserialize(val)
		txhs = append(txhs, txh)
	}
	sort.Sort(txhs)
	return txhs
}

func (c ChainStoreExtend) GetTxHistoryByPage(addr string, pageNum, pageSize uint32) types.TransactionHistorySorter {
	txhs := c.GetTxHistory(addr)
	from := (pageNum - 1) * pageSize
	return txhs.Filter(from, pageSize)
}
