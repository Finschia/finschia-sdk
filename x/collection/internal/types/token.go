package types

import (
	"encoding/json"

	sdk "github.com/line/lbm-sdk/types"
)

type Token interface {
	GetName() string
	SetName(name string)
	GetContractID() string
	GetTokenID() string
	GetTokenType() string
	GetTokenIndex() string
	String() string
	GetMeta() string
	SetMeta(meta string)
}

type FT interface {
	Token
	GetMintable() bool
	GetDecimals() sdk.Int
}

type NFT interface {
	Token
	GetOwner() sdk.AccAddress
	SetOwner(sdk.AccAddress)
}

var _ Token = (*BaseNFT)(nil)

type BaseNFT struct {
	ContractID string         `json:"contract_id"`
	TokenID    string         `json:"token_id"`
	Owner      sdk.AccAddress `json:"owner"`
	Name       string         `json:"name"`
	Meta       string         `json:"meta"`
}

func NewNFT(contractID, tokenID, name, meta string, owner sdk.AccAddress) NFT {
	return &BaseNFT{
		ContractID: contractID,
		TokenID:    tokenID,
		Owner:      owner,
		Name:       name,
		Meta:       meta,
	}
}
func (t BaseNFT) GetName() string          { return t.Name }
func (t BaseNFT) GetContractID() string    { return t.ContractID }
func (t BaseNFT) GetOwner() sdk.AccAddress { return t.Owner }
func (t BaseNFT) GetTokenID() string       { return t.TokenID }
func (t BaseNFT) GetTokenType() string     { return t.TokenID[:TokenTypeLength] }
func (t BaseNFT) GetTokenIndex() string    { return t.TokenID[TokenTypeLength:] }
func (t *BaseNFT) SetName(name string) {
	t.Name = name
}
func (t *BaseNFT) SetOwner(owner sdk.AccAddress) {
	t.Owner = owner
}
func (t BaseNFT) String() string {
	b, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(b)
}
func (t BaseNFT) GetMeta() string { return t.Meta }
func (t *BaseNFT) SetMeta(meta string) {
	t.Meta = meta
}

var _ Token = (*BaseFT)(nil)
var _ FT = (*BaseFT)(nil)

type BaseFT struct {
	ContractID string  `json:"contract_id"`
	TokenID    string  `json:"token_id"`
	Decimals   sdk.Int `json:"decimals"`
	Mintable   bool    `json:"mintable"`
	Name       string  `json:"name"`
	Meta       string  `json:"meta"`
}

func NewFT(contractID, tokenID, name, meta string, decimals sdk.Int, mintable bool) FT {
	return &BaseFT{
		ContractID: contractID,
		TokenID:    tokenID,
		Decimals:   decimals,
		Mintable:   mintable,
		Name:       name,
		Meta:       meta,
	}
}
func (t BaseFT) GetName() string       { return t.Name }
func (t BaseFT) GetContractID() string { return t.ContractID }
func (t BaseFT) GetMintable() bool     { return t.Mintable }
func (t BaseFT) GetDecimals() sdk.Int  { return t.Decimals }
func (t BaseFT) GetTokenID() string    { return t.TokenID }
func (t BaseFT) GetTokenType() string  { return t.TokenID[:TokenTypeLength] }
func (t BaseFT) GetTokenIndex() string { return t.TokenID[TokenTypeLength:] }
func (t *BaseFT) SetName(name string) {
	t.Name = name
}
func (t BaseFT) String() string {
	b, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(b)
}
func (t BaseFT) GetMeta() string { return t.Meta }
func (t *BaseFT) SetMeta(meta string) {
	t.Meta = meta
}

type Tokens []Token

func (ts Tokens) String() string {
	b, err := json.Marshal(ts)
	if err != nil {
		panic(err)
	}
	return string(b)
}
