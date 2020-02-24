package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/client"
)

const (
	QuerierRoute = "token"
	QueryTokens  = "tokens"
	QueryPerms   = "perms"
	QueryBalance = "balance"
	QuerySupply  = "supply"
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

type QueryAccAddressParams struct {
	Addr sdk.AccAddress `json:"addr"`
}

func NewQueryAccAddressParams(addr sdk.AccAddress) QueryAccAddressParams {
	return QueryAccAddressParams{Addr: addr}
}

type QuerySymbolAccAddressParams struct {
	Symbol string         `json:"symbol"`
	Addr   sdk.AccAddress `json:"addr"`
}

func NewQuerySymbolAccAddressParams(symbol string, addr sdk.AccAddress) QuerySymbolAccAddressParams {
	return QuerySymbolAccAddressParams{Symbol: symbol, Addr: addr}
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
