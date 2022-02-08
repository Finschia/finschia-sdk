package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/line/lbm-sdk/store/prefix"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/query"
	"github.com/line/lbm-sdk/x/token"
)

type queryServer struct {
	keeper Keeper
}

// NewQueryServer returns an implementation of the token QueryServer interface
// for the provided Keeper.
func NewQueryServer(keeper Keeper) token.QueryServer {
	return &queryServer{
		keeper: keeper,
	}
}

var _ token.QueryServer = queryServer{}

// Balance queries the number of tokens of a given class owned by the owner.
func (s queryServer) Balance(c context.Context, req *token.QueryBalanceRequest) (*token.QueryBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	balance := s.keeper.GetBalance(ctx, sdk.AccAddress(req.Address), req.ClassId)

	return &token.QueryBalanceResponse{Amount: balance.Amount}, nil
}

// Supply queries the number of tokens from the given class id.
func (s queryServer) Supply(c context.Context, req *token.QuerySupplyRequest) (*token.QuerySupplyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	supply := s.keeper.GetSupply(ctx, req.ClassId)

	return &token.QuerySupplyResponse{Amount: supply.Amount}, nil
}

// Token queries an token metadata based on its class id.
func (s queryServer) Token(c context.Context, req *token.QueryTokenRequest) (*token.QueryTokenResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	class, err := s.keeper.GetClass(ctx, req.ClassId)
	if err != nil {
		return nil, err
	}

	return &token.QueryTokenResponse{Token: *class}, nil
}

// Tokens queries all token metadata.
func (s queryServer) Tokens(c context.Context, req *token.QueryTokensRequest) (*token.QueryTokensResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(s.keeper.storeKey)
	classStore := prefix.NewStore(store, classKeyPrefix)
	var classes []token.Token
	pageRes, err := query.Paginate(classStore, req.Pagination, func(key []byte, value []byte) error {
		var class token.Token
		s.keeper.cdc.MustUnmarshalBinaryBare(value, &class)
		classes = append(classes, class)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &token.QueryTokensResponse{Tokens: classes, Pagination: pageRes}, nil
}
