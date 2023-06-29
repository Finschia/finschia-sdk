package keeper

import (
	"context"

	"github.com/Finschia/finschia-sdk/store/prefix"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/types/query"
	"github.com/Finschia/finschia-sdk/x/or/rollup/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Rollup(ctx context.Context, req *types.QueryRollupRequest) (*types.QueryRollupResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	val, found := k.GetRollup(
		sdkCtx,
		req.RollupName,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}
	res := types.QueryRollupResponse{Rollup: val}

	return &res, nil
}

func (k Keeper) AllRollup(ctx context.Context, req *types.QueryAllRollupRequest) (*types.QueryAllRollupResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var rollups []types.Rollup
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := sdkCtx.KVStore(k.storeKey)
	rollupStore := prefix.NewStore(store, types.RollupKeyPrefix)

	pageRes, err := query.Paginate(rollupStore, req.Pagination, func(key []byte, value []byte) error {
		var rollapp types.Rollup
		if err := k.cdc.Unmarshal(value, &rollapp); err != nil {
			return err
		}

		rollups = append(rollups, rollapp)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	res := types.QueryAllRollupResponse{Rollup: rollups, Pagination: pageRes}

	return &res, nil
}

func (k Keeper) Sequencer(ctx context.Context, req *types.QuerySequencerRequest) (*types.QuerySequencerResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	k.GetSequencersByRollupName(sdkCtx, req.RollupName)
	res := types.QuerySequencerResponse{}

	return &res, nil
}

func (k Keeper) SequencersByRollup(ctx context.Context, req *types.QuerySequencersByRollupRequest) (*types.QuerySequencersByRollupResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sequencersByRollup, found := k.GetSequencersByRollupName(sdkCtx, req.RollupName)
	if !found {
		return &types.QuerySequencersByRollupResponse{}, status.Error(codes.NotFound, "not found")
	}

	res := types.QuerySequencersByRollupResponse{RollupName: req.RollupName, SequencerList: sequencersByRollup.Sequencers}

	return &res, nil
}
