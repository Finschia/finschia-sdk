package types

import (
	"encoding/json"
)

const (
	TokenTypeLength    = 8
	TokenIDLength      = 16
	SmallestAlphanum   = "0"
	LargestAlphanum    = "z"
	FungibleFlag       = SmallestAlphanum
	ReservedEmpty      = "00000000"
	SmallestFTType     = "00000001"
	SmallestNFTType    = "10000001"
	SmallestTokenIndex = "00000001"
)

type TokenType interface {
	GetName() string
	SetName(string) TokenType
	GetSymbol() string
	GetTokenType() string
	String() string
}

type BaseTokenType struct {
	Symbol    string `json:"symbol"`
	TokenType string `json:"token_type"`
	Name      string `json:"name"`
}

func NewBaseTokenType(symbol, tokenType, name string) TokenType {
	return &BaseTokenType{
		Symbol:    symbol,
		TokenType: tokenType,
		Name:      name,
	}
}
func (t BaseTokenType) GetName() string { return t.Name }
func (t BaseTokenType) SetName(name string) TokenType {
	t.Name = name
	return t
}
func (t BaseTokenType) GetSymbol() string    { return t.Symbol }
func (t BaseTokenType) GetTokenType() string { return t.TokenType }
func (t BaseTokenType) String() string {
	b, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(b)
}

type TokenTypes []TokenType

func (ts TokenTypes) String() string {
	b, err := json.Marshal(ts)
	if err != nil {
		panic(err)
	}
	return string(b)
}
