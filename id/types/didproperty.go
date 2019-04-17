package types

import (
	"encoding/hex"
	"errors"
	common2 "github.com/elastos/Elastos.ELA.Elephant.Node/common"
	"github.com/elastos/Elastos.ELA/common"
	"io"
)

type DidPropertyDisplay struct {
	Did                 string
	Did_status          string
	Public_key          string
	Property_key        string
	Property_key_status string
	Property_value      string
	Txid                string
	Block_time          uint32
	Height              uint32
}

type DidProperty struct {
	Did                 []byte
	Did_status          []byte
	Public_key          []byte
	Property_key        []byte
	Property_key_status []byte
	Property_value      []byte
	Txid                common.Uint256
	Block_time          uint32
	Height              uint32
}

type DidMemo struct {
	Msg string `json:"msg"`
	Pub string `json:"pub"`
	Sig string `json:"sig"`
}

type Properties struct {
	Key    string
	Value  string
	Status string
}

type DidInfo struct {
	Tag        string
	Ver        string
	Status     string
	Properties []Properties
}

func (dp *DidProperty) Serialize(w io.Writer) error {
	err := common.WriteVarBytes(w, dp.Did)
	if err != nil {
		return errors.New("[DidProperty], did serialize failed.")
	}
	common.WriteVarBytes(w, dp.Did_status)
	if err != nil {
		return errors.New("[DidProperty], did_status serialize failed.")
	}
	err = common.WriteVarBytes(w, dp.Public_key)
	if err != nil {
		return errors.New("[DidProperty], Public_key serialize failed.")
	}
	err = common.WriteVarBytes(w, dp.Property_key)
	if err != nil {
		return errors.New("[DidProperty], Property_key serialize failed.")
	}
	err = common.WriteVarBytes(w, dp.Property_key_status)
	if err != nil {
		return errors.New("[DidProperty], property_key_status serialize failed.")
	}
	err = common.WriteVarBytes(w, dp.Property_value)
	if err != nil {
		return errors.New("[DidProperty], Property_value serialize failed.")
	}
	err = common.WriteVarBytes(w, dp.Txid.Bytes())
	if err != nil {
		return errors.New("[DidProperty], Txid serialize failed.")
	}
	err = common.WriteUint32(w, dp.Block_time)
	if err != nil {
		return errors.New("[DidProperty], Block_time serialize failed.")
	}
	err = common.WriteUint32(w, dp.Height)
	if err != nil {
		return errors.New("[DidProperty], Height serialize failed.")
	}
	return nil
}

func (dp *DidProperty) Deserialize(r io.Reader) (*DidPropertyDisplay, error) {
	rst := new(DidPropertyDisplay)
	did, err := common.ReadVarBytes(r, 1024, "did")
	if err != nil {
		return rst, errors.New("[DidProperty], did deserialize failed.")
	}
	rst.Did, err = common2.GenDid(did)
	if err != nil {
		return rst, errors.New("[DidProperty], did deserialize failed.")
	}
	Did_status, err := common.ReadVarBytes(r, 1024, "did status")
	if err != nil {
		return rst, errors.New("[DidProperty], did_status deserialize failed.")
	}
	rst.Did_status = string(Did_status)
	Public_key, err := common.ReadVarBytes(r, 1024, "public key")
	if err != nil {
		return rst, errors.New("[DidProperty], Public_key deserialize failed.")
	}
	rst.Public_key = string(Public_key)
	Property_key, err := common.ReadVarBytes(r, common.MaxVarStringLength, "property key")
	if err != nil {
		return rst, errors.New("[DidProperty], Property_key deserialize failed.")
	}
	rst.Property_key = string(Property_key)
	Property_key_status, err := common.ReadVarBytes(r, 1024, "property_key_status")
	if err != nil {
		return rst, errors.New("[DidProperty], property_key_status deserialize failed.")
	}
	rst.Property_key_status = string(Property_key_status)
	Property_value, err := common.ReadVarBytes(r, common.MaxVarStringLength, "Property_value")
	if err != nil {
		return rst, errors.New("[DidProperty], Property_value deserialize failed.")
	}
	rst.Property_value = string(Property_value)
	Txid, err := common.ReadVarBytes(r, 1024, "txid")
	if err != nil {
		return rst, errors.New("[DidProperty], Txid deserialize failed.")
	}
	rst.Txid, err = common2.ReverseHexString(hex.EncodeToString(Txid))
	if err != nil {
		return rst, errors.New("[DidProperty], Txid deserialize failed.")
	}
	Block_time, err := common.ReadUint32(r)
	if err != nil {
		return rst, errors.New("[DidProperty], Block_time deserialize failed.")
	}
	rst.Block_time = Block_time
	Height, err := common.ReadUint32(r)
	if err != nil {
		return rst, errors.New("[DidProperty], Height deserialize failed.")
	}
	rst.Height = Height
	return rst, nil
}
