package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Findable interface {
	IDAtIndex(index int) string
	Len() int
}
type Collection interface {
	GetSymbol() string
	GetName() string
	GetAllTokens() Tokens
	GetFTokens() Tokens
	GetNFTokens() Tokens
	GetToken(string) (Token, sdk.Error)
	AddToken(Token) (Collection, sdk.Error)
	UpdateToken(Token) (Collection, sdk.Error)
	DeleteToken(Token) (Collection, sdk.Error)
	HasToken(string) bool
	NextTokenID(string) string
	NextTokenTypeNFT() string
	NextTokenTypeFT() string
	String() string
}
type BaseCollection struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
	Tokens Tokens `json:"tokens"`
}

func NewCollection(symbol, name string) Collection {
	return &BaseCollection{
		Symbol: symbol,
		Name:   name,
	}
}

func (c BaseCollection) GetSymbol() string { return c.Symbol }
func (c BaseCollection) GetName() string   { return c.Name }

func (c BaseCollection) GetAllTokens() Tokens {
	var tokens Tokens
	tokens = append(tokens, c.GetFTokens()...)
	tokens = append(tokens, c.GetNFTokens()...)
	return tokens.Sort()
}

func (c BaseCollection) GetFTokens() Tokens {
	var tokens Tokens
	for _, token := range c.Tokens.GetFTs() {
		token.(Token).SetCollection(&c)
		tokens = append(tokens, token)
	}
	return tokens.Sort()
}

func (c BaseCollection) GetNFTokens() Tokens {
	var tokens Tokens
	for _, token := range c.Tokens.GetNFTs() {
		token.(Token).SetCollection(&c)
		tokens = append(tokens, token)
	}
	return tokens.Sort()
}
func (c BaseCollection) GetToken(tokenID string) (Token, sdk.Error) {
	token, found := c.Tokens.Find(tokenID)
	if found {
		token.(Token).SetCollection(&c)
		return token, nil
	}
	return nil, ErrCollectionTokenNotExist(DefaultCodespace, c.Symbol, tokenID)
}

func (c BaseCollection) AddToken(token Token) (Collection, sdk.Error) {
	if c.HasToken(token.GetTokenID()) {
		return c, ErrCollectionTokenExist(DefaultCodespace, c.Symbol, token.GetTokenID())
	}
	c.Tokens = c.Tokens.Append(token)
	return c, nil
}

func (c BaseCollection) UpdateToken(token Token) (Collection, sdk.Error) {
	tokens, ok := c.Tokens.Update(token)
	if !ok {
		return c, ErrCollectionTokenNotExist(DefaultCodespace, c.Symbol, token.GetTokenID())
	}
	c.Tokens = tokens
	return c, nil
}

func (c BaseCollection) DeleteToken(token Token) (Collection, sdk.Error) {
	ids, ok := c.Tokens.Remove(token.GetTokenID())
	if !ok {
		return c, ErrCollectionTokenNotExist(DefaultCodespace, c.Symbol, token.GetTokenID())
	}
	c.Tokens = ids
	return c, nil
}

func (c BaseCollection) HasToken(tokenID string) bool {
	_, err := c.GetToken(tokenID)
	return err == nil
}

func (c BaseCollection) NextTokenID(prefix string) string {
	return c.Tokens.NextTokenID(prefix)
}

func (c BaseCollection) NextTokenTypeNFT() string {
	return c.Tokens.NextTokenTypeForNFT()
}

func (c BaseCollection) NextTokenTypeFT() string {
	return c.Tokens.NextTokenTypeForFT()
}

func (c BaseCollection) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(b)
}

type Collections []Collection

func (collections Collections) String() string {
	b, err := json.Marshal(collections)
	if err != nil {
		panic(err)
	}
	return string(b)
}
