package types

import (
	sdk "github.com/line/lbm-sdk/types"
)

// QueryBalanceParams defines the params for querying an account balance.
type QueryBalanceParams struct {
	Address sdk.AccAddress
}

// NewQueryBalanceParams creates a new instance of QueryBalanceParams.
func NewQueryBalanceParams(addr sdk.AccAddress) QueryBalanceParams {
	return QueryBalanceParams{Address: addr}
}
