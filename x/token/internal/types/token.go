package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Token interface {
	GetName() string
	GetSymbol() string
	GetDenom() string
	GetTokenID() string
	GetTokenURI() string
	SetTokenURI(tokenURI string)
	String() string
}

type FT interface {
	Token
	GetMintable() bool
	GetDecimals() sdk.Int
}

var _ Token = (*BaseFT)(nil)
var _ FT = (*BaseFT)(nil)

type BaseFT struct {
	Name     string  `json:"name"`
	Symbol   string  `json:"symbol"`
	TokenURI string  `json:"token_uri"`
	Decimals sdk.Int `json:"decimals"`
	Mintable bool    `json:"mintable"`
}

func NewFT(name, symbol, tokenURI string, decimals sdk.Int, mintable bool) FT {
	return &BaseFT{
		Name:     name,
		Symbol:   symbol,
		TokenURI: tokenURI,
		Decimals: decimals,
		Mintable: mintable,
	}
}
func (t BaseFT) GetName() string      { return t.Name }
func (t BaseFT) GetSymbol() string    { return t.Symbol }
func (t BaseFT) GetTokenURI() string  { return t.TokenURI }
func (t BaseFT) GetDenom() string     { return t.Symbol }
func (t BaseFT) GetMintable() bool    { return t.Mintable }
func (t BaseFT) GetDecimals() sdk.Int { return t.Decimals }
func (t BaseFT) GetTokenID() string   { return "" }
func (t *BaseFT) SetTokenURI(tokenURI string) {
	t.TokenURI = tokenURI
}
func (t BaseFT) String() string {
	b, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(b)
}
