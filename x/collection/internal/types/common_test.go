package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	defaultName        = "name"
	defaultContractID  = "abcdef01"
	defaultBaseImgURI  = "base-img-uri"
	defaultMeta        = "{}"
	defaultDecimals    = 6
	defaultAmount      = 1000
	defaultTokenType   = "10000001"
	defaultTokenIndex  = "00000001"
	defaultTokenID1    = defaultTokenType + defaultTokenIndex
	defaultTokenID2    = defaultTokenType + "00000002"
	defaultTokenTypeFT = "00000001"
	defaultTokenIDFT   = defaultTokenTypeFT + "00000000"
)

var (
	addr1 = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	addr2 = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
)
