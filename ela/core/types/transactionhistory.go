package types

import (
	"github.com/elastos/Elastos.ELA/common"
	"github.com/pkg/errors"
	"io"
)

type TransactionHistory struct {
	Address    string
	Txid       string
	Type       string
	Value      uint64
	CreateTime uint64
	Height     uint64
	Fee        uint64
	Inputs     []string
	Outputs    []string
	TxType     string
	Memo       string
}

func (th *TransactionHistory) Serialize(w io.Writer) error {
	err := common.WriteVarString(w, th.Address)
	if err != nil {
		return errors.New("[TransactionHistory], Address serialize failed.")
	}
	err = common.WriteVarString(w, th.Txid)
	if err != nil {
		return errors.New("[TransactionHistory], Txid serialize failed.")
	}
	err = common.WriteVarString(w, th.Type)
	if err != nil {
		return errors.New("[TransactionHistory], Type serialize failed.")
	}
	err = common.WriteUint64(w, th.Value)
	if err != nil {
		return errors.New("[TransactionHistory], Value serialize failed.")
	}
	err = common.WriteUint64(w, th.CreateTime)
	if err != nil {
		return errors.New("[TransactionHistory], CreateTime serialize failed.")
	}
	err = common.WriteUint64(w, th.Height)
	if err != nil {
		return errors.New("[TransactionHistory], Height serialize failed.")
	}
	err = common.WriteUint64(w, th.Fee)
	if err != nil {
		return errors.New("[TransactionHistory], Fee serialize failed.")
	}
	err = common.WriteVarUint(w, uint64(len(th.Inputs)))
	if err != nil {
		return errors.New("[TransactionHistory], Length of inputs serialize failed.")
	}
	for i := 0; i < len(th.Inputs); i++ {
		err = common.WriteVarString(w, th.Inputs[i])
		if err != nil {
			return errors.New("[TransactionHistory], input:" + th.Inputs[i] + " serialize failed.")
		}
	}
	err = common.WriteVarUint(w, uint64(len(th.Outputs)))
	if err != nil {
		return errors.New("[TransactionHistory], Length of outputs serialize failed.")
	}
	for i := 0; i < len(th.Outputs); i++ {
		err = common.WriteVarString(w, th.Outputs[i])
		if err != nil {
			return errors.New("[TransactionHistory], output:" + th.Outputs[i] + " serialize failed.")
		}
	}
	err = common.WriteVarString(w, th.TxType)
	if err != nil {
		return errors.New("[TransactionHistory], TxType serialize failed.")
	}
	err = common.WriteVarString(w, th.Memo)
	if err != nil {
		return errors.New("[TransactionHistory], Memo serialize failed.")
	}
	return nil
}

func (th *TransactionHistory) Deserialize(r io.Reader) error {
	var err error
	th.Address, err = common.ReadVarString(r)
	if err != nil {
		return errors.New("[TransactionHistory], Address deserialize failed.")
	}
	th.Txid, err = common.ReadVarString(r)
	if err != nil {
		return errors.New("[TransactionHistory], Txid deserialize failed.")
	}
	th.Type, err = common.ReadVarString(r)
	if err != nil {
		return errors.New("[TransactionHistory], Type deserialize failed.")
	}
	th.Value, err = common.ReadUint64(r)
	if err != nil {
		return errors.New("[TransactionHistory], Value deserialize failed.")
	}
	th.CreateTime, err = common.ReadUint64(r)
	if err != nil {
		return errors.New("[TransactionHistory], CreateTime deserialize failed.")
	}
	th.Height, err = common.ReadUint64(r)
	if err != nil {
		return errors.New("[TransactionHistory], Height deserialize failed.")
	}
	th.Fee, err = common.ReadUint64(r)
	if err != nil {
		return errors.New("[TransactionHistory], Fee deserialize failed.")
	}
	n, err := common.ReadVarUint(r, 0)
	if err != nil {
		return errors.New("[TransactionHistory], length of inputs deserialize failed.")
	}
	for i := uint64(0); i < n; i++ {
		str, err := common.ReadVarString(r)
		if err != nil {
			return errors.New("[TransactionHistory], input:" + str + " deserialize failed.")
		}
		th.Inputs = append(th.Inputs, str)
	}
	n, err = common.ReadVarUint(r, 0)
	if err != nil {
		return errors.New("[TransactionHistory], length of outputs deserialize failed.")
	}
	for i := uint64(0); i < n; i++ {
		str, err := common.ReadVarString(r)
		if err != nil {
			return errors.New("[TransactionHistory], output:" + str + " deserialize failed.")
		}
		th.Outputs = append(th.Outputs, str)
	}
	th.TxType , err = common.ReadVarString(r)
	if err != nil {
		return errors.New("[TransactionHistory], TxType serialize failed.")
	}
	th.Memo , err = common.ReadVarString(r)
	if err != nil {
		return errors.New("[TransactionHistory], Memo serialize failed.")
	}
	return nil
}
