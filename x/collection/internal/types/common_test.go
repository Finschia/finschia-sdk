package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	defaultName      = "name"
	defaultSymbol    = "token001"
	defaultTokenURI  = "token-uri"
	defaultDecimals  = 6
	defaultAmount    = 1000
	defaultTokenType = "1001"
	defaultTokenID1  = defaultTokenType + "0001"
	defaultTokenIDFT = "00010000"
)

var (
	addr1 = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	addr2 = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
)
