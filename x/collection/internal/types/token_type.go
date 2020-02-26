package types

import "encoding/json"

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
