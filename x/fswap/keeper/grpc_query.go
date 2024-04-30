package keeper

import (
	"context"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

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
	swapped, err := s.Keeper.getSwapped(c)
	if err != nil {
		return nil, err
	}
	return &types.QuerySwappedResponse{
		OldCoinAmount: swapped.GetOldCoinAmount(),
		NewCoinAmount: swapped.GetNewCoinAmount(),
	}, nil
}

func (s QueryServer) TotalSwappableAmount(ctx context.Context, _ *types.QueryTotalSwappableAmountRequest) (*types.QueryTotalSwappableAmountResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	amount, err := s.Keeper.getSwappableNewCoinAmount(c)
	if err != nil {
		return &types.QueryTotalSwappableAmountResponse{}, err
	}
	return &types.QueryTotalSwappableAmountResponse{
		SwappableNewCoin: amount,
	}, nil
}
