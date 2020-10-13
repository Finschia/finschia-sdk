package types

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	QueryApprovers  = "approvers"
)

type NodeQuerier interface {
	QueryWithData(path string, data []byte) ([]byte, int64, error)
	WithHeight(height int64) context.CLIContext
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

type QueryProxyParams struct {
	Proxy sdk.AccAddress `json:"proxy"`
}

func NewQueryApproverParams(proxy sdk.AccAddress) QueryProxyParams {
	return QueryProxyParams{Proxy: proxy}
}

func IsAddressContains(addresses []sdk.AccAddress, address sdk.AccAddress) bool {
	for _, it := range addresses {
		if address.Equals(it) {
			return true
		}
	}
	return false
}
