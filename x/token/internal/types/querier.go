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
	QueryMint    = "mint"
	QueryBurn    = "burn"
)

type NodeQuerier interface {
	QueryWithData(path string, data []byte) ([]byte, int64, error)
	WithHeight(height int64) client.CLIContext
}

type QueryContractIDAccAddressParams struct {
	Addr sdk.AccAddress `json:"addr"`
}

func NewQueryContractIDAccAddressParams(addr sdk.AccAddress) QueryContractIDAccAddressParams {
	return QueryContractIDAccAddressParams{Addr: addr}
}
