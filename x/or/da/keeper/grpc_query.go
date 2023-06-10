package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/or/da/types"
)

var _ types.QueryServer = Keeper{}

// Params queries params of da module
func (k Keeper) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	p := k.GetParams(ctx)
	return &types.QueryParamsResponse{Params: p}, nil
}

func (k Keeper) CCBatches(context.Context, *types.QueryCCBatchesRequest) (*types.QueryCCBatchesResponse, error) {
	panic("implement me")
}

func (k Keeper) CCBatch(context.Context, *types.QueryCCBatchRequest) (*types.QueryCCBatchResponse, error) {
	panic("implement me")
}
