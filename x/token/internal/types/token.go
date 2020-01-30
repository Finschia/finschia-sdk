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

type NFT interface {
	Token
	GetOwner() sdk.AccAddress
	SetOwner(sdk.AccAddress)
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
func (t BaseToken) GetName() string     { return t.Name }
func (t BaseToken) GetSymbol() string   { return t.Symbol }
func (t BaseToken) GetTokenURI() string { return t.TokenURI }
func (t *BaseToken) SetTokenURI(tokenURI string) {
	t.TokenURI = tokenURI
}

var _ Token = (*BaseFT)(nil)
var _ FT = (*BaseFT)(nil)
var _ json.Marshaler = (*BaseFT)(nil)
var _ json.Unmarshaler = (*BaseFT)(nil)

type BaseFT struct {
	*BaseToken
	Decimals sdk.Int `json:"decimals"`
	Mintable bool    `json:"mintable"`
}

func NewFT(name, symbol, tokenURI string, decimals sdk.Int, mintable bool) FT {
	return NewBaseFT(NewBaseToken(name, symbol, tokenURI), decimals, mintable)
}
func NewBaseFT(baseToken *BaseToken, decimals sdk.Int, mintable bool) FT {
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

var _ Token = (*BaseNFT)(nil)
var _ NFT = (*BaseNFT)(nil)
var _ json.Marshaler = (*BaseNFT)(nil)
var _ json.Unmarshaler = (*BaseNFT)(nil)

type BaseNFT struct {
	*BaseToken
	Owner sdk.AccAddress `json:"owner"`
}

func NewNFT(name, symbol, tokenURI string, owner sdk.AccAddress) NFT {
	return NewBaseNFT(NewBaseToken(name, symbol, tokenURI), owner)
}
func NewBaseNFT(baseToken *BaseToken, owner sdk.AccAddress) NFT {
	return &BaseNFT{
		BaseToken: baseToken,
		Owner:     owner,
	}
}
func (t BaseNFT) GetDenom() string               { return t.Symbol }
func (t BaseNFT) GetOwner() sdk.AccAddress       { return t.Owner }
func (t *BaseNFT) SetOwner(owner sdk.AccAddress) { t.Owner = owner }
func (t BaseNFT) GetTokenID() string             { return "" }

func (t BaseNFT) String() string {
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
		Name:     t.GetName(),
		Symbol:   t.GetSymbol(),
		Decimals: t.GetDecimals(),
		TokenURI: t.GetTokenURI(),
		Mintable: t.GetMintable(),
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
		Name:     t.GetName(),
		Symbol:   t.GetSymbol(),
		TokenURI: t.GetTokenURI(),
		Owner:    t.GetOwner(),
	})
}
func (t *BaseNFT) UnmarshalJSON(data []byte) error {
	type msgAlias *BaseNFT
	return json.Unmarshal(data, msgAlias(t))
}
func (t BaseCollectiveFT) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name     string  `json:"name"`
		Symbol   string  `json:"symbol"`
		TokenURI string  `json:"token_uri"`
		Decimals sdk.Int `json:"decimals"`
		Mintable bool    `json:"mintable"`
		TokenID  string  `json:"token_id"`
	}{
		Name:     t.GetName(),
		Symbol:   t.GetSymbol(),
		TokenURI: t.GetTokenURI(),
		Decimals: t.GetDecimals(),
		Mintable: t.GetMintable(),
		TokenID:  t.GetTokenID(),
	})
}
func (t *BaseCollectiveFT) UnmarshalJSON(data []byte) error {
	rawStruct := struct {
		Name     string  `json:"name"`
		Symbol   string  `json:"symbol"`
		TokenURI string  `json:"token_uri"`
		Decimals sdk.Int `json:"decimals"`
		Mintable bool    `json:"mintable"`
		TokenID  string  `json:"token_id"`
	}{}
	if err := json.Unmarshal(data, &rawStruct); err != nil {
		return err
	}
	t.Name = rawStruct.Name
	t.symbol = rawStruct.Symbol
	t.TokenURI = rawStruct.TokenURI
	t.Decimals = rawStruct.Decimals
	t.Mintable = rawStruct.Mintable
	t.TokenID = rawStruct.TokenID
	return nil
}
func (t BaseCollectiveNFT) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name     string         `json:"name"`
		Symbol   string         `json:"symbol"`
		TokenURI string         `json:"token_uri"`
		Owner    sdk.AccAddress `json:"owner"`
		TokenID  string         `json:"token_id"`
	}{
		Name:     t.GetName(),
		Symbol:   t.GetSymbol(),
		TokenURI: t.GetTokenURI(),
		Owner:    t.GetOwner(),
		TokenID:  t.GetTokenID(),
	})
}
func (t *BaseCollectiveNFT) UnmarshalJSON(data []byte) error {
	rawStruct := struct {
		Name     string         `json:"name"`
		Symbol   string         `json:"symbol"`
		TokenURI string         `json:"token_uri"`
		Owner    sdk.AccAddress `json:"owner"`
		TokenID  string         `json:"token_id"`
	}{}
	if err := json.Unmarshal(data, &rawStruct); err != nil {
		return err
	}
	t.Name = rawStruct.Name
	t.symbol = rawStruct.Symbol
	t.TokenURI = rawStruct.TokenURI
	t.Owner = rawStruct.Owner
	t.TokenID = rawStruct.TokenID
	return nil
}
