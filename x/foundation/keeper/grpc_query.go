package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/line/lbm-sdk/store/prefix"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/query"
	"github.com/line/lbm-sdk/x/foundation"
)

type queryServer struct {
	keeper Keeper
}

func NewQueryServer(keeper Keeper) foundation.QueryServer {
	return &queryServer{
		keeper: keeper,
	}
}

var _ foundation.QueryServer = (*queryServer)(nil)

func (s queryServer) Params(c context.Context, req *foundation.QueryParamsRequest) (*foundation.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	return &foundation.QueryParamsResponse{Params: s.keeper.GetParams(ctx)}, nil
}

func (s queryServer) ValidatorAuth(c context.Context, req *foundation.QueryValidatorAuthRequest) (*foundation.QueryValidatorAuthResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ValidatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty validator address")
	}

	ctx := sdk.UnwrapSDKContext(c)

	addr := sdk.ValAddress(req.ValidatorAddress)
	auth, err := s.keeper.GetValidatorAuth(ctx, addr)
	if err != nil {
		return nil, err
	}

	return &foundation.QueryValidatorAuthResponse{Auth: auth}, nil
}

func (s queryServer) ValidatorAuths(c context.Context, req *foundation.QueryValidatorAuthsRequest) (*foundation.QueryValidatorAuthsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var auths []*foundation.ValidatorAuth
	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(s.keeper.storeKey)
	validatorStore := prefix.NewStore(store, validatorAuthKeyPrefix)
	pageRes, err := query.Paginate(validatorStore, req.Pagination, func(key []byte, value []byte) error {
		var auth foundation.ValidatorAuth
		s.keeper.cdc.MustUnmarshal(value, &auth)
		auths = append(auths, &auth)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &foundation.QueryValidatorAuthsResponse{Auths: auths, Pagination: pageRes}, nil
}
