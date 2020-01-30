package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// QueryBalanceOfParams defines the params for querying an account balance.
type QueryBalanceParams struct {
	Address sdk.AccAddress
	Denom   string
}

// NewQueryBalanceOfParams creates a new instance of QueryBalanceParams.
func NewQueryBalanceParams(addr sdk.AccAddress, denom string) QueryBalanceParams {
	return QueryBalanceParams{Address: addr, Denom: denom}
}
