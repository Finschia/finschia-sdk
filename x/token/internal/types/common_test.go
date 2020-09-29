package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	defaultName       = "name"
	defaultContractID = "linktkn"
	defaultSymbol     = "BTC"
	defaultMeta       = "{}"
	defaultImageURI   = "image-uri"
	defaultDecimals   = 6
	defaultAmount     = 1000
)

var (
	addr1 = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	addr2 = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
)
