package keeper

import (
	"context"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

//var _ types.QueryServer = Keeper{}
//
//// todo this files is just for test
//// Swapped implements types.QueryServer.
//func (k Keeper) Swapped(ctx context.Context, req *types.QuerySwappedRequest) (*types.QuerySwappedResponse, error) {
//	sdkCtx := sdk.UnwrapSDKContext(ctx)
//	check := k.GetFswap(sdkCtx)
//
//	return &types.QuerySwappedResponse{
//		OldCoinAmount: sdk.NewCoin(check.FromDenom, sdk.ZeroInt()),
//		NewCoinAmount: sdk.NewCoin(check.ToDenom, sdk.ZeroInt()),
//	}, nil
//}
//
//// TotalSwappableAmount implements types.QueryServer.
//func (k Keeper) TotalSwappableAmount(context.Context, *types.QueryTotalSwappableAmountRequest) (*types.QueryTotalSwappableAmountResponse, error) {
//	//TODO implement me
//	panic("unimplemented")
//	// sdkCtx := sdk.UnwrapSDKContext(ctx)
//	// totalSwappableAmount := k.GetParams(sdkCtx).SwappableNewCoinAmount
//	// swapped := k.GetSwapped(sdkCtx)
//	// swappableNewCoin := sdk.NewCoin(k.config.NewCoinDenom, totalSwappableAmount.Sub(swapped.NewCoinAmount))
//	// return &types.QueryTotalSwappableAmountResponse{
//	// 	SwappableNewCoin: swappableNewCoin,
//	// }, nil
//}

var _ types.QueryServer = QueryServer{}

type QueryServer struct {
	Keeper
}

func NewQueryServer(keeper Keeper) *QueryServer {
	return &QueryServer{
		keeper,
	}
}

func (s QueryServer) Swapped(ctx context.Context, _ *types.QuerySwappedRequest) (*types.QuerySwappedResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	allSwapped := s.Keeper.GetAllSwapped(c)
	if len(allSwapped) == 0 {
		return nil, types.ErrSwappedNotFound
	}
	swapped := allSwapped[0]
	return &types.QuerySwappedResponse{
		OldCoinAmount: swapped.GetOldCoinAmount(),
		NewCoinAmount: swapped.GetNewCoinAmount(),
	}, nil
}

func (s QueryServer) TotalSwappableAmount(_ context.Context, _ *types.QueryTotalSwappableAmountRequest) (*types.QueryTotalSwappableAmountResponse, error) {
	// s.Keeper.GetSwappableNewCoinAmount()
	return nil, nil
}
