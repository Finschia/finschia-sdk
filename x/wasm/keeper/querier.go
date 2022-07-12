package keeper

import (
	"context"
	"encoding/binary"
	"runtime/debug"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/line/lbm-sdk/codec"

	"github.com/line/lbm-sdk/store/prefix"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/types/query"
	"github.com/line/lbm-sdk/x/wasm/types"
)

var _ types.QueryServer = &GrpcQuerier{}

type GrpcQuerier struct {
	cdc           codec.Codec
	storeKey      sdk.StoreKey
	keeper        types.ViewKeeper
	queryGasLimit sdk.Gas
}

// NewGrpcQuerier constructor
func NewGrpcQuerier(cdc codec.Codec, storeKey sdk.StoreKey, keeper types.ViewKeeper, queryGasLimit sdk.Gas) *GrpcQuerier { //nolint:revive
	return &GrpcQuerier{cdc: cdc, storeKey: storeKey, keeper: keeper, queryGasLimit: queryGasLimit}
}

func (q GrpcQuerier) ContractInfo(c context.Context, req *types.QueryContractInfoRequest) (*types.QueryContractInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	err := sdk.ValidateAccAddress(req.Address)
	if err != nil {
		return nil, err
	}
	rsp, err := queryContractInfo(sdk.UnwrapSDKContext(c), sdk.AccAddress(req.Address), q.keeper)
	switch {
	case err != nil:
		return nil, err
	case rsp == nil:
		return nil, types.ErrNotFound
	}
	return rsp, nil
}

func (q GrpcQuerier) ContractHistory(c context.Context, req *types.QueryContractHistoryRequest) (*types.QueryContractHistoryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	err := sdk.ValidateAccAddress(req.Address)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	r := make([]types.ContractCodeHistoryEntry, 0)

	prefixStore := prefix.NewStore(ctx.KVStore(q.storeKey), types.GetContractCodeHistoryElementPrefix(sdk.AccAddress(req.Address)))
	pageRes, err := query.FilteredPaginate(prefixStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		if accumulate {
			var e types.ContractCodeHistoryEntry
			if err := q.cdc.Unmarshal(value, &e); err != nil {
				return false, err
			}
			e.Updated = nil // redact
			r = append(r, e)
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return &types.QueryContractHistoryResponse{
		Entries:    r,
		Pagination: pageRes,
	}, nil
}

// ContractsByCode lists all smart contracts for a code id
func (q GrpcQuerier) ContractsByCode(c context.Context, req *types.QueryContractsByCodeRequest) (*types.QueryContractsByCodeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.CodeId == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "code id")
	}
	ctx := sdk.UnwrapSDKContext(c)
	r := make([]string, 0)

	prefixStore := prefix.NewStore(ctx.KVStore(q.storeKey), types.GetContractByCodeIDSecondaryIndexPrefix(req.CodeId))
	pageRes, err := query.FilteredPaginate(prefixStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		if accumulate {
			var contractAddr = sdk.AccAddress(string(key[types.AbsoluteTxPositionLen:]))
			r = append(r, contractAddr.String())
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return &types.QueryContractsByCodeResponse{
		Contracts:  r,
		Pagination: pageRes,
	}, nil
}

func (q GrpcQuerier) AllContractState(c context.Context, req *types.QueryAllContractStateRequest) (*types.QueryAllContractStateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	err := sdk.ValidateAccAddress(req.Address)
	if err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(c)
	if !q.keeper.HasContractInfo(ctx, sdk.AccAddress(req.Address)) {
		return nil, types.ErrNotFound
	}

	r := make([]types.Model, 0)
	prefixStore := prefix.NewStore(ctx.KVStore(q.storeKey), types.GetContractStorePrefix(sdk.AccAddress(req.Address)))
	pageRes, err := query.FilteredPaginate(prefixStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		if accumulate {
			r = append(r, types.Model{
				Key:   key,
				Value: value,
			})
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return &types.QueryAllContractStateResponse{
		Models:     r,
		Pagination: pageRes,
	}, nil
}

func (q GrpcQuerier) RawContractState(c context.Context, req *types.QueryRawContractStateRequest) (*types.QueryRawContractStateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	err := sdk.ValidateAccAddress(req.Address)
	if err != nil {
		return nil, err
	}

	if !q.keeper.HasContractInfo(ctx, sdk.AccAddress(req.Address)) {
		return nil, types.ErrNotFound
	}
	rsp := q.keeper.QueryRaw(ctx, sdk.AccAddress(req.Address), req.QueryData)
	return &types.QueryRawContractStateResponse{Data: rsp}, nil
}

func (q GrpcQuerier) SmartContractState(c context.Context, req *types.QuerySmartContractStateRequest) (rsp *types.QuerySmartContractStateResponse, err error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if err := req.QueryData.ValidateBasic(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid query data")
	}
	err = sdk.ValidateAccAddress(req.Address)
	if err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(c).WithGasMeter(sdk.NewGasMeter(q.queryGasLimit))
	// recover from out-of-gas panic
	defer func() {
		if r := recover(); r != nil {
			switch rType := r.(type) {
			case sdk.ErrorOutOfGas:
				err = sdkerrors.Wrapf(sdkerrors.ErrOutOfGas,
					"out of gas in location: %v; gasWanted: %d, gasUsed: %d",
					rType.Descriptor, ctx.GasMeter().Limit(), ctx.GasMeter().GasConsumed(),
				)
			default:
				err = sdkerrors.ErrPanic
			}
			rsp = nil
			moduleLogger(ctx).
				Debug("smart query contract",
					"error", "recovering panic",
					"contract-address", req.Address,
					"stacktrace", string(debug.Stack()))
		}
	}()

	bz, err := q.keeper.QuerySmart(ctx, sdk.AccAddress(req.Address), req.QueryData)
	switch {
	case err != nil:
		return nil, err
	case bz == nil:
		return nil, types.ErrNotFound
	}
	return &types.QuerySmartContractStateResponse{Data: bz}, nil

}

func (q GrpcQuerier) Code(c context.Context, req *types.QueryCodeRequest) (*types.QueryCodeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.CodeId == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "code id")
	}
	rsp, err := queryCode(sdk.UnwrapSDKContext(c), req.CodeId, q.keeper)
	switch {
	case err != nil:
		return nil, err
	case rsp == nil:
		return nil, types.ErrNotFound
	}
	return &types.QueryCodeResponse{
		CodeInfoResponse: rsp.CodeInfoResponse,
		Data:             rsp.Data,
	}, nil
}

func (q GrpcQuerier) Codes(c context.Context, req *types.QueryCodesRequest) (*types.QueryCodesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	r := make([]types.CodeInfoResponse, 0)
	prefixStore := prefix.NewStore(ctx.KVStore(q.storeKey), types.CodeKeyPrefix)
	pageRes, err := query.FilteredPaginate(prefixStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		if accumulate {
			var c types.CodeInfo
			if err := q.cdc.Unmarshal(value, &c); err != nil {
				return false, err
			}
			r = append(r, types.CodeInfoResponse{
				CodeID:                binary.BigEndian.Uint64(key),
				Creator:               c.Creator,
				DataHash:              c.CodeHash,
				InstantiatePermission: c.InstantiateConfig,
			})
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return &types.QueryCodesResponse{CodeInfos: r, Pagination: pageRes}, nil
}

func queryContractInfo(ctx sdk.Context, addr sdk.AccAddress, keeper types.ViewKeeper) (*types.QueryContractInfoResponse, error) {
	info := keeper.GetContractInfo(ctx, addr)
	if info == nil {
		return nil, types.ErrNotFound
	}
	// redact the Created field (just used for sorting, not part of public API)
	info.Created = nil
	return &types.QueryContractInfoResponse{
		Address:      addr.String(),
		ContractInfo: *info,
	}, nil
}

func queryCode(ctx sdk.Context, codeID uint64, keeper types.ViewKeeper) (*types.QueryCodeResponse, error) {
	if codeID == 0 {
		return nil, nil
	}
	res := keeper.GetCodeInfo(ctx, codeID)
	if res == nil {
		// nil, nil leads to 404 in rest handler
		return nil, nil
	}
	info := types.CodeInfoResponse{
		CodeID:                codeID,
		Creator:               res.Creator,
		DataHash:              res.CodeHash,
		InstantiatePermission: res.InstantiateConfig,
	}

	code, err := keeper.GetByteCode(ctx, codeID)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "loading wasm code")
	}

	return &types.QueryCodeResponse{CodeInfoResponse: &info, Data: code}, nil
}

func (q GrpcQuerier) PinnedCodes(c context.Context, req *types.QueryPinnedCodesRequest) (*types.QueryPinnedCodesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	r := make([]uint64, 0)

	prefixStore := prefix.NewStore(ctx.KVStore(q.storeKey), types.PinnedCodeIndexPrefix)
	pageRes, err := query.FilteredPaginate(prefixStore, req.Pagination, func(key []byte, _ []byte, accumulate bool) (bool, error) {
		if accumulate {

			r = append(r, sdk.BigEndianToUint64(key))
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return &types.QueryPinnedCodesResponse{
		CodeIDs:    r,
		Pagination: pageRes,
	}, nil

}
