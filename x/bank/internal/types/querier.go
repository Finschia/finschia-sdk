package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// QueryBalanceOfParams defines the params for querying an account balance.
type QueryBalanceOfParams struct {
	Address sdk.AccAddress
	Denom   string
}

// NewQueryBalanceOfParams creates a new instance of QueryBalanceParams.
func NewQueryBalanceOfParams(addr sdk.AccAddress, denom string) QueryBalanceOfParams {
	return QueryBalanceOfParams{Address: addr, Denom: denom}
}
