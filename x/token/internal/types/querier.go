package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link-modules/client"
)

const (
	QuerierRoute    = "token"
	QueryTokens     = "tokens"
	QueryPerms      = "perms"
	QueryBalance    = "balance"
	QuerySupply     = "supply"
	QueryMint       = "mint"
	QueryBurn       = "burn"
	QueryIsApproved = "approved"
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

type QueryIsApprovedParams struct {
	Proxy    sdk.AccAddress `json:"proxy"`
	Approver sdk.AccAddress `json:"approver"`
}

func NewQueryIsApprovedParams(proxy sdk.AccAddress, approver sdk.AccAddress) QueryIsApprovedParams {
	return QueryIsApprovedParams{
		Proxy:    proxy,
		Approver: approver,
	}
}
