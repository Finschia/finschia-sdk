package wallet

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	DefaultBIP39Passphrase = ""
)

type HDWallet struct {
	masterPrivateKey [32]byte
	masterChainCode  [32]byte
	coinType         uint32
}

/*
	crypto.keys.dbKeybase encrypt and armor the private key when generating it.
	So it is too slow to generate a large number of keys.
	Therefore, defined a new class for key generation.
*/
func NewHDWallet(mnemonic string) (*HDWallet, error) {
	coinType := sdk.GetConfig().GetCoinType()

	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, DefaultBIP39Passphrase)
	if err != nil {
		return nil, err
	}
	masterPriv, ch := hd.ComputeMastersFromSeed(seed)

	return &HDWallet{
		masterPrivateKey: masterPriv,
		masterChainCode:  ch,
		coinType:         coinType,
	}, nil
}

func (hdw HDWallet) GeneratePrivateKey(accountNumber uint32, index uint32) (priv secp256k1.PrivKeySecp256k1,
	err error) {
	hdPath := hd.NewFundraiserParams(accountNumber, hdw.coinType, index)
	priv, err = hd.DerivePrivateKeyForPath(hdw.masterPrivateKey, hdw.masterChainCode, hdPath.String())
	return
}

func (hdw HDWallet) GetKeyWallet(accountNumber uint32, index uint32) (keyWallet *KeyWallet, err error) {
	priv, err := hdw.GeneratePrivateKey(accountNumber, index)
	keyWallet = NewKeyWallet(priv)
	return
}
