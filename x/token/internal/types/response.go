package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

type TokenWithSupply struct {
	Token
	Supply sdk.Int `json:"supply"`
}

func (t TokenWithSupply) String() string {
	return string(codec.MustMarshalJSONIndent(ModuleCdc, t))
}

type TokensWithSupply []TokenWithSupply

func (tokens TokensWithSupply) String() string {
	var tokenStrings []string
	for _, t := range tokens {
		tokenStrings = append(tokenStrings, t.String())
	}
	return strings.Join(tokenStrings, ",")
}

type CollectionWithTokens struct {
	Collection `json:"collection"`
	Tokens     `json:"tokens"`
}

func (c CollectionWithTokens) String() string {
	return string(codec.MustMarshalJSONIndent(ModuleCdc, c))
}

type CollectionsWithTokens []CollectionWithTokens

func (cs CollectionsWithTokens) String() string {
	var collectionStrings []string
	for _, c := range cs {
		collectionStrings = append(collectionStrings, c.String())
	}
	return strings.Join(collectionStrings, ",")
}
