package types

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Token struct {
	Name     string         `json:"name"`
	Symbol   string         `json:"symbol"`
	Owner    sdk.AccAddress `json:"owner"`
	Mintable bool           `json:"mintable"`
}

func (t Token) String() string {
	return string(codec.MustMarshalJSONIndent(ModuleCdc, t))
}

type Tokens []Token

func (tokens Tokens) String() string {
	var tokenStrings []string
	for _, t := range tokens {
		tokenStrings = append(tokenStrings, t.String())
	}
	return strings.Join(tokenStrings, ",")
}
