package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QuerierRoute     = "token"
	QueryTokens      = "tokens"
	QueryPerms       = "perms"
	QueryCollections = "collections"
	QuerySupply      = "supply"
)

type NodeQuerier interface {
	QueryWithData(path string, data []byte) ([]byte, int64, error)
}

type QueryTokenParams struct {
	Symbol string `json:"symbol"`
}

func (r QueryTokenParams) String() string {
	return r.Symbol
}

func NewQueryTokenParams(symbol string) QueryTokenParams {
	return QueryTokenParams{Symbol: symbol}
}

type QuerySupplyParams struct {
	Symbol string `json:"symbol"`
}

func (r QuerySupplyParams) String() string {
	return r.Symbol
}

func NewQuerySupplyParams(symbol string) QuerySupplyParams {
	return QuerySupplyParams{Symbol: symbol}
}

type QueryAccountPermissionParams struct {
	Addr sdk.AccAddress `json:"addr"`
}

func (r QueryAccountPermissionParams) String() string {
	return r.Addr.String()
}

func NewQueryAccountPermissionParams(addr sdk.AccAddress) QueryAccountPermissionParams {
	return QueryAccountPermissionParams{Addr: addr}
}

type QueryCollectionParams struct {
	Symbol string `json:"symbol"`
}

func (r QueryCollectionParams) String() string {
	return r.Symbol
}

func NewQueryCollectionParams(symbol string) QueryCollectionParams {
	return QueryCollectionParams{Symbol: symbol}
}
