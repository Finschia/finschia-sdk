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
	SetName(name string)
	GetSymbol() string
	GetMeta() string
	SetMeta(meta string)
	GetImageURI() string
	SetImageURI(tokenURI string)
	GetMintable() bool
	GetDecimals() sdk.Int
	String() string
}

var _ Token = (*BaseToken)(nil)

type BaseToken struct {
	ContractID string  `json:"contract_id"`
	Name       string  `json:"name"`
	Symbol     string  `json:"symbol"`
	Meta       string  `json:"meta"`
	ImageURI   string  `json:"img_uri"`
	Decimals   sdk.Int `json:"decimals"`
	Mintable   bool    `json:"mintable"`
}

func NewToken(contractID, name, symbol, meta string, imageURI string, decimals sdk.Int, mintable bool) Token {
	return &BaseToken{
		ContractID: contractID,
		Name:       name,
		Symbol:     symbol,
		Meta:       meta,
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
func (t *BaseToken) SetName(name string) {
	t.Name = name
}
func (t BaseToken) GetMeta() string { return t.Meta }
func (t *BaseToken) SetMeta(meta string) {
	t.Meta = meta
}
func (t *BaseToken) SetImageURI(tokenURI string) {
	t.ImageURI = tokenURI
}

func (t BaseToken) String() string {
	b, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(b)
}
