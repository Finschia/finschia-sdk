package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Token interface {
	GetName() string
	GetSymbol() string
	GetTokenID() string
	GetTokenURI() string
	SetTokenURI(tokenURI string)
	GetTokenType() string
	GetTokenIndex() string
	SetCollection(Collection)
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

var _ Token = (*BaseNFT)(nil)
var _ json.Marshaler = (*BaseNFT)(nil)
var _ json.Unmarshaler = (*BaseNFT)(nil)

type BaseNFT struct {
	TokenID  string         `json:"token_id"`
	TokenURI string         `json:"token_uri"`
	Owner    sdk.AccAddress `json:"owner"`
	Name     string         `json:"Name"`

	// volatile
	symbol string
}

func NewNFT(collection Collection, name, tokenType, tokenURI string, owner sdk.AccAddress) NFT {
	tokenID := collection.NextTokenID(tokenType)
	return &BaseNFT{
		TokenID:  tokenID,
		TokenURI: tokenURI,
		Owner:    owner,
		Name:     name,
		symbol:   collection.GetSymbol(),
	}
}
func (t BaseNFT) GetName() string                { return t.Name }
func (t BaseNFT) GetSymbol() string              { return t.symbol }
func (t BaseNFT) GetTokenURI() string            { return t.TokenURI }
func (t BaseNFT) GetOwner() sdk.AccAddress       { return t.Owner }
func (t BaseNFT) GetTokenID() string             { return t.TokenID }
func (t BaseNFT) GetTokenType() string           { return t.TokenID[:TokenTypeLength] }
func (t BaseNFT) GetTokenIndex() string          { return t.TokenID[TokenTypeLength:] }
func (t *BaseNFT) SetOwner(owner sdk.AccAddress) { t.Owner = owner }
func (t *BaseNFT) SetTokenURI(tokenURI string)   { t.TokenURI = tokenURI }
func (t *BaseNFT) SetCollection(c Collection)    { t.symbol = c.GetSymbol() }
func (t BaseNFT) String() string                 { return "" }

var _ Token = (*BaseFT)(nil)
var _ FT = (*BaseFT)(nil)
var _ json.Marshaler = (*BaseFT)(nil)
var _ json.Unmarshaler = (*BaseFT)(nil)

type BaseFT struct {
	TokenID  string  `json:"token_id"`
	TokenURI string  `json:"token_uri"`
	Decimals sdk.Int `json:"decimals"`
	Mintable bool    `json:"mintable"`
	Name     string  `json:"Name"`

	// volatile
	symbol string
}

func NewFT(collection Collection, name, tokenURI string, decimals sdk.Int, mintable bool) FT {
	tokenID := collection.NextTokenTypeFT() + "0000"
	return &BaseFT{
		TokenID:  tokenID,
		TokenURI: tokenURI,
		Decimals: decimals,
		Mintable: mintable,
		Name:     name,
		symbol:   collection.GetSymbol(),
	}
}
func (t BaseFT) GetName() string              { return t.Name }
func (t BaseFT) GetSymbol() string            { return t.symbol }
func (t BaseFT) GetTokenURI() string          { return t.TokenURI }
func (t BaseFT) GetMintable() bool            { return t.Mintable }
func (t BaseFT) GetDecimals() sdk.Int         { return t.Decimals }
func (t BaseFT) GetTokenID() string           { return t.TokenID }
func (t BaseFT) GetTokenType() string         { return t.TokenID[:TokenTypeLength] }
func (t BaseFT) GetTokenIndex() string        { return t.TokenID[TokenTypeLength:] }
func (t *BaseFT) SetTokenURI(tokenURI string) { t.TokenURI = tokenURI }
func (t *BaseFT) SetCollection(c Collection)  { t.symbol = c.GetSymbol() }
func (t BaseFT) String() string               { return "" }

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
func (t *BaseFT) UnmarshalJSON(data []byte) error {
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
func (t BaseNFT) MarshalJSON() ([]byte, error) {
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
func (t *BaseNFT) UnmarshalJSON(data []byte) error {
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
