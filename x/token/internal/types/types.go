package types

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Token struct {
	Name     string  `json:"name"`
	Symbol   string  `json:"symbol"`
	Mintable bool    `json:"mintable"`
	Decimals sdk.Int `json:"decimals"`
	TokenURI string  `json:"token_uri"`
}

func NewFT(name, symbol string, decimals sdk.Int, mintable bool) Token {
	return Token{
		Name:     name,
		Symbol:   symbol,
		Decimals: decimals,
		Mintable: mintable,
	}
}

func NewNFT(name, symbol, tokenURI string) Token {
	return Token{
		Name:     name,
		Symbol:   symbol,
		Decimals: sdk.NewInt(0),
		Mintable: false,
		TokenURI: tokenURI,
	}
}

func (t Token) GetTokenURI() string { return t.TokenURI }

func (t *Token) EditMetadata(tokenURI string) {
	t.TokenURI = tokenURI
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

type Collection struct {
	Symbol string `json:"symbol"`
}

func NewCollection(symbol string) Collection {
	return Collection{
		Symbol: symbol,
	}
}

func (c Collection) String() string {
	return string(codec.MustMarshalJSONIndent(ModuleCdc, c))
}

type Collections []Collection

func (collections Collections) String() string {
	var collectionStrings []string
	for _, t := range collections {
		collectionStrings = append(collectionStrings, t.String())
	}
	return strings.Join(collectionStrings, ",")
}
