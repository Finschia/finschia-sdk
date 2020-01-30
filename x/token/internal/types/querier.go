package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/client"
)

const (
	QuerierRoute     = "token"
	QueryTokens      = "tokens"
	QueryPerms       = "perms"
	QueryCollections = "collections"
	QuerySupply      = "supply"
	QueryNFTCount    = "nftcount"
	QueryParent      = "parent"
	QueryRoot        = "root"
	QueryChildren    = "children"
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

type QueryAccAddressParams struct {
	Addr sdk.AccAddress `json:"addr"`
}

func NewQueryAccAddressParams(addr sdk.AccAddress) QueryAccAddressParams {
	return QueryAccAddressParams{Addr: addr}
}
