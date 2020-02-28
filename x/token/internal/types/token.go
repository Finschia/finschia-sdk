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
	GetContractID() string
	GetName() string
	SetName(name string) Token
	GetSymbol() string
	GetImageURI() string
	SetImageURI(tokenURI string) Token
	GetMintable() bool
	GetDecimals() sdk.Int
	String() string
}

var _ Token = (*BaseToken)(nil)

type BaseToken struct {
	ContractID string  `json:"contract_id"`
	Name       string  `json:"name"`
	Symbol     string  `json:"symbol"`
	ImageURI   string  `json:"image_uri"`
	Decimals   sdk.Int `json:"decimals"`
	Mintable   bool    `json:"mintable"`
}

func NewToken(contractID, name, symbol, imageURI string, decimals sdk.Int, mintable bool) Token {
	return &BaseToken{
		ContractID: contractID,
		Name:       name,
		Symbol:     symbol,
		ImageURI:   imageURI,
		Decimals:   decimals,
		Mintable:   mintable,
	}
}

func (t BaseToken) GetContractID() string { return t.ContractID }
func (t BaseToken) GetName() string       { return t.Name }
func (t BaseToken) GetSymbol() string     { return t.Symbol }
func (t BaseToken) GetImageURI() string   { return t.ImageURI }
func (t BaseToken) GetMintable() bool     { return t.Mintable }
func (t BaseToken) GetDecimals() sdk.Int  { return t.Decimals }
func (t *BaseToken) SetName(name string) Token {
	t.Name = name
	return t
}
func (t *BaseToken) SetImageURI(tokenURI string) Token {
	t.ImageURI = tokenURI
	return t
}
func (t BaseToken) String() string {
	b, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(b)
}
