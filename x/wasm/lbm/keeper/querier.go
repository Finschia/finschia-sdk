package keeper

import (
	"context"

	"github.com/line/lbm-sdk/codec"
	"github.com/line/lbm-sdk/store/prefix"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/query"
	lbmwasmtypes "github.com/line/lbm-sdk/x/wasm/lbm/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ lbmwasmtypes.QueryServer = &GrpcQuerier{}

type GrpcQuerier struct {
	cdc      codec.Codec
	storeKey sdk.StoreKey
	keeper   lbmwasmtypes.ViewKeeper
}

// NewGrpcQuerier constructor
func NewGrpcQuerier(cdc codec.Codec, storeKey sdk.StoreKey, keeper lbmwasmtypes.ViewKeeper) *GrpcQuerier {
	return &GrpcQuerier{cdc: cdc, storeKey: storeKey, keeper: keeper}
}

func (q GrpcQuerier) InactiveContracts(c context.Context, req *lbmwasmtypes.QueryInactiveContractsRequest) (*lbmwasmtypes.QueryInactiveContractsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	addresses := make([]string, 0)
	prefixStore := prefix.NewStore(ctx.KVStore(q.storeKey), inactiveContractPrefix)
	pageRes, err := query.FilteredPaginate(prefixStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		if accumulate {
			contractAddress := sdk.AccAddress(value)
			addresses = append(addresses, contractAddress.String())
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return &lbmwasmtypes.QueryInactiveContractsResponse{
		Addresses:  addresses,
		Pagination: pageRes,
	}, nil
}
