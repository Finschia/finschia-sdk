package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Finschia/finschia-sdk/store/prefix"
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/types/query"
	"github.com/Finschia/finschia-sdk/x/token"
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

func (s queryServer) addressFromBech32GRPC(address, context string) (sdk.AccAddress, error) {
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress.Wrap(address), context).Error())
	}

	return addr, nil
}

// Balance queries the number of tokens of a given class owned by the owner.
func (s queryServer) Balance(c context.Context, req *token.QueryBalanceRequest) (*token.QueryBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := token.ValidateContractID(req.ContractId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	addr, err := s.addressFromBech32GRPC(req.Address, "address")
	if err != nil {
		return nil, err
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
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)
	supply := s.keeper.GetSupply(ctx, req.ContractId)

	return &token.QuerySupplyResponse{Amount: supply}, nil
}

// Minted queries the number of tokens from the given contract id.
func (s queryServer) Minted(c context.Context, req *token.QueryMintedRequest) (*token.QueryMintedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := token.ValidateContractID(req.ContractId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)
	minted := s.keeper.GetMinted(ctx, req.ContractId)

	return &token.QueryMintedResponse{Amount: minted}, nil
}

// Burnt queries the number of tokens from the given contract id.
func (s queryServer) Burnt(c context.Context, req *token.QueryBurntRequest) (*token.QueryBurntResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := token.ValidateContractID(req.ContractId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)
	burnt := s.keeper.GetBurnt(ctx, req.ContractId)

	return &token.QueryBurntResponse{Amount: burnt}, nil
}

// Contract queries an token metadata based on its contract id.
func (s queryServer) Contract(c context.Context, req *token.QueryContractRequest) (*token.QueryContractResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := token.ValidateContractID(req.ContractId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)
	class, err := s.keeper.GetClass(ctx, req.ContractId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &token.QueryContractResponse{Contract: *class}, nil
}

func (s queryServer) GranteeGrants(c context.Context, req *token.QueryGranteeGrantsRequest) (*token.QueryGranteeGrantsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := token.ValidateContractID(req.ContractId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	grantee, err := s.addressFromBech32GRPC(req.Grantee, "grantee")
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(s.keeper.storeKey)
	grantStore := prefix.NewStore(store, grantKeyPrefixByGrantee(req.ContractId, grantee))
	var grants []token.Grant
	pageRes, err := query.Paginate(grantStore, req.Pagination, func(key, _ []byte) error {
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

func (s queryServer) IsOperatorFor(c context.Context, req *token.QueryIsOperatorForRequest) (*token.QueryIsOperatorForResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := token.ValidateContractID(req.ContractId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	operator, err := s.addressFromBech32GRPC(req.Operator, "operator")
	if err != nil {
		return nil, err
	}
	holder, err := s.addressFromBech32GRPC(req.Holder, "holder")
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	_, err = s.keeper.GetAuthorization(ctx, req.ContractId, holder, operator)
	authorized := err == nil

	return &token.QueryIsOperatorForResponse{Authorized: authorized}, nil
}

func (s queryServer) HoldersByOperator(c context.Context, req *token.QueryHoldersByOperatorRequest) (*token.QueryHoldersByOperatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := token.ValidateContractID(req.ContractId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	operator, err := s.addressFromBech32GRPC(req.Operator, "operator")
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(s.keeper.storeKey)
	authorizationStore := prefix.NewStore(store, authorizationKeyPrefixByOperator(req.ContractId, operator))
	var holders []string
	pageRes, err := query.Paginate(authorizationStore, req.Pagination, func(key, value []byte) error {
		holder := sdk.AccAddress(key)
		holders = append(holders, holder.String())
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &token.QueryHoldersByOperatorResponse{Holders: holders, Pagination: pageRes}, nil
}
