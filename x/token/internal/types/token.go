package types

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Token interface {
	GetName() string
	SetName(name string)
	GetSymbol() string
	GetDenom() string
	GetTokenID() string
	GetTokenURI() string
	SetTokenURI(tokenURI string)
	String() string
}

type Tokens []Token

func (tokens Tokens) String() string {
	b, err := json.Marshal(tokens)
	if err != nil {
		panic(err)
	}
	return string(b)
}

type FT interface {
	Token
	GetMintable() bool
	GetDecimals() sdk.Int
}

type NFT interface {
	Token
	GetOwner() sdk.AccAddress
}
type BaseToken struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	TokenURI string `json:"token_uri"`
}

func NewBaseToken(name, symbol, tokenURI string) *BaseToken {
	return &BaseToken{
		Name:     name,
		Symbol:   symbol,
		TokenURI: tokenURI,
	}
}
func (t BaseToken) GetName() string { return t.Name }
func (t *BaseToken) SetName(name string) {
	t.Name = name
}
func (t BaseToken) GetSymbol() string   { return t.Symbol }
func (t BaseToken) GetTokenURI() string { return t.TokenURI }
func (t *BaseToken) SetTokenURI(tokenURI string) {
	t.TokenURI = tokenURI
}

var _ FT = (*BaseFT)(nil)
var _ json.Marshaler = (*BaseFT)(nil)
var _ json.Unmarshaler = (*BaseFT)(nil)

type BaseFT struct {
	*BaseToken
	Decimals sdk.Int `json:"decimals"`
	Mintable bool    `json:"mintable"`
}

func NewFT(name, symbol, tokenURI string, decimals sdk.Int, mintable bool) *BaseFT {
	return NewBaseFT(NewBaseToken(name, symbol, tokenURI), decimals, mintable)
}
func NewBaseFT(baseToken *BaseToken, decimals sdk.Int, mintable bool) *BaseFT {
	return &BaseFT{
		BaseToken: baseToken,
		Decimals:  decimals,
		Mintable:  mintable,
	}
}
func (t BaseFT) GetDenom() string     { return t.Symbol }
func (t BaseFT) GetMintable() bool    { return t.Mintable }
func (t BaseFT) GetDecimals() sdk.Int { return t.Decimals }
func (t BaseFT) GetTokenID() string   { return "" }
func (t BaseFT) String() string {
	b, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(b)
}

var _ NFT = (*BaseNFT)(nil)
var _ json.Marshaler = (*BaseNFT)(nil)
var _ json.Unmarshaler = (*BaseFT)(nil)

type BaseNFT struct {
	*BaseToken
	Owner sdk.AccAddress `json:"owner"`
}

func NewNFT(name, symbol, tokenURI string, owner sdk.AccAddress) *BaseNFT {
	return NewBaseNFT(NewBaseToken(name, symbol, tokenURI), owner)
}
func NewBaseNFT(baseToken *BaseToken, owner sdk.AccAddress) *BaseNFT {
	return &BaseNFT{
		BaseToken: baseToken,
		Owner:     owner,
	}
}
func (t BaseNFT) GetDenom() string         { return t.Symbol }
func (t BaseNFT) GetOwner() sdk.AccAddress { return t.Owner }
func (t BaseNFT) GetTokenID() string       { return "" }

func (t BaseNFT) String() string {
	b, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(b)
}

var _ FT = (*BaseIDFT)(nil)
var _ json.Marshaler = (*BaseIDFT)(nil)
var _ json.Unmarshaler = (*BaseIDFT)(nil)

type BaseIDFT struct {
	*BaseFT
	TokenID string `json:"token_id"`
}

func NewIDFT(name, symbol, tokenURI string, decimals sdk.Int, mintable bool, tokenID string) *BaseIDFT {
	return NewBaseIDFTWithBaseFT(NewFT(name, symbol, tokenURI, decimals, mintable), tokenID)
}
func NewBaseIDFTWithBaseFT(baseFT *BaseFT, tokenID string) *BaseIDFT {
	return &BaseIDFT{
		BaseFT:  baseFT,
		TokenID: tokenID,
	}
}

func (t BaseIDFT) GetDenom() string   { return t.Symbol + t.TokenID }
func (t BaseIDFT) GetTokenID() string { return t.TokenID }
func (t BaseIDFT) String() string {
	b, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(b)
}

var _ NFT = (*BaseIDNFT)(nil)
var _ json.Marshaler = (*BaseIDNFT)(nil)
var _ json.Unmarshaler = (*BaseIDNFT)(nil)

type BaseIDNFT struct {
	*BaseNFT
	TokenID string `json:"token_id"`
}

func NewIDNFT(name, symbol, tokenURI string, owner sdk.AccAddress, tokenID string) *BaseIDNFT {
	return NewBaseIDNFTWithBaseNFT(NewNFT(name, symbol, tokenURI, owner), tokenID)
}
func NewBaseIDNFTWithBaseNFT(baseNFT *BaseNFT, tokenID string) *BaseIDNFT {
	return &BaseIDNFT{
		BaseNFT: baseNFT,
		TokenID: tokenID,
	}
}
func (t BaseIDNFT) GetDenom() string   { return t.Symbol + t.TokenID }
func (t BaseIDNFT) GetTokenID() string { return t.TokenID }
func (t BaseIDNFT) String() string {
	b, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(b)
}

//---------------------------------------------------------
// Custom json (un)marshaler to avoid go-amino, go-json bug
// see https://github.com/line/link/pull/354
//---------------------------------------------------------
func (t BaseFT) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name     string  `json:"name"`
		Symbol   string  `json:"symbol"`
		TokenURI string  `json:"token_uri"`
		Decimals sdk.Int `json:"decimals"`
		Mintable bool    `json:"mintable"`
	}{
		Name:     t.Name,
		Symbol:   t.Symbol,
		Decimals: t.Decimals,
		Mintable: t.Mintable,
	})
}
func (t *BaseFT) UnmarshalJSON(data []byte) error {
	type msgAlias *BaseFT
	return json.Unmarshal(data, msgAlias(t))
}

func (t BaseNFT) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name     string         `json:"name"`
		Symbol   string         `json:"symbol"`
		TokenURI string         `json:"token_uri"`
		Owner    sdk.AccAddress `json:"owner"`
	}{
		Name:     t.Name,
		Symbol:   t.Symbol,
		TokenURI: t.TokenURI,
		Owner:    t.Owner,
	})
}
func (t *BaseNFT) UnmarshalJSON(data []byte) error {
	type msgAlias *BaseNFT
	return json.Unmarshal(data, msgAlias(t))
}
func (t BaseIDFT) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name     string  `json:"name"`
		Symbol   string  `json:"symbol"`
		TokenURI string  `json:"token_uri"`
		Decimals sdk.Int `json:"decimals"`
		Mintable bool    `json:"mintable"`
		TokenID  string  `json:"token_id"`
	}{
		Name:     t.Name,
		Symbol:   t.Symbol,
		Decimals: t.Decimals,
		Mintable: t.Mintable,
		TokenID:  t.TokenID,
	})
}
func (t *BaseIDFT) UnmarshalJSON(data []byte) error {
	type msgAlias *BaseIDFT
	return json.Unmarshal(data, msgAlias(t))
}
func (t BaseIDNFT) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name     string         `json:"name"`
		Symbol   string         `json:"symbol"`
		TokenURI string         `json:"token_uri"`
		Owner    sdk.AccAddress `json:"owner"`
		TokenID  string         `json:"token_id"`
	}{
		Name:     t.Name,
		Symbol:   t.Symbol,
		TokenURI: t.TokenURI,
		Owner:    t.Owner,
		TokenID:  t.TokenID,
	})
}
func (t *BaseIDNFT) UnmarshalJSON(data []byte) error {
	type msgAlias *BaseIDNFT
	return json.Unmarshal(data, msgAlias(t))
}
