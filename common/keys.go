package common

import (
	"bytes"
	"errors"
	"github.com/elastos/Elastos.ELA.SideChain.ID/pact"
	"github.com/elastos/Elastos.ELA/common"
	"github.com/elastos/Elastos.ELA/crypto"
)

func GenDid(pub []byte) (string, error) {
	publicKey, err := crypto.DecodePoint(pub)
	if err != nil {
		return "", err
	}
	redeemScript, err := createRegistedRedeemedScript(publicKey)
	if err != nil {
		return "", err
	}
	uint168 := common.ToProgramHash(pact.PrefixRegisterId, redeemScript)
	if err != nil {
		return "", err
	}
	did, err := uint168.ToAddress()
	if err != nil {
		return "", err
	}
	return did, nil
}

func createRegistedRedeemedScript(publicKey *crypto.PublicKey) ([]byte, error) {
	content, err := publicKey.EncodePoint(true)
	if err != nil {
		return nil, errors.New("create standard redeem script, encode public key failed")
	}
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(len(content)))
	buf.Write(content)
	buf.WriteByte(0xAD)

	return buf.Bytes(), nil
}
