package blockchain

import (
	"bytes"
	"encoding/hex"
	"github.com/elastos/Elastos.ELA.Elephant.Node/common"
	"github.com/elastos/Elastos.ELA.Elephant.Node/core/types"
	. "github.com/elastos/Elastos.ELA/blockchain"
	common2 "github.com/elastos/Elastos.ELA/common"
	"github.com/elastos/Elastos.ELA/common/log"
	. "github.com/elastos/Elastos.ELA/core/types"
	"github.com/elastos/Elastos.ELA/core/types/outputpayload"
	"github.com/elastos/Elastos.ELA/core/types/payload"
	"github.com/robfig/cron"
	"io"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	INCOME string = "income"
	SPEND  string = "spend"
	ELA    uint64 = 100000000
)

var MINING_ADDR = common2.Uint168{}

type ChainStoreExtend struct {
	IChainStore
	IStore
	taskChEx chan interface{}
	quitEx   chan chan bool
	mu       sync.Mutex
	*cron.Cron
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
		taskChEx:    make(chan interface{}, 1000),
		quitEx:      make(chan chan bool, 1),
		Cron:        cron.New(),
		mu:          sync.Mutex{},
	}
	DefaultChainStoreEx = c
	go c.loop()
	go c.initCmc()
	return c, nil
}

func (c ChainStoreExtend) Close() {

}

func (c ChainStoreExtend) persistTxHistory(block *Block) error {
	txs := block.Transactions
	txhs := make([]types.TransactionHistory, 0)
	pubks := make(map[common2.Uint168][]byte)
	for i := 0; i < len(txs); i++ {
		tx := txs[i]
		var memo []byte
		if len(tx.Attributes) > 0 {
			memo = tx.Attributes[0].Data
		}
		if tx.TxType == CoinBase {
			vouts := txs[i].Outputs
			var to []common2.Uint168
			hold := make(map[common2.Uint168]uint64)
			txhscoinbase := make([]types.TransactionHistory, 0)
			for _, vout := range vouts {
				if !common.ContainsU168(vout.ProgramHash, to) {
					to = append(to, vout.ProgramHash)
					txh := types.TransactionHistory{}
					txh.Address = vout.ProgramHash
					txh.Inputs = []common2.Uint168{MINING_ADDR}
					txh.TxType = tx.TxType
					txh.Txid = tx.Hash()
					txh.Height = uint64(block.Height)
					txh.CreateTime = uint64(block.Header.Timestamp)
					txh.Type = []byte(INCOME)
					txh.Fee = 0
					txh.Memo = memo
					hold[vout.ProgramHash] = uint64(vout.Value)
					txhscoinbase = append(txhscoinbase, txh)
				} else {
					hold[vout.ProgramHash] += uint64(vout.Value)
				}
			}
			for i := 0; i < len(txhscoinbase); i++ {
				txhscoinbase[i].Outputs = []common2.Uint168{txhscoinbase[i].Address}
				txhscoinbase[i].Value = hold[txhscoinbase[i].Address]
			}
			txhs = append(txhs, txhscoinbase...)
		} else {
			for _, program := range tx.Programs {
				code := program.Code
				programHash, err := common.GetProgramHash(code[1 : len(code)-1])
				if err != nil {
					continue
				}
				pubks[*programHash] = code[1 : len(code)-1]
			}

			isCrossTx := false
			if tx.TxType == TransferCrossChainAsset {
				isCrossTx = true
			}
			spend := make(map[common2.Uint168]int64)
			var totalInput int64 = 0
			var from []common2.Uint168
			var to []common2.Uint168
			for _, input := range tx.Inputs {
				txid := input.Previous.TxID
				index := input.Previous.Index
				//txResp, err := get("http://" + config.Conf.Ela.Host + TransactionDetail + vintxid)
				referTx, _, err := c.GetTransaction(txid)
				if err != nil {
					return err
				}
				address := referTx.Outputs[index].ProgramHash
				totalInput += int64(referTx.Outputs[index].Value)
				v, ok := spend[address]
				if ok {
					spend[address] = v + int64(referTx.Outputs[index].Value)
				} else {
					spend[address] = int64(referTx.Outputs[index].Value)
				}
				if !common.ContainsU168(address, from) {
					from = append(from, address)
				}
			}
			receive := make(map[common2.Uint168]int64)
			var totalOutput int64 = 0
			vote := outputpayload.VoteOutput{}
			for _, output := range tx.Outputs {
				outputPayload := output.Payload
				if tx.TxType != types.Vote && outputPayload != nil && outputPayload.Validate() == nil {
					var buf bytes.Buffer
					err := outputPayload.Deserialize(&buf)
					if err == nil || err == io.EOF {
						err = vote.Serialize(&buf)
						if err == nil || err == io.EOF {
							tx.TxType = types.Vote
						}
					}
				}

				address, _ := output.ProgramHash.ToAddress()
				var valueCross int64
				if isCrossTx == true && (output.ProgramHash == MINING_ADDR || strings.Index(address, "X") == 0 || address == "4oLvT2") {
					switch pl := tx.Payload.(type) {
					case *payload.TransferCrossChainAsset:
						valueCross = int64(pl.CrossChainAmounts[0])
					}
				}
				if valueCross != 0 {
					totalOutput += valueCross
				} else {
					totalOutput += int64(output.Value)
				}
				v, ok := receive[output.ProgramHash]
				if ok {
					receive[output.ProgramHash] = v + int64(output.Value)
				} else {
					receive[output.ProgramHash] = int64(output.Value)
				}
				if !common.ContainsU168(output.ProgramHash, to) {
					to = append(to, output.ProgramHash)
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
					to = []common2.Uint168{k}
				}

				if transferType == SPEND {
					from = []common2.Uint168{k}
				}

				txh := types.TransactionHistory{}
				txh.Value = uint64(value)
				txh.Address = k
				txh.Inputs = from
				txh.TxType = tx.TxType
				txh.Txid = tx.Hash()
				txh.Height = uint64(block.Height)
				txh.CreateTime = uint64(block.Header.Timestamp)
				txh.Type = []byte(transferType)
				txh.Fee = realFee
				txh.Outputs = to
				txh.Memo = memo
				txhs = append(txhs, txh)
			}

			for k, r := range spend {
				txh := types.TransactionHistory{}
				txh.Value = uint64(r)
				txh.Address = k
				txh.Inputs = []common2.Uint168{k}
				txh.TxType = tx.TxType
				txh.Txid = tx.Hash()
				txh.Height = uint64(block.Height)
				txh.CreateTime = uint64(block.Header.Timestamp)
				txh.Type = []byte(SPEND)
				txh.Fee = uint64(fee)
				txh.Outputs = to
				txh.Memo = memo
				txhs = append(txhs, txh)
			}
		}
	}
	c.persistTransactionHistory(txhs)
	c.persistPbks(pubks)
	return nil
}

func (c ChainStoreExtend) CloseEx() {
	closed := make(chan bool)
	c.quitEx <- closed
	<-closed
	c.Stop()
	log.Info("Extend chainStore shutting down")
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
	txhs := make(types.TransactionHistorySorter,0)
	address, err := common2.Uint168FromAddress(addr)
	if err != nil {
		return txhs
	}
	common2.WriteVarBytes(key, address[:])
	iter := c.NewIterator(key.Bytes())
	defer iter.Release()

	for iter.Next() {
		val := new(bytes.Buffer)
		val.Write(iter.Value())
		txh := types.TransactionHistory{}
		txhd, _ := txh.Deserialize(val)
		txhs = append(txhs, *txhd)
	}
	sort.Sort(txhs)
	return txhs
}

func (c ChainStoreExtend) GetTxHistoryByPage(addr string, pageNum, pageSize uint32) (types.TransactionHistorySorter, int) {
	txhs := c.GetTxHistory(addr)
	from := (pageNum - 1) * pageSize
	return txhs.Filter(from, pageSize), txhs.Len()
}

func (c ChainStoreExtend) GetCmcPrice() types.Cmcs {
	key := new(bytes.Buffer)
	key.WriteByte(byte(DataCmcPrefix))
	common2.WriteVarString(key, "CMC")
	cmcs := types.Cmcs{}
	buf, err := c.Get(key.Bytes())
	if err != nil {
		log.Warn("Can not get Cmc Price data")
		return cmcs
	}
	val := new(bytes.Buffer)
	val.Write(buf)
	cmcs.Deserialize(val)
	return cmcs
}

func (c ChainStoreExtend) GetPublicKey(addr string) string {
	key := new(bytes.Buffer)
	key.WriteByte(byte(DataPkPrefix))
	k, _ := common2.Uint168FromAddress(addr)
	k.Serialize(key)
	buf, err := c.Get(key.Bytes())
	if err != nil {
		log.Warn("No public key find")
		return ""
	}
	return hex.EncodeToString(buf[1:])
}
