package keeper

import (
	"context"

	"github.com/Finschia/finschia-rdk/store/prefix"
	sdk "github.com/Finschia/finschia-rdk/types"
	"github.com/Finschia/finschia-rdk/types/query"
	"github.com/Finschia/finschia-rdk/x/or/rollup/types"
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

	sequencer, found := k.GetSequencer(sdkCtx, req.SequencerAddress)
	if !found {
		return &types.QuerySequencerResponse{}, status.Error(codes.NotFound, "not found")
	}
	res := types.QuerySequencerResponse{
		Sequencer: sequencer,
	}

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

func (k Keeper) Deposit(ctx context.Context, req *types.QueryDepositRequest) (*types.QueryDepositResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sequencers, found := k.GetSequencersByRollupName(sdkCtx, req.RollupName)
	if !found {
		return &types.QueryDepositResponse{}, status.Error(codes.NotFound, "not found")
	}

	sequencerContains := false
	for _, v := range sequencers.Sequencers {
		if v.SequencerAddress == req.SequencerAddress {
			sequencerContains = true
			break
		}
	}

	if !sequencerContains {
		return nil, types.ErrNotExistSequencer
	}

	deposit, found := k.GetDeposit(sdkCtx, req.RollupName, req.SequencerAddress)
	if !found {
		return nil, types.ErrNotFoundDeposit
	}

	res := types.QueryDepositResponse{
		Deposit: &deposit,
	}

	return &res, nil
}
