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

type QueryContractIDParams struct {
	ContractID string `json:"contract_id"`
}

func NewQueryContractIDParams(contractID string) QueryContractIDParams {
	return QueryContractIDParams{ContractID: contractID}
}

type QueryAccAddressParams struct {
	Addr sdk.AccAddress `json:"addr"`
}

func NewQueryAccAddressParams(addr sdk.AccAddress) QueryAccAddressParams {
	return QueryAccAddressParams{Addr: addr}
}

type QueryAccAddressContractIDParams struct {
	Addr       sdk.AccAddress `json:"addr"`
	ContractID string         `json:"contract_id"`
}

func NewQueryAccAddressContractIDParams(contractID string, addr sdk.AccAddress) QueryAccAddressContractIDParams {
	return QueryAccAddressContractIDParams{Addr: addr, ContractID: contractID}
}
