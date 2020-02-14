package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const RequestGetsLimit = 100

// QueryBalanceOfParams defines the params for querying an account balance.
type QueryBalanceParams struct {
	Address sdk.AccAddress
	Denom   string
}

type QueryBulkBalancesParams struct {
	Addresses []sdk.AccAddress `json:"addresses"`
}

type QueryBulkBalancesResult struct {
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}

// NewQueryBalanceOfParams creates a new instance of QueryBalanceParams.
func NewQueryBalanceParams(addr sdk.AccAddress, denom string) QueryBalanceParams {
	return QueryBalanceParams{Address: addr, Denom: denom}
}

func NewQueryBulkBalanceParams(addrs []sdk.AccAddress) QueryBulkBalancesParams {
	return QueryBulkBalancesParams{Addresses: addrs}
}

func NewQueryBulkBalancesResult(address sdk.AccAddress, coins sdk.Coins) QueryBulkBalancesResult {
	return QueryBulkBalancesResult{Address: address, Coins: coins}
}
