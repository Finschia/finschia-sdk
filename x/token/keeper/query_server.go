package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/line/lbm-sdk/store/prefix"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/query"
	"github.com/line/lbm-sdk/x/token"
	"github.com/line/lbm-sdk/x/token/class"
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

	if err := class.ValidateID(req.ContractId); err != nil {
		return nil, err
	}
	if err := sdk.ValidateAccAddress(req.Address); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	balance := s.keeper.GetBalance(ctx, req.ContractId, sdk.AccAddress(req.Address))

	return &token.QueryBalanceResponse{Amount: balance}, nil
}

// Supply queries the number of tokens from the given contract id.
func (s queryServer) Supply(c context.Context, req *token.QuerySupplyRequest) (*token.QuerySupplyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := class.ValidateID(req.ContractId); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	// daphne compat.
	if _, err := s.keeper.GetClass(ctx, req.ContractId); err != nil {
		return nil, err
	}
	supply := s.keeper.GetSupply(ctx, req.ContractId)

	return &token.QuerySupplyResponse{Amount: supply}, nil
}

// Minted queries the number of tokens from the given contract id.
func (s queryServer) Minted(c context.Context, req *token.QueryMintedRequest) (*token.QueryMintedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := class.ValidateID(req.ContractId); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	// daphne compat.
	if _, err := s.keeper.GetClass(ctx, req.ContractId); err != nil {
		return nil, err
	}
	minted := s.keeper.GetMinted(ctx, req.ContractId)

	return &token.QueryMintedResponse{Amount: minted}, nil
}

// Burnt queries the number of tokens from the given contract id.
func (s queryServer) Burnt(c context.Context, req *token.QueryBurntRequest) (*token.QueryBurntResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := class.ValidateID(req.ContractId); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	// daphne compat.
	if _, err := s.keeper.GetClass(ctx, req.ContractId); err != nil {
		return nil, err
	}
	burnt := s.keeper.GetBurnt(ctx, req.ContractId)

	return &token.QueryBurntResponse{Amount: burnt}, nil
}

// TokenClass queries an token metadata based on its contract id.
func (s queryServer) TokenClass(c context.Context, req *token.QueryTokenClassRequest) (*token.QueryTokenClassResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := class.ValidateID(req.ContractId); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	class, err := s.keeper.GetClass(ctx, req.ContractId)
	if err != nil {
		return nil, err
	}

	return &token.QueryTokenClassResponse{Class: class}, nil
}

// TokenClasses queries all token metadata.
func (s queryServer) TokenClasses(c context.Context, req *token.QueryTokenClassesRequest) (*token.QueryTokenClassesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(s.keeper.storeKey)
	classStore := prefix.NewStore(store, classKeyPrefix)
	var classes []token.TokenClass
	pageRes, err := query.Paginate(classStore, req.Pagination, func(key []byte, value []byte) error {
		var class token.TokenClass
		s.keeper.cdc.MustUnmarshal(value, &class)
		classes = append(classes, class)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &token.QueryTokenClassesResponse{Classes: classes, Pagination: pageRes}, nil
}

func (s queryServer) Grant(c context.Context, req *token.QueryGrantRequest) (*token.QueryGrantResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := class.ValidateID(req.ContractId); err != nil {
		return nil, err
	}
	if err := sdk.ValidateAccAddress(req.Grantee); err != nil {
		return nil, err
	}
	permission := token.Permission(token.Permission_value[req.Permission])
	if permission == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid permission")
	}

	ctx := sdk.UnwrapSDKContext(c)
	grant := s.keeper.GetGrant(ctx, req.ContractId, sdk.AccAddress(req.Grantee), permission)

	return &token.QueryGrantResponse{Grant: grant}, nil
}

func (s queryServer) GranteeGrants(c context.Context, req *token.QueryGranteeGrantsRequest) (*token.QueryGranteeGrantsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := class.ValidateID(req.ContractId); err != nil {
		return nil, err
	}
	if err := sdk.ValidateAccAddress(req.Grantee); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(s.keeper.storeKey)
	grantStore := prefix.NewStore(store, grantKeyPrefixByGrantee(req.ContractId, sdk.AccAddress(req.Grantee)))
	var grants []token.Grant
	pageRes, err := query.Paginate(grantStore, req.Pagination, func(key []byte, _ []byte) error {
		permission := token.Permission(key[0])
		grants = append(grants, token.Grant{
			Grantee:    req.Grantee,
			Permission: permission.String(),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &token.QueryGranteeGrantsResponse{Grants: grants, Pagination: pageRes}, nil
}

func (s queryServer) Authorization(c context.Context, req *token.QueryAuthorizationRequest) (*token.QueryAuthorizationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := class.ValidateID(req.ContractId); err != nil {
		return nil, err
	}
	if err := sdk.ValidateAccAddress(req.Proxy); err != nil {
		return nil, err
	}
	if err := sdk.ValidateAccAddress(req.Approver); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	authorization := s.keeper.GetAuthorization(ctx, req.ContractId, sdk.AccAddress(req.Approver), sdk.AccAddress(req.Proxy))

	return &token.QueryAuthorizationResponse{Authorization: authorization}, nil
}

func (s queryServer) OperatorAuthorizations(c context.Context, req *token.QueryOperatorAuthorizationsRequest) (*token.QueryOperatorAuthorizationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := class.ValidateID(req.ContractId); err != nil {
		return nil, err
	}
	if err := sdk.ValidateAccAddress(req.Proxy); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(s.keeper.storeKey)
	authorizationStore := prefix.NewStore(store, authorizationKeyPrefixByOperator(req.ContractId, sdk.AccAddress(req.Proxy)))
	var authorizations []token.Authorization
	pageRes, err := query.Paginate(authorizationStore, req.Pagination, func(key []byte, value []byte) error {
		approver := sdk.AccAddress(key)
		authorizations = append(authorizations, token.Authorization{
			Approver: approver.String(),
			Proxy:    req.Proxy,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &token.QueryOperatorAuthorizationsResponse{Authorizations: authorizations, Pagination: pageRes}, nil
}
