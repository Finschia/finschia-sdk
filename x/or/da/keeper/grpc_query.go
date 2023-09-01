package keeper

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Finschia/finschia-sdk/store/prefix"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/types/query"
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

func (k Keeper) CCState(goCtx context.Context, req *types.QueryCCStateRequest) (*types.QueryCCStateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.RollupName == "" {
		return nil, status.Error(codes.InvalidArgument, "empty rollup name")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	ccState, err := k.GetCCState(ctx, req.RollupName)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCCStateResponse{State: ccState}, nil
}

func (k Keeper) CCRef(goCtx context.Context, req *types.QueryCCRefRequest) (*types.QueryCCRefResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.RollupName == "" {
		return nil, status.Error(codes.InvalidArgument, "empty rollup name")
	}

	if req.BatchHeight < 1 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid batch height %d: must be greater than 1", req.BatchHeight)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	ccRef, err := k.GetCCRef(ctx, req.RollupName, req.BatchHeight)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCCRefResponse{Ref: ccRef}, nil
}

func (k Keeper) CCRefs(goCtx context.Context, req *types.QueryCCRefsRequest) (*types.QueryCCRefsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.RollupName == "" {
		return nil, status.Error(codes.InvalidArgument, "empty rollup name")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	key := types.GenRollupPrefix(req.RollupName, types.CCBatchIndexPrefix)
	qtxStore := prefix.NewStore(ctx.KVStore(k.storeKey), key)
	ccRefs, pageRes, err := query.GenericFilteredPaginate(k.cdc, qtxStore, req.Pagination,
		func(key []byte, value *types.CCRef) (*types.CCRef, error) {
			if len(key) != 8 {
				panic(fmt.Sprintf("unexpected key length %d", len(key)))
			}

			return value, nil
		},
		func() *types.CCRef { return &types.CCRef{} },
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCCRefsResponse{Refs: ccRefs, Pagination: pageRes}, nil
}

func (k Keeper) QueueTxState(goCtx context.Context, req *types.QueryQueueTxStateRequest) (*types.QueryQueueTxStateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.RollupName == "" {
		return nil, status.Error(codes.InvalidArgument, "empty rollup name")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	qtxState, err := k.GetQueueTxState(ctx, req.RollupName)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryQueueTxStateResponse{State: qtxState}, nil
}

func (k Keeper) QueueTx(goCtx context.Context, req *types.QueryQueueTxRequest) (*types.QueryQueueTxResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.RollupName == "" {
		return nil, status.Error(codes.InvalidArgument, "empty rollup name")
	}

	if req.QueueIndex < 1 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid batch height %d: must be greater than 1", req.QueueIndex)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	qtx, err := k.GetQueueTx(ctx, req.RollupName, req.QueueIndex)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryQueueTxResponse{Tx: qtx}, nil
}

func (k Keeper) QueueTxs(goCtx context.Context, req *types.QueryQueueTxsRequest) (*types.QueryQueueTxsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.RollupName == "" {
		return nil, status.Error(codes.InvalidArgument, "empty rollup name")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	key := types.GenRollupPrefix(req.RollupName, types.CCQueueTxPrefix)
	qtxStore := prefix.NewStore(ctx.KVStore(k.storeKey), key)
	qtxs, pageRes, err := query.GenericFilteredPaginate(k.cdc, qtxStore, req.Pagination,
		func(key []byte, value *types.L1ToL2Queue) (*types.L1ToL2Queue, error) {
			if len(key) != 8 {
				panic(fmt.Sprintf("unexpected key length %d", len(key)))
			}

			return value, nil
		},
		func() *types.L1ToL2Queue { return &types.L1ToL2Queue{} },
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryQueueTxsResponse{Txs: qtxs, Pagination: pageRes}, nil
}

func (k Keeper) MappedBatch(goCtx context.Context, req *types.QueryMappedBatchRequest) (*types.QueryMappedBatchResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.RollupName == "" {
		return nil, status.Error(codes.InvalidArgument, "empty rollup name")
	}

	if req.L2Height < 1 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid rollup height %d: must be greater than 1", req.L2Height)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	batchIdx, err := k.GetL2HeightBatchMap(ctx, req.RollupName, req.L2Height)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	ccRef, err := k.GetCCRef(ctx, req.RollupName, batchIdx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryMappedBatchResponse{Ref: ccRef}, nil
}

func (k Keeper) SCCState(goCtx context.Context, req *types.QuerySCCStateRequest) (*types.QuerySCCStateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.RollupName == "" {
		return nil, status.Error(codes.InvalidArgument, "empty rollup name")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	sccState, err := k.GetSCCState(ctx, req.RollupName)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QuerySCCStateResponse{State: sccState}, nil
}

func (k Keeper) SCCRef(goCtx context.Context, req *types.QuerySCCRefRequest) (*types.QuerySCCRefResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.RollupName == "" {
		return nil, status.Error(codes.InvalidArgument, "empty rollup name")
	}

	if req.BatchHeight < 1 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid batch height %d: must be greater than 1", req.BatchHeight)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	sccRef, err := k.GetSCCRef(ctx, req.RollupName, req.BatchHeight)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	res := new(types.SCCRes)
	res.Ref = *sccRef
	isSCCConfirmed(ctx, res)

	return &types.QuerySCCRefResponse{Ref: res}, nil
}

func (k Keeper) SCCRefs(goCtx context.Context, req *types.QuerySCCRefsRequest) (*types.QuerySCCRefsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.RollupName == "" {
		return nil, status.Error(codes.InvalidArgument, "empty rollup name")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	key := types.GenRollupPrefix(req.RollupName, types.SCCBatchIndexPrefix)
	qtxStore := prefix.NewStore(ctx.KVStore(k.storeKey), key)
	sccRefs, pageRes, err := query.GenericFilteredPaginate(k.cdc, qtxStore, req.Pagination,
		func(key []byte, value *types.SCCRef) (*types.SCCRef, error) {
			if len(key) != 8 {
				panic(fmt.Sprintf("unexpected key length %d", len(key)))
			}

			return value, nil
		},
		func() *types.SCCRef { return &types.SCCRef{} },
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	allRes := make([]*types.SCCRes, len(sccRefs))
	for _, sccRef := range sccRefs {
		res := new(types.SCCRes)
		res.Ref = *sccRef
		isSCCConfirmed(ctx, res)
		allRes = append(allRes, res)
	}

	return &types.QuerySCCRefsResponse{Refs: allRes, Pagination: pageRes}, nil
}

func (k Keeper) LastSequencerTimestamp(goCtx context.Context, req *types.QueryLastSequencerTimestampRequest) (*types.QueryLastSequencerTimestampResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.RollupName == "" {
		return nil, status.Error(codes.InvalidArgument, "empty rollup name")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	state, err := k.GetSCCState(ctx, req.RollupName)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryLastSequencerTimestampResponse{LastSequencerSubmit: state.LastSequencerSubmit}, nil
}

func isSCCConfirmed(ctx sdk.Context, target *types.SCCRes) {
	if ctx.BlockTime().UTC().After(target.Ref.Deadline) {
		target.IsConfirmed = true
	} else {
		target.IsConfirmed = false
	}
}
