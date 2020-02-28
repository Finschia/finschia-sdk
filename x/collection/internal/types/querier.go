package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/client"
)

const (
	QuerierRoute     = ModuleName
	QueryBalance     = "balance"
	QueryTokens      = "tokens"
	QueryTokenTypes  = "tokentypess"
	QueryPerms       = "perms"
	QueryCollections = "collections"
	QuerySupply      = "supply"
	QueryMint        = "mint"
	QueryBurn        = "burn"
	QueryNFTCount    = "nftcount"
	QueryParent      = "parent"
	QueryRoot        = "root"
	QueryChildren    = "children"
	QueryIsApproved  = "approved"
)

type NodeQuerier interface {
	QueryWithData(path string, data []byte) ([]byte, int64, error)
	WithHeight(height int64) client.CLIContext
}

type QuerySymbolParams struct {
	Symbol string `json:"symbol"`
}

func NewQuerySymbolParams(symbol string) QuerySymbolParams {
	return QuerySymbolParams{Symbol: symbol}
}

type QuerySymbolTokenIDParams struct {
	Symbol  string `json:"symbol"`
	TokenID string `json:"token_id"`
}

func NewQuerySymbolTokenIDParams(symbol, tokenID string) QuerySymbolTokenIDParams {
	return QuerySymbolTokenIDParams{Symbol: symbol, TokenID: tokenID}
}

type QuerySymbolTokenIDAccAddressParams struct {
	Symbol  string         `json:"symbol"`
	TokenID string         `json:"token_id"`
	Addr    sdk.AccAddress `json:"addr"`
}

func NewQuerySymbolTokenIDAccAddressParams(symbol, tokenID string, addr sdk.AccAddress) QuerySymbolTokenIDAccAddressParams {
	return QuerySymbolTokenIDAccAddressParams{Symbol: symbol, TokenID: tokenID, Addr: addr}
}

type QueryAccAddressParams struct {
	Addr sdk.AccAddress `json:"addr"`
}

func NewQueryAccAddressParams(addr sdk.AccAddress) QueryAccAddressParams {
	return QueryAccAddressParams{Addr: addr}
}

type QueryIsApprovedParams struct {
	Proxy    sdk.AccAddress `json:"proxy"`
	Approver sdk.AccAddress `json:"approver"`
	Symbol   string         `json:"symbol"`
}

func NewQueryIsApprovedParams(proxy sdk.AccAddress, approver sdk.AccAddress, symbol string) QueryIsApprovedParams {
	return QueryIsApprovedParams{
		Proxy:    proxy,
		Approver: approver,
		Symbol:   symbol,
	}
}
