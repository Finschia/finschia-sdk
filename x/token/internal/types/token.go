package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Tokens []Token

func (ts Tokens) String() string {
	b, err := json.Marshal(ts)
	if err != nil {
		panic(err)
	}
	return string(b)
}

type Token interface {
	GetName() string
	GetSymbol() string
	GetTokenURI() string
	SetTokenURI(tokenURI string)
	GetMintable() bool
	GetDecimals() sdk.Int
	String() string
}

var _ Token = (*BaseToken)(nil)

type BaseToken struct {
	Name     string  `json:"name"`
	Symbol   string  `json:"symbol"`
	TokenURI string  `json:"token_uri"`
	Decimals sdk.Int `json:"decimals"`
	Mintable bool    `json:"mintable"`
}

func NewToken(name, symbol, tokenURI string, decimals sdk.Int, mintable bool) Token {
	return &BaseToken{
		Name:     name,
		Symbol:   symbol,
		TokenURI: tokenURI,
		Decimals: decimals,
		Mintable: mintable,
	}
}
func (t BaseToken) GetName() string      { return t.Name }
func (t BaseToken) GetSymbol() string    { return t.Symbol }
func (t BaseToken) GetTokenURI() string  { return t.TokenURI }
func (t BaseToken) GetMintable() bool    { return t.Mintable }
func (t BaseToken) GetDecimals() sdk.Int { return t.Decimals }
func (t *BaseToken) SetTokenURI(tokenURI string) {
	t.TokenURI = tokenURI
}
func (t BaseToken) String() string {
	b, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(b)
}
