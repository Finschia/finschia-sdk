package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/line/lbm-sdk/store/prefix"
	"github.com/line/lbm-sdk/types/query"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/consortium/types"
)

var _ types.QueryServer = Keeper{}

func (q Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{q.GetParams(ctx)}, nil
}

func (q Keeper) ValidatorAuths(c context.Context, req *types.QueryValidatorAuthsRequest) (*types.QueryValidatorAuthsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	auths := []*types.ValidatorAuth{}
	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(q.storeKey)
	validatorStore := prefix.NewStore(store, types.ValidatorAuthKeyPrefix)
	pageRes, err := query.Paginate(validatorStore, req.Pagination, func(key []byte, value []byte) error {
		var auth types.ValidatorAuth
		q.cdc.MustUnmarshalBinaryBare(value, &auth)
		auths = append(auths, &auth)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryValidatorAuthsResponse{Auths: auths, Pagination: pageRes}, nil
}

func (q Keeper) ValidatorAuth(c context.Context, req *types.QueryValidatorAuthRequest) (*types.QueryValidatorAuthResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ValidatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty validator address")
	}

	ctx := sdk.UnwrapSDKContext(c)

	addr := sdk.ValAddress(req.ValidatorAddress)
	auth, err := q.GetValidatorAuth(ctx, addr)
	if err != nil {
		return nil, err
	}

	return &types.QueryValidatorAuthResponse{Auth: auth}, nil
}
