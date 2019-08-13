package blockchain

import (
	"bytes"
	"database/sql"
	"encoding/hex"
	"github.com/elastos/Elastos.ELA.Elephant.Node/common"
	"github.com/elastos/Elastos.ELA.Elephant.Node/core/types"
	. "github.com/elastos/Elastos.ELA/blockchain"
	common2 "github.com/elastos/Elastos.ELA/common"
	"github.com/elastos/Elastos.ELA/common/log"
	. "github.com/elastos/Elastos.ELA/core/types"
	"github.com/elastos/Elastos.ELA/core/types/outputpayload"
	"github.com/elastos/Elastos.ELA/core/types/payload"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron"
	"io"
	"os"
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

var (
	MINING_ADDR  = common2.Uint168{}
	ELA_ASSET, _ = common2.Uint256FromHexString("a3d0eaa466df74983b5d7c543de6904f4c9418ead5ffd6d25814234a96db37b0")
)

type ChainStoreExtend struct {
	IChainStore
	IStore
	sql      *sql.DB
	chain    *BlockChain
	taskChEx chan interface{}
	quitEx   chan chan bool
	mu       sync.Mutex
	*cron.Cron
}

func (c ChainStoreExtend) AddTask(task interface{}) {
	c.taskChEx <- task
}

func NewChainStoreEx(chain *BlockChain, chainstore IChainStore, filePath string) (ChainStoreExtend, error) {
	st, err := NewLevelDB(filePath)
	if err != nil {
		return ChainStoreExtend{}, err
	}
	db, err := sql.Open("sqlite3", filePath+"/dpos/dpos.db")
	if err != nil {
		log.Fatal(err)
	}
	c := ChainStoreExtend{
		IChainStore: chainstore,
		IStore:      st,
		sql:         db,
		chain:       chain,
		taskChEx:    make(chan interface{}, 1000),
		quitEx:      make(chan chan bool, 1),
		Cron:        cron.New(),
		mu:          sync.Mutex{},
	}
	DefaultChainStoreEx = c
	go c.loop()
	go c.initTask()
	return c, nil
}

func (c ChainStoreExtend) Close() {

}

func processVote(block *Block, voteTxHolder *map[string]TxType, db *sql.Tx) error {
	for i := 0; i < len(block.Transactions); i++ {
		tx := block.Transactions[i]
		version := tx.Version
		txid, err := common.ReverseHexString(tx.Hash().String())
		if err != nil {
			return err
		}
		if version == 0x09 {
			vout := tx.Outputs
			stmt, err := db.Prepare("insert into chain_vote_info (producer_public_key,vote_type,txid,n,`value`,outputlock,address,block_time,height) values(?,?,?,?,?,?,?,?,?)")
			if err != nil {
				return err
			}
			for _, v := range vout {
				if v.Type == 0x01 && v.AssetID == *ELA_ASSET {
					payload, ok := v.Payload.(*outputpayload.VoteOutput)
					if !ok || payload == nil {
						continue
					}
					contents := payload.Contents
					if !ok {
						continue
					}
					value := v.Value.String()
					n := i
					address, err := v.ProgramHash.ToAddress()
					if err != nil {
						return err
					}
					outputlock := v.OutputLock
					for _, cv := range contents {
						votetype := cv.VoteType
						votetypeStr := ""
						if votetype == 0x00 {
							votetypeStr = "Delegate"
						} else if votetype == 0x01 {
							votetypeStr = "CRC"
						}
						candidates := cv.Candidates
						for _, pub := range candidates {
							_, err := stmt.Exec(common2.BytesToHexString(pub), votetypeStr, txid, n, value, outputlock, address, block.Header.Timestamp, block.Header.Height)
							if err != nil {
								return err
							}
							(*voteTxHolder)[txid] = types.Vote
						}
					}
				}
			}
			stmt.Close()
		}
		// remove canceled vote
		vin := tx.Inputs
		stmt, err := db.Prepare("update chain_vote_info set is_valid = 'NO',cancel_height=? where txid = ? and n = ? ")
		if err != nil {
			return err
		}
		for _, v := range vin {
			txhash := v.Previous.TxID
			vout := v.Previous.Index
			_, err := stmt.Exec(block.Header.Height, txhash, vout)
			if err != nil {
				return err
			}
		}
		stmt.Close()
	}
	return nil
}

func (c ChainStoreExtend) persistTxHistory(block *Block) error {
	txs := block.Transactions
	txhs := make([]types.TransactionHistory, 0)
	pubks := make(map[common2.Uint168][]byte)

	//process vote
	db, err := c.sql.Begin()
	if err != nil {
		return err
	}
	voteTxHolder := new(map[string]TxType)
	err = processVote(block, voteTxHolder, db)
	if err != nil {
		db.Rollback()
		return err
	} else {
		err = db.Commit()
		if err != nil {
			return err
		}
	}

	for i := 0; i < len(txs); i++ {
		tx := txs[i]
		txid, err := common.ReverseHexString(tx.Hash().String())
		if err != nil {
			return err
		}
		var memo []byte
		if len(tx.Attributes) > 0 {
			memo = tx.Attributes[0].Data
		}

		if tx.TxType == CoinBase {
			var to []common2.Uint168
			hold := make(map[common2.Uint168]uint64)
			txhscoinbase := make([]types.TransactionHistory, 0)
			for _, vout := range tx.Outputs {
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
			if (*voteTxHolder)[txid] == types.Vote {
				tx.TxType = types.Vote
			}
			spend := make(map[common2.Uint168]int64)
			var totalInput int64 = 0
			var from []common2.Uint168
			var to []common2.Uint168
			for _, input := range tx.Inputs {
				txid := input.Previous.TxID
				index := input.Previous.Index
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
				err := c.persistTxHistory(kind)
				if err != nil {
					log.Error("Error persist transaction history %s", err.Error())
					os.Exit(-1)
				}
				tcall := float64(time.Now().Sub(now)) / float64(time.Second)
				log.Debugf("handle SaveHistory time cost: %g num transactions:%d", tcall, len(kind.Transactions))
			}
		case closed := <-c.quitEx:
			closed <- true
			return
		}
	}
}

func (c ChainStoreExtend) GetTxHistory(addr string, order string) interface{} {
	key := new(bytes.Buffer)
	key.WriteByte(byte(DataTxHistoryPrefix))
	var txhs interface{}
	if order == "desc" {
		txhs = make(types.TransactionHistorySorterDesc, 0)
	} else {
		txhs = make(types.TransactionHistorySorter, 0)
	}
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
		if txhd.Type == "income" {
			if len(txhd.Inputs) > 0 {
				txhd.Inputs = []string{txhd.Inputs[0]}
			} else {
				txhd.Inputs = []string{}
			}
			txhd.Outputs = []string{txhd.Address}
		} else {
			txhd.Inputs = []string{txhd.Address}
			txhd.Outputs = []string{txhd.Outputs[0]}
		}
		if order == "desc" {
			txhs = append(txhs.(types.TransactionHistorySorterDesc), *txhd)
		} else {
			txhs = append(txhs.(types.TransactionHistorySorter), *txhd)
		}
	}
	if order == "desc" {
		sort.Sort(txhs.(types.TransactionHistorySorterDesc))
	} else {
		sort.Sort(txhs.(types.TransactionHistorySorter))
	}
	return txhs
}

func (c ChainStoreExtend) GetTxHistoryByPage(addr, order string, pageNum, pageSize uint32) (interface{}, int) {
	txhs := c.GetTxHistory(addr, order)
	from := (pageNum - 1) * pageSize
	if order == "desc" {
		return txhs.(types.TransactionHistorySorterDesc).Filter(from, pageSize), len(txhs.(types.TransactionHistorySorterDesc))
	} else {
		return txhs.(types.TransactionHistorySorter).Filter(from, pageSize), len(txhs.(types.TransactionHistorySorter))
	}
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
