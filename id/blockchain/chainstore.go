package blockchain

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	common2 "github.com/elastos/Elastos.ELA.Elephant.Node/common"
	types2 "github.com/elastos/Elastos.ELA.Elephant.Node/id/types"
	"github.com/elastos/Elastos.ELA.SideChain.ID/blockchain"
	blockchain2 "github.com/elastos/Elastos.ELA.SideChain/blockchain"
	"github.com/elastos/Elastos.ELA.SideChain/database"
	"github.com/elastos/Elastos.ELA.SideChain/types"
	"github.com/elastos/Elastos.ELA/common"
	"github.com/elastos/Elastos.ELA/common/log"
	"github.com/elastos/Elastos.ELA/crypto"
	"sort"
)

const DID_PropertyPrefix blockchain2.EntryPrefix = 0x99

type IDChainStoreEx struct {
	*blockchain.IDChainStore
}

func NewChainStoreEx(idChainStore *blockchain.IDChainStore) (*IDChainStoreEx, error) {
	store := &IDChainStoreEx{
		idChainStore,
	}
	store.RegisterFunctions(true, blockchain2.StoreFuncNames.PersistCurrentBlock, store.persistCurrentBlock)
	store.RegisterFunctions(false, blockchain2.StoreFuncNames.RollbackCurrentBlock, store.rollbackCurrentBlock)
	return store, nil
}

// key: SYS_CurrentBlock
// value: current block hash || height
func (s *IDChainStoreEx) persistCurrentBlock(batch database.Batch, b *types.Block) error {
	key := new(bytes.Buffer)
	key.WriteByte(byte(blockchain2.SYS_CurrentBlock))

	value := new(bytes.Buffer)
	blockHash := b.Hash()
	if err := blockHash.Serialize(value); err != nil {
		return err
	}
	if err := common.WriteUint32(value, b.Header.Height); err != nil {
		return err
	}
	batch.Put(key.Bytes(), value.Bytes())
	go s.externalBlockAction(b)
	return nil
}

func (s *IDChainStoreEx) rollbackCurrentBlock(batch database.Batch, b *types.Block) error {
	key := new(bytes.Buffer)
	key.WriteByte(byte(blockchain2.SYS_CurrentBlock))

	value := bytes.NewBuffer(nil)
	previous := b.Header.Previous
	if err := previous.Serialize(value); err != nil {
		return err
	}
	if err := common.WriteUint32(value, b.Header.Height-1); err != nil {
		return err
	}
	batch.Put(key.Bytes(), value.Bytes())
	//TODO rollback
	return nil
}

func (s *IDChainStoreEx) externalBlockAction(b *types.Block) {

	var didPropertys []types2.DidProperty
	for _, tx := range b.Transactions {
		if len(tx.Attributes) > 0 {
			if types.Memo == tx.Attributes[0].Usage {
				didMemo := types2.DidMemo{}
				err := json.Unmarshal(tx.Attributes[0].Data, &didMemo)
				if err != nil {
					log.Warn("[parsing did property]: Not a valid property")
					continue
				}

				if len(didMemo.Msg) == 0 || len(didMemo.Pub) == 0 || len(didMemo.Sig) == 0 {
					log.Warn("[parsing did property]: invalid 'msg' or 'pub' or 'sig' key in memo")
					continue
				}

				pub, err := hex.DecodeString(didMemo.Pub)
				if err != nil {
					log.Warn("[parsing did property]: invalid memo pub")
					continue
				}
				publicKey, err := crypto.DecodePoint(pub)
				if err != nil {
					log.Warn("[parsing did property]: invalid memo public key")
					continue
				}
				msg, err := hex.DecodeString(didMemo.Msg)
				if err != nil {
					log.Warn("[parsing did property]: invalid memo msg")
					continue
				}
				sig, err := hex.DecodeString(didMemo.Sig)
				if err != nil {
					log.Warn("[parsing did property]: invalid memo sig")
					continue
				}
				err = crypto.Verify(*publicKey, msg, sig)
				if err != nil {
					log.Warn("[parsing did property]: verify Error")
					continue
				}
				raw := types2.DidInfo{}
				err = json.Unmarshal(msg, &raw)
				if err != nil {
					log.Warn("[parsing did property]: RawData is not Json")
					continue
				}
				did, err := common2.GenDid(pub)
				for _, v := range raw.Properties {
					didPropertys = append(didPropertys, types2.DidProperty{
						Did:                 []byte(did),
						Did_status:          []byte(v.Status),
						Public_key:          pub,
						Property_key:        []byte(v.Key),
						Property_key_status: []byte(v.Status),
						Property_value:      []byte(v.Value),
						Txid:                tx.Hash(),
						Block_time:          b.Timestamp,
						Height:              b.Height,
					})
				}
			}
		}
	}
	if len(didPropertys) > 0 {
		batch := s.NewBatch()
		for _, v := range didPropertys {
			err := persistDidProperty(batch, v, b)
			if err != nil {
				log.Warn("[parsing did property]: Unexpected happend , rollback database")
				batch.Rollback()
				return
			}
		}
		batch.Commit()
	}
}

// key: DID || property_key || height
// value: property
func persistDidProperty(batch database.Batch, property types2.DidProperty, b *types.Block) error {
	key := new(bytes.Buffer)
	key.WriteByte(byte(DID_PropertyPrefix))
	common.WriteVarBytes(key, property.Did)
	common.WriteVarBytes(key, property.Property_key)
	common.WriteUint32(key, property.Height)
	value := new(bytes.Buffer)
	if err := property.Serialize(value); err != nil {
		return err
	}
	batch.Put(key.Bytes(), value.Bytes())
	return nil
}

func (c *IDChainStoreEx) GetDidPropertyByKey(encodeDid string, propertyKey string) types2.DidPropertyDisplaySorter {
	key := new(bytes.Buffer)
	key.WriteByte(byte(DID_PropertyPrefix))
	common.WriteVarBytes(key, []byte(encodeDid))
	if propertyKey != "" {
		common.WriteVarBytes(key, []byte(propertyKey))
	}
	it := c.NewIterator([]byte(key.Bytes()))
	var dpds types2.DidPropertyDisplaySorter
	for it.Next() {
		val := new(bytes.Buffer)
		val.Write(it.Value())
		dp := types2.DidProperty{}
		dpd, _ := dp.Deserialize(val)
		dpds = append(dpds, *dpd)
	}
	it.Release()
	sort.Sort(dpds)
	return dpds
}

func (c *IDChainStoreEx) GetDidProperty(hexDid string) types2.DidPropertyDisplaySorter {
	return c.GetDidPropertyByKey(hexDid, "")
}
