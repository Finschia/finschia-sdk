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
	defaultTokenType = "10000001"
	defaultTokenID1  = defaultTokenType + "00000001"
	defaultTokenID2  = defaultTokenType + "00000002"
	defaultTokenIDFT = "0000000100000000"
)

var (
	addr1 = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	addr2 = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
)
