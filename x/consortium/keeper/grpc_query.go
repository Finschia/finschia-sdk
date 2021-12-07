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

func (q Keeper) Enabled(c context.Context, req *types.QueryEnabledRequest) (*types.QueryEnabledResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	enabled := q.GetEnabled(ctx)
	return &types.QueryEnabledResponse{Enabled: enabled}, nil
}

func (q Keeper) AllowedValidators(c context.Context, req *types.QueryAllowedValidatorsRequest) (*types.QueryAllowedValidatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	addrs := []string{}
	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(q.storeKey)
	validatorStore := prefix.NewStore(store, types.AllowedValidatorKeyPrefix)
	pageRes, err := query.Paginate(validatorStore, req.Pagination, func(key []byte, value []byte) error {
		addr := string(key)
		addrs = append(addrs, addr)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryAllowedValidatorsResponse{ValidatorAddresses: addrs, Pagination: pageRes}, nil
}

func (q Keeper) AllowedValidator(c context.Context, req *types.QueryAllowedValidatorRequest) (*types.QueryAllowedValidatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ValidatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty validator address")
	}

	ctx := sdk.UnwrapSDKContext(c)

	addr := sdk.ValAddress(req.ValidatorAddress)
	allowed := q.GetAllowedValidator(ctx, addr)

	return &types.QueryAllowedValidatorResponse{Allowed: allowed}, nil
}
