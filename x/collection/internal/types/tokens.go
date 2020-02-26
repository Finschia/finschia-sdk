package types

import (
	"encoding/json"

	linktype "github.com/line/link/types"
)

const (
	TokenTypeLength    = 4
	SmallestAlphanum   = "0"
	LargestAlphanum    = "z"
	TokenIDLength      = linktype.TokenIDLen
	FungibleFlag       = SmallestAlphanum
	ReservedEmpty      = "0000"
	SmallestFTType     = "0001"
	SmallestNFTType    = "1001"
	SmallestTokenIndex = "0001"
)

type Tokens []Token

func (ts Tokens) String() string {
	b, err := json.Marshal(ts)
	if err != nil {
		panic(err)
	}
	return string(b)
}
