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
	swapped, err := s.Keeper.GetSwapped(c)
	if err != nil {
		return nil, err
	}
	return &types.QuerySwappedResponse{Swapped: swapped}, nil
}

func (s QueryServer) TotalNewCurrencySwapLimit(_ context.Context, _ *types.QueryTotalSwappableAmountRequest) (*types.QueryTotalSwappableAmountResponse, error) {
	// s.Keeper.GetSwappableNewCoinAmount()
	return nil, nil
}
