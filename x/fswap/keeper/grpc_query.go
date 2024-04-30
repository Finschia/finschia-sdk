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
		FromCoinAmount: swapped.GetFromCoinAmount(),
		ToCoinAmount:   swapped.GetToCoinAmount(),
	}, nil
}

func (s QueryServer) TotalSwappableToCoinAmount(ctx context.Context, _ *types.QueryTotalSwappableToCoinAmountRequest) (*types.QueryTotalSwappableToCoinAmountResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)

	amount, err := s.Keeper.getSwappableNewCoinAmount(c)
	if err != nil {
		return &types.QueryTotalSwappableToCoinAmountResponse{}, err
	}

	return &types.QueryTotalSwappableToCoinAmountResponse{SwappableAmount: amount}, nil
}
