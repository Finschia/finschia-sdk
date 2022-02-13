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
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := sdk.ValidateAccAddress(req.Address); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	balance := s.keeper.GetBalance(ctx, sdk.AccAddress(req.Address), req.ClassId)

	return &token.QueryBalanceResponse{Amount: balance.Amount}, nil
}

// Supply queries the number of tokens from the given class id.
func (s queryServer) Supply(c context.Context, req *token.QuerySupplyRequest) (*token.QuerySupplyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	queriers := map[string]func(ctx sdk.Context, classId string) token.FT{
		"supply": s.keeper.GetSupply,
		"mint":   s.keeper.GetMint,
		"burn":   s.keeper.GetBurn,
	}
	querier, ok := queriers[req.Type]
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid supply type: %s", req.Type)
	}

	supply := querier(ctx, req.ClassId)

	return &token.QuerySupplyResponse{Amount: supply.Amount}, nil
}

// Token queries an token metadata based on its class id.
func (s queryServer) Token(c context.Context, req *token.QueryTokenRequest) (*token.QueryTokenResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
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
		return nil, status.Error(codes.InvalidArgument, "empty request")
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

func (s queryServer) Grants(c context.Context, req *token.QueryGrantsRequest) (*token.QueryGrantsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := sdk.ValidateAccAddress(req.Grantee); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	var grants []token.Grant
	actions := []string{"mint", "burn", "modify"}
	for _, action := range actions {
		granted := s.keeper.GetGrant(ctx, sdk.AccAddress(req.Grantee), req.ClassId, action)
		if granted {
			grants = append(grants, token.Grant{
				ClassId: req.ClassId,
				Grantee: req.Grantee,
				Action:  action,
			})
		}
	}

	return &token.QueryGrantsResponse{Grants: grants}, nil
}

func (s queryServer) Approve(c context.Context, req *token.QueryApproveRequest) (*token.QueryApproveResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := sdk.ValidateAccAddress(req.Proxy); err != nil {
		return nil, err
	}
	if err := sdk.ValidateAccAddress(req.Approver); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	approved := s.keeper.GetApprove(ctx, sdk.AccAddress(req.Approver), sdk.AccAddress(req.Proxy), req.ClassId)
	var approve *token.Approve
	if approved {
		approve = &token.Approve{
			ClassId:  req.ClassId,
			Approver: req.Approver,
			Proxy:    req.Proxy,
		}
	}

	return &token.QueryApproveResponse{Approve: approve}, nil
}

func (s queryServer) Approves(c context.Context, req *token.QueryApprovesRequest) (*token.QueryApprovesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := sdk.ValidateAccAddress(req.Proxy); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(s.keeper.storeKey)
	approveStore := prefix.NewStore(store, approveKeyPrefixByProxy(req.ClassId, sdk.AccAddress(req.Proxy)))
	var approves []token.Approve
	pageRes, err := query.Paginate(approveStore, req.Pagination, func(key []byte, value []byte) error {
		approver := string(key)
		approves = append(approves, token.Approve{
			ClassId:  req.ClassId,
			Approver: approver,
			Proxy:    req.Proxy,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &token.QueryApprovesResponse{Approves: approves, Pagination: pageRes}, nil
}
