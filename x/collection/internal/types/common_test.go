package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	defaultName      = "name"
	defaultSymbol    = "linktkn"
	defaultTokenURI  = "token-uri"
	defaultDecimals  = 6
	defaultTokenType = "a000"
)

var (
	defaultAddr = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
)
