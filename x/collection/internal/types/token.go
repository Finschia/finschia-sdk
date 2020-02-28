package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Token interface {
	GetName() string
	SetName(name string) Token
	GetSymbol() string
	GetTokenID() string
	GetTokenType() string
	GetTokenIndex() string
	String() string
}

type FT interface {
	Token
	GetMintable() bool
	GetDecimals() sdk.Int
}

type NFT interface {
	Token
	GetOwner() sdk.AccAddress
	SetOwner(sdk.AccAddress) Token
}

var _ Token = (*BaseNFT)(nil)

type BaseNFT struct {
	Symbol  string         `json:"symbol"`
	TokenID string         `json:"token_id"`
	Owner   sdk.AccAddress `json:"owner"`
	Name    string         `json:"name"`
}

func NewNFT(symbol, tokenID, name string, owner sdk.AccAddress) NFT {
	return &BaseNFT{
		Symbol:  symbol,
		TokenID: tokenID,
		Owner:   owner,
		Name:    name,
	}
}
func (t BaseNFT) GetName() string          { return t.Name }
func (t BaseNFT) GetSymbol() string        { return t.Symbol }
func (t BaseNFT) GetOwner() sdk.AccAddress { return t.Owner }
func (t BaseNFT) GetTokenID() string       { return t.TokenID }
func (t BaseNFT) GetTokenType() string     { return t.TokenID[:TokenTypeLength] }
func (t BaseNFT) GetTokenIndex() string    { return t.TokenID[TokenTypeLength:] }
func (t *BaseNFT) SetName(name string) Token {
	t.Name = name
	return t
}
func (t *BaseNFT) SetOwner(owner sdk.AccAddress) Token {
	t.Owner = owner
	return t
}
func (t BaseNFT) String() string {
	b, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(b)
}

var _ Token = (*BaseFT)(nil)
var _ FT = (*BaseFT)(nil)

type BaseFT struct {
	Symbol   string  `json:"symbol"`
	TokenID  string  `json:"token_id"`
	Decimals sdk.Int `json:"decimals"`
	Mintable bool    `json:"mintable"`
	Name     string  `json:"name"`
}

func NewFT(symbol, tokenID, name string, decimals sdk.Int, mintable bool) FT {
	return &BaseFT{
		Symbol:   symbol,
		TokenID:  tokenID,
		Decimals: decimals,
		Mintable: mintable,
		Name:     name,
	}
}
func (t BaseFT) GetName() string       { return t.Name }
func (t BaseFT) GetSymbol() string     { return t.Symbol }
func (t BaseFT) GetMintable() bool     { return t.Mintable }
func (t BaseFT) GetDecimals() sdk.Int  { return t.Decimals }
func (t BaseFT) GetTokenID() string    { return t.TokenID }
func (t BaseFT) GetTokenType() string  { return t.TokenID[:TokenTypeLength] }
func (t BaseFT) GetTokenIndex() string { return t.TokenID[TokenTypeLength:] }
func (t *BaseFT) SetName(name string) Token {
	t.Name = name
	return t
}
func (t BaseFT) String() string {
	b, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(b)
}

type Tokens []Token

func (ts Tokens) String() string {
	b, err := json.Marshal(ts)
	if err != nil {
		panic(err)
	}
	return string(b)
}
