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
	NextBaseID() string
	NexTokenIDForFT() string
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
	for _, token := range c.Tokens {
		token.(CollectiveToken).SetCollection(&c)
		tokens = append(tokens, token)
	}
	return tokens.Sort()
}

func (c BaseCollection) GetFTokens() Tokens {
	var tokens Tokens
	for _, token := range c.Tokens.GetFTs() {
		token.(CollectiveToken).SetCollection(&c)
		tokens = append(tokens, token)
	}
	return tokens.Sort()
}

func (c BaseCollection) GetNFTokens() Tokens {
	var tokens Tokens
	for _, token := range c.Tokens.GetNFTs() {
		token.(CollectiveToken).SetCollection(&c)
		tokens = append(tokens, token)
	}
	return tokens.Sort()
}
func (c BaseCollection) GetToken(tokenID string) (Token, sdk.Error) {
	token, found := c.Tokens.Find(tokenID)
	if found {
		token.(CollectiveToken).SetCollection(&c)
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

func (c BaseCollection) NextBaseID() string {
	return c.Tokens.NextBaseID()
}

func (c BaseCollection) NexTokenIDForFT() string {
	return c.Tokens.NextTokenID(FungibleFlag)
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

var _ Token = (*BaseCollectiveNFT)(nil)
var _ NFT = (*BaseCollectiveNFT)(nil)
var _ CollectiveToken = (*BaseCollectiveNFT)(nil)
var _ json.Marshaler = (*BaseCollectiveNFT)(nil)
var _ json.Unmarshaler = (*BaseCollectiveNFT)(nil)

type CollectiveToken interface {
	Token
	SetCollection(Collection)
}

type CollectiveNFT interface {
	NFT
	SetCollection(Collection)
}

type CollectiveFT interface {
	FT
	SetCollection(Collection)
}

type BaseCollectiveNFT struct {
	TokenID  string         `json:"token_id"`
	TokenURI string         `json:"token_uri"`
	Owner    sdk.AccAddress `json:"owner"`
	Name     string         `json:"Name"`

	// volatile
	symbol string
}

func NewCollectiveNFT(collection Collection, name, tokenID, tokenURI string, owner sdk.AccAddress) NFT {
	if tokenID == "" {
		tokenID = collection.NextBaseID()
	}
	if len(tokenID) == BaseTokenIDLength {
		tokenID = collection.NextTokenID(tokenID)
	}
	return &BaseCollectiveNFT{
		TokenID:  tokenID,
		TokenURI: tokenURI,
		Owner:    owner,
		Name:     name,
		symbol:   collection.GetSymbol(),
	}
}
func (t BaseCollectiveNFT) GetName() string                { return t.Name }
func (t BaseCollectiveNFT) GetSymbol() string              { return t.symbol }
func (t BaseCollectiveNFT) GetDenom() string               { return t.symbol + t.TokenID }
func (t BaseCollectiveNFT) GetTokenURI() string            { return t.TokenURI }
func (t BaseCollectiveNFT) GetOwner() sdk.AccAddress       { return t.Owner }
func (t BaseCollectiveNFT) GetTokenID() string             { return t.TokenID }
func (t *BaseCollectiveNFT) SetOwner(owner sdk.AccAddress) { t.Owner = owner }
func (t *BaseCollectiveNFT) SetTokenURI(tokenURI string)   { t.TokenURI = tokenURI }
func (t *BaseCollectiveNFT) SetCollection(c Collection)    { t.symbol = c.GetSymbol() }
func (t BaseCollectiveNFT) String() string                 { return "" }

var _ Token = (*BaseCollectiveFT)(nil)
var _ FT = (*BaseCollectiveFT)(nil)
var _ CollectiveToken = (*BaseCollectiveFT)(nil)
var _ json.Marshaler = (*BaseCollectiveFT)(nil)
var _ json.Unmarshaler = (*BaseCollectiveFT)(nil)

type BaseCollectiveFT struct {
	TokenID  string  `json:"token_id"`
	TokenURI string  `json:"token_uri"`
	Decimals sdk.Int `json:"decimals"`
	Mintable bool    `json:"mintable"`
	Name     string  `json:"Name"`

	// volatile
	symbol string
}

func NewCollectiveFT(collection Collection, name, tokenID string, tokenURI string, decimals sdk.Int, mintable bool) FT {
	if tokenID == "" {
		tokenID = collection.NexTokenIDForFT()
	}
	return &BaseCollectiveFT{
		TokenID:  tokenID,
		TokenURI: tokenURI,
		Decimals: decimals,
		Mintable: mintable,
		Name:     name,
		symbol:   collection.GetSymbol(),
	}
}
func (t BaseCollectiveFT) GetName() string              { return t.Name }
func (t BaseCollectiveFT) GetSymbol() string            { return t.symbol }
func (t BaseCollectiveFT) GetDenom() string             { return t.symbol + t.TokenID }
func (t BaseCollectiveFT) GetTokenURI() string          { return t.TokenURI }
func (t BaseCollectiveFT) GetMintable() bool            { return t.Mintable }
func (t BaseCollectiveFT) GetDecimals() sdk.Int         { return t.Decimals }
func (t BaseCollectiveFT) GetTokenID() string           { return t.TokenID }
func (t *BaseCollectiveFT) SetTokenURI(tokenURI string) { t.TokenURI = tokenURI }
func (t *BaseCollectiveFT) SetCollection(c Collection)  { t.symbol = c.GetSymbol() }
func (t BaseCollectiveFT) String() string               { return "" }
