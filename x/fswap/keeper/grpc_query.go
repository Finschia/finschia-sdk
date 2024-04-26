package keeper

import (
	"context"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

var _ types.QueryServer = Keeper{}

// todo this files is just for test
// Swapped implements types.QueryServer.
func (k Keeper) Swapped(ctx context.Context, req *types.QuerySwappedRequest) (*types.QuerySwappedResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	check := k.GetFswapInit(sdkCtx)

	return &types.QuerySwappedResponse{
		OldCoin: sdk.NewCoin(check.FromDenom, sdk.ZeroInt()),
		NewCoin: sdk.NewCoin(check.ToDenom, sdk.ZeroInt()),
	}, nil
}

// TotalNewCurrencySwapLimit implements types.QueryServer.
func (k Keeper) TotalNewCurrencySwapLimit(context.Context, *types.QueryTotalSwappableAmountRequest) (*types.QueryTotalSwappableAmountResponse, error) {
	panic("unimplemented")
	// sdkCtx := sdk.UnwrapSDKContext(ctx)
	// totalSwappableAmount := k.GetParams(sdkCtx).SwappableNewCoinAmount
	// swapped := k.GetSwapped(sdkCtx)
	// swappableNewCoin := sdk.NewCoin(k.config.NewCoinDenom, totalSwappableAmount.Sub(swapped.NewCoinAmount))
	// return &types.QueryTotalSwappableAmountResponse{
	// 	SwappableNewCoin: swappableNewCoin,
	// }, nil
}
