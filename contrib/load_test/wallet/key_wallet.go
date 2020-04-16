package wallet

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type KeyWallet struct {
	privateKey secp256k1.PrivKeySecp256k1
	address    sdk.AccAddress
}

func NewKeyWallet(privateKey secp256k1.PrivKeySecp256k1) *KeyWallet {
	return &KeyWallet{
		privateKey: privateKey,
		address:    privateKey.PubKey().Address().Bytes(),
	}
}

func (k *KeyWallet) PrivateKey() secp256k1.PrivKeySecp256k1 {
	return k.privateKey
}

func (k *KeyWallet) Address() sdk.AccAddress {
	return k.address
}
