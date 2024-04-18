package keeper

import (
	"context"

	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

var _ types.QueryServer = Keeper{}

// Swapped implements types.QueryServer.
func (k Keeper) Swapped(context.Context, *types.QuerySwappedRequest) (*types.QuerySwappedResponse, error) {
	panic("unimplemented")
}

// TotalNewCurrencySwapLimit implements types.QueryServer.
func (k Keeper) TotalNewCurrencySwapLimit(context.Context, *types.QueryTotalSwappableAmountRequest) (*types.QueryTotalSwappableAmountResponse, error) {
	panic("unimplemented")
}
