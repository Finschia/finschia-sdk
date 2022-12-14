package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/line/lbm-sdk/store/prefix"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
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

	if err := token.ValidateContractID(req.ContractId); err != nil {
		return nil, token.SDKErrorToGRPCError(err)
	}
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, token.SDKErrorToGRPCError(sdkerrors.ErrInvalidAddress.Wrapf("invalid address: %s", req.Address))
	}

	ctx := sdk.UnwrapSDKContext(c)
	balance := s.keeper.GetBalance(ctx, req.ContractId, addr)

	return &token.QueryBalanceResponse{Amount: balance}, nil
}

// Supply queries the number of tokens from the given contract id.
func (s queryServer) Supply(c context.Context, req *token.QuerySupplyRequest) (*token.QuerySupplyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := token.ValidateContractID(req.ContractId); err != nil {
		return nil, token.SDKErrorToGRPCError(err)
	}

	ctx := sdk.UnwrapSDKContext(c)
	// daphne compat.
	if _, err := s.keeper.GetClass(ctx, req.ContractId); err != nil {
		return nil, token.SDKErrorToGRPCError(err)
	}
	supply := s.keeper.GetSupply(ctx, req.ContractId)

	return &token.QuerySupplyResponse{Amount: supply}, nil
}

// Minted queries the number of tokens from the given contract id.
func (s queryServer) Minted(c context.Context, req *token.QueryMintedRequest) (*token.QueryMintedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := token.ValidateContractID(req.ContractId); err != nil {
		return nil, token.SDKErrorToGRPCError(err)
	}

	ctx := sdk.UnwrapSDKContext(c)
	// daphne compat.
	if _, err := s.keeper.GetClass(ctx, req.ContractId); err != nil {
		return nil, token.SDKErrorToGRPCError(err)
	}
	minted := s.keeper.GetMinted(ctx, req.ContractId)

	return &token.QueryMintedResponse{Amount: minted}, nil
}

// Burnt queries the number of tokens from the given contract id.
func (s queryServer) Burnt(c context.Context, req *token.QueryBurntRequest) (*token.QueryBurntResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := token.ValidateContractID(req.ContractId); err != nil {
		return nil, token.SDKErrorToGRPCError(err)
	}

	ctx := sdk.UnwrapSDKContext(c)
	// daphne compat.
	if _, err := s.keeper.GetClass(ctx, req.ContractId); err != nil {
		return nil, token.SDKErrorToGRPCError(err)
	}
	burnt := s.keeper.GetBurnt(ctx, req.ContractId)

	return &token.QueryBurntResponse{Amount: burnt}, nil
}

// TokenClass queries an token metadata based on its contract id.
func (s queryServer) TokenClass(c context.Context, req *token.QueryTokenClassRequest) (*token.QueryTokenClassResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := token.ValidateContractID(req.ContractId); err != nil {
		return nil, token.SDKErrorToGRPCError(err)
	}

	ctx := sdk.UnwrapSDKContext(c)
	class, err := s.keeper.GetClass(ctx, req.ContractId)
	if err != nil {
		return nil, token.SDKErrorToGRPCError(err)
	}

	return &token.QueryTokenClassResponse{Class: *class}, nil
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

func (s queryServer) GranteeGrants(c context.Context, req *token.QueryGranteeGrantsRequest) (*token.QueryGranteeGrantsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := token.ValidateContractID(req.ContractId); err != nil {
		return nil, token.SDKErrorToGRPCError(err)
	}
	grantee, err := sdk.AccAddressFromBech32(req.Grantee)
	if err != nil {
		return nil, token.SDKErrorToGRPCError(sdkerrors.ErrInvalidAddress.Wrapf("invalid grantee address: %s", req.Grantee))
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(s.keeper.storeKey)
	grantStore := prefix.NewStore(store, grantKeyPrefixByGrantee(req.ContractId, grantee))
	var grants []token.Grant
	pageRes, err := query.Paginate(grantStore, req.Pagination, func(key []byte, _ []byte) error {
		permission := token.Permission(key[0])
		grants = append(grants, token.Grant{
			Grantee:    req.Grantee,
			Permission: permission,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &token.QueryGranteeGrantsResponse{Grants: grants, Pagination: pageRes}, nil
}

func (s queryServer) Approved(c context.Context, req *token.QueryApprovedRequest) (*token.QueryApprovedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := token.ValidateContractID(req.ContractId); err != nil {
		return nil, token.SDKErrorToGRPCError(err)
	}
	proxy, err := sdk.AccAddressFromBech32(req.Proxy)
	if err != nil {
		return nil, token.SDKErrorToGRPCError(sdkerrors.ErrInvalidAddress.Wrapf("invalid proxy address: %s", req.Proxy))
	}
	approver, err := sdk.AccAddressFromBech32(req.Approver)
	if err != nil {
		return nil, token.SDKErrorToGRPCError(sdkerrors.ErrInvalidAddress.Wrapf("invalid approver address: %s", req.Approver))
	}

	ctx := sdk.UnwrapSDKContext(c)
	_, err = s.keeper.GetAuthorization(ctx, req.ContractId, approver, proxy)
	approved := err == nil

	return &token.QueryApprovedResponse{Approved: approved}, nil
}

func (s queryServer) Approvers(c context.Context, req *token.QueryApproversRequest) (*token.QueryApproversResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := token.ValidateContractID(req.ContractId); err != nil {
		return nil, token.SDKErrorToGRPCError(err)
	}
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, token.SDKErrorToGRPCError(sdkerrors.ErrInvalidAddress.Wrapf("invalid address: %s", req.Address))
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(s.keeper.storeKey)
	authorizationStore := prefix.NewStore(store, authorizationKeyPrefixByOperator(req.ContractId, addr))
	var approvers []string
	pageRes, err := query.Paginate(authorizationStore, req.Pagination, func(key []byte, value []byte) error {
		holder := sdk.AccAddress(key)
		approvers = append(approvers, holder.String())
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &token.QueryApproversResponse{Approvers: approvers, Pagination: pageRes}, nil
}
