package common

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"github.com/elastos/Elastos.ELA.SideChain.ID/pact"
	"github.com/elastos/Elastos.ELA/common"
	"github.com/elastos/Elastos.ELA/core/contract"
	"github.com/elastos/Elastos.ELA/crypto"
	"math/big"
)

func GetProgramHash(public []byte) (*common.Uint168, error) {
	hash, err := contract.PublicKeyToStandardProgramHash(public)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func GetAddress(public []byte) (string, error) {
	hash, err := GetProgramHash(public)
	if err != nil {
		return "", err
	}
	return hash.ToAddress()
}

//GetAddressFromPrivKey get address using privateKey
func GetPublicKeyFromPrivKey(privKey string) (*crypto.PublicKey, error) {
	priv := new(ecdsa.PrivateKey)
	c := elliptic.P256()
	priv.PublicKey.Curve = c
	k := new(big.Int)
	privKeyBytes, err := hex.DecodeString(privKey)
	if err != nil {
		return nil, err
	}
	k.SetBytes(privKeyBytes)
	priv.D = k
	priv.PublicKey.X, priv.PublicKey.Y = c.ScalarBaseMult(k.Bytes())
	publicKey := new(crypto.PublicKey)
	publicKey.X = new(big.Int).Set(priv.PublicKey.X)
	publicKey.Y = new(big.Int).Set(priv.PublicKey.Y)
	return publicKey, nil
}

func GetAddressFromPrivKey(privKey string) (string, error) {
	pub, err := GetPublicKeyFromPrivKey(privKey)
	if err != nil {
		return "", err
	}
	return getAddressFromPublicKey(pub)
}

func GetDIDFromPrivKey(privKey string) (string, error) {
	pub, err := GetPublicKeyFromPrivKey(privKey)
	if err != nil {
		return "", err
	}
	return getDIDFromPublicKey(pub)
}

//GetAddressFromPublicKey get address from public key
func getAddressFromPublicKey(publicKey *crypto.PublicKey) (string, error) {
	ct, err := contract.CreateStandardContract(publicKey)
	if err != nil {
		return "", err
	}
	programHash := ct.ToProgramHash()
	addr, _ := programHash.ToAddress()
	return addr, nil
}

//GetAddressFromPublicKey get did from public key
func getDIDFromPublicKey(publicKey *crypto.PublicKey) (string, error) {
	content, err := publicKey.EncodePoint(true)
	if err != nil {
		return "", err
	}
	redeemScript := new(bytes.Buffer)
	redeemScript.WriteByte(byte(len(content)))
	redeemScript.Write(content)
	redeemScript.WriteByte(0xAD)
	if err != nil {
		return "", err
	}
	uint168 := common.ToProgramHash(pact.PrefixRegisterId, redeemScript.Bytes())
	if err != nil {
		return "", err
	}
	did, err := uint168.ToAddress()
	if err != nil {
		return "", err
	}
	return did, nil
}
