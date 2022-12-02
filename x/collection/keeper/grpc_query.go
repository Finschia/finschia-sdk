package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gogo/protobuf/proto"

	codectypes "github.com/line/lbm-sdk/codec/types"
	"github.com/line/lbm-sdk/store/prefix"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/types/query"
	"github.com/line/lbm-sdk/x/collection"
)

type queryServer struct {
	keeper Keeper
}

// NewQueryServer returns an implementation of the token QueryServer interface
// for the provided Keeper.
func NewQueryServer(keeper Keeper) collection.QueryServer {
	return &queryServer{
		keeper: keeper,
	}
}

var _ collection.QueryServer = queryServer{}

// Balance queries the number of tokens of a given token id owned by the owner.
func (s queryServer) Balance(c context.Context, req *collection.QueryBalanceRequest) (*collection.QueryBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, collection.SDKErrorToGRPCError(sdkerrors.ErrInvalidAddress.Wrapf("invalid address: %s", req.Address))
	}

	if err := collection.ValidateTokenID(req.TokenId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	ctx := sdk.UnwrapSDKContext(c)
	balance := s.keeper.GetBalance(ctx, req.ContractId, addr, req.TokenId)
	coin := collection.NewCoin(req.TokenId, balance)

	return &collection.QueryBalanceResponse{Balance: coin}, nil
}

// AllBalances queries all tokens owned by owner.
func (s queryServer) AllBalances(c context.Context, req *collection.QueryAllBalancesRequest) (*collection.QueryAllBalancesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, collection.SDKErrorToGRPCError(sdkerrors.ErrInvalidAddress.Wrapf("invalid address: %s", req.Address))
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(s.keeper.storeKey)
	balanceStore := prefix.NewStore(store, balanceKeyPrefixByAddress(req.ContractId, addr))
	var balances []collection.Coin
	pageRes, err := query.Paginate(balanceStore, req.Pagination, func(key []byte, value []byte) error {
		tokenID := string(key)

		var balance sdk.Int
		if err := balance.Unmarshal(value); err != nil {
			panic(err)
		}

		coin := collection.NewCoin(tokenID, balance)
		balances = append(balances, coin)

		if err := collection.ValidateNFTID(tokenID); err == nil {
			s.keeper.iterateDescendants(ctx, req.ContractId, tokenID, func(tokenID string, _ int) (stop bool) {
				coin := collection.NewCoin(tokenID, sdk.OneInt())
				balances = append(balances, coin)

				return false
			})
		}

		return nil
	})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &collection.QueryAllBalancesResponse{Balances: balances, Pagination: pageRes}, nil
}

func (s queryServer) FTSupply(c context.Context, req *collection.QueryFTSupplyRequest) (*collection.QueryFTSupplyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	if err := collection.ValidateTokenID(req.TokenId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	classID := collection.SplitTokenID(req.TokenId)

	ctx := sdk.UnwrapSDKContext(c)
	if _, err := s.keeper.GetTokenClass(ctx, req.ContractId, classID); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}
	supply := s.keeper.GetSupply(ctx, req.ContractId, classID)

	return &collection.QueryFTSupplyResponse{Supply: supply}, nil
}

func (s queryServer) FTMinted(c context.Context, req *collection.QueryFTMintedRequest) (*collection.QueryFTMintedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	if err := collection.ValidateTokenID(req.TokenId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	classID := collection.SplitTokenID(req.TokenId)

	ctx := sdk.UnwrapSDKContext(c)
	if _, err := s.keeper.GetTokenClass(ctx, req.ContractId, classID); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}
	minted := s.keeper.GetMinted(ctx, req.ContractId, classID)

	return &collection.QueryFTMintedResponse{Minted: minted}, nil
}

func (s queryServer) FTBurnt(c context.Context, req *collection.QueryFTBurntRequest) (*collection.QueryFTBurntResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	if err := collection.ValidateTokenID(req.TokenId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	classID := collection.SplitTokenID(req.TokenId)

	ctx := sdk.UnwrapSDKContext(c)
	if _, err := s.keeper.GetTokenClass(ctx, req.ContractId, classID); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}
	burnt := s.keeper.GetBurnt(ctx, req.ContractId, classID)

	return &collection.QueryFTBurntResponse{Burnt: burnt}, nil
}

func (s queryServer) NFTSupply(c context.Context, req *collection.QueryNFTSupplyRequest) (*collection.QueryNFTSupplyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	classID := req.TokenType
	if err := collection.ValidateClassID(classID); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	ctx := sdk.UnwrapSDKContext(c)
	if _, err := s.keeper.GetTokenClass(ctx, req.ContractId, classID); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}
	supply := s.keeper.GetSupply(ctx, req.ContractId, classID)

	return &collection.QueryNFTSupplyResponse{Supply: supply}, nil
}

func (s queryServer) NFTMinted(c context.Context, req *collection.QueryNFTMintedRequest) (*collection.QueryNFTMintedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	classID := req.TokenType
	if err := collection.ValidateClassID(classID); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	ctx := sdk.UnwrapSDKContext(c)
	if _, err := s.keeper.GetTokenClass(ctx, req.ContractId, classID); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}
	minted := s.keeper.GetMinted(ctx, req.ContractId, classID)

	return &collection.QueryNFTMintedResponse{Minted: minted}, nil
}

func (s queryServer) NFTBurnt(c context.Context, req *collection.QueryNFTBurntRequest) (*collection.QueryNFTBurntResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	classID := req.TokenType
	if err := collection.ValidateClassID(classID); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	ctx := sdk.UnwrapSDKContext(c)
	if _, err := s.keeper.GetTokenClass(ctx, req.ContractId, classID); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}
	burnt := s.keeper.GetBurnt(ctx, req.ContractId, classID)

	return &collection.QueryNFTBurntResponse{Burnt: burnt}, nil
}

func (s queryServer) Contract(c context.Context, req *collection.QueryContractRequest) (*collection.QueryContractResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	ctx := sdk.UnwrapSDKContext(c)
	contract, err := s.keeper.GetContract(ctx, req.ContractId)
	if err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	return &collection.QueryContractResponse{Contract: *contract}, nil
}

// TokenClassTypeName queries the fully qualified message type name of a token class based on its class id.
func (s queryServer) TokenClassTypeName(c context.Context, req *collection.QueryTokenClassTypeNameRequest) (*collection.QueryTokenClassTypeNameResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	if err := collection.ValidateClassID(req.ClassId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	ctx := sdk.UnwrapSDKContext(c)
	class, err := s.keeper.GetTokenClass(ctx, req.ContractId, req.ClassId)
	if err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}
	name := proto.MessageName(class)

	return &collection.QueryTokenClassTypeNameResponse{Name: name}, nil
}

func (s queryServer) TokenType(c context.Context, req *collection.QueryTokenTypeRequest) (*collection.QueryTokenTypeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	classID := req.TokenType
	if err := collection.ValidateClassID(classID); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	ctx := sdk.UnwrapSDKContext(c)
	class, err := s.keeper.GetTokenClass(ctx, req.ContractId, classID)
	if err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	nftClass, ok := class.(*collection.NFTClass)
	if !ok {
		return nil, collection.SDKErrorToGRPCError(sdkerrors.ErrInvalidType.Wrapf("not a class of non-fungible token: %s", classID))
	}

	tokenType := collection.TokenType{
		ContractId: req.ContractId,
		TokenType:  nftClass.Id,
		Name:       nftClass.Name,
		Meta:       nftClass.Meta,
	}

	return &collection.QueryTokenTypeResponse{TokenType: tokenType}, nil
}

func (s queryServer) TokenTypes(c context.Context, req *collection.QueryTokenTypesRequest) (*collection.QueryTokenTypesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(s.keeper.storeKey)
	tokenTypeStore := prefix.NewStore(store, legacyTokenTypeKeyPrefixByContractID(req.ContractId))
	var tokenTypes []collection.TokenType
	pageRes, err := query.Paginate(tokenTypeStore, req.Pagination, func(key []byte, value []byte) error {
		classID := string(key)
		class, err := s.keeper.GetTokenClass(ctx, req.ContractId, classID)
		if err != nil {
			panic(err)
		}

		nftClass, ok := class.(*collection.NFTClass)
		if !ok {
			panic(sdkerrors.ErrInvalidType.Wrapf("not a class of non-fungible token: %s", key))
		}

		tokenType := collection.TokenType{
			ContractId: req.ContractId,
			TokenType:  nftClass.Id,
			Name:       nftClass.Name,
			Meta:       nftClass.Meta,
		}
		tokenTypes = append(tokenTypes, tokenType)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &collection.QueryTokenTypesResponse{TokenTypes: tokenTypes, Pagination: pageRes}, nil
}

func (s queryServer) getToken(ctx sdk.Context, contractID string, tokenID string) (collection.Token, error) {
	switch {
	case collection.ValidateNFTID(tokenID) == nil:
		token, err := s.keeper.GetNFT(ctx, contractID, tokenID)
		if err != nil {
			return nil, err
		}

		owner := s.keeper.GetRootOwner(ctx, contractID, token.Id)
		return &collection.OwnerNFT{
			ContractId: contractID,
			TokenId:    token.Id,
			Name:       token.Name,
			Meta:       token.Meta,
			Owner:      owner.String(),
		}, nil
	case collection.ValidateFTID(tokenID) == nil:
		classID := collection.SplitTokenID(tokenID)
		class, err := s.keeper.GetTokenClass(ctx, contractID, classID)
		if err != nil {
			return nil, err
		}

		ftClass, ok := class.(*collection.FTClass)
		if !ok {
			panic(sdkerrors.ErrInvalidType.Wrapf("not a class of fungible token: %s", classID))
		}

		return &collection.FT{
			ContractId: contractID,
			TokenId:    collection.NewFTID(ftClass.Id),
			Name:       ftClass.Name,
			Meta:       ftClass.Meta,
			Decimals:   ftClass.Decimals,
			Mintable:   ftClass.Mintable,
		}, nil
	default:
		panic("cannot reach here: token must be ft or nft")
	}
}

func (s queryServer) Token(c context.Context, req *collection.QueryTokenRequest) (*collection.QueryTokenResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	if err := collection.ValidateTokenID(req.TokenId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	ctx := sdk.UnwrapSDKContext(c)
	legacyToken, err := s.getToken(ctx, req.ContractId, req.TokenId)
	if err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	any, err := codectypes.NewAnyWithValue(legacyToken)
	if err != nil {
		panic(err)
	}

	return &collection.QueryTokenResponse{Token: *any}, nil
}

func (s queryServer) TokensWithTokenType(c context.Context, req *collection.QueryTokensWithTokenTypeRequest) (*collection.QueryTokensWithTokenTypeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	if err := collection.ValidateClassID(req.TokenType); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(s.keeper.storeKey)
	tokenStore := prefix.NewStore(store, legacyTokenKeyPrefixByTokenType(req.ContractId, req.TokenType))
	var tokens []codectypes.Any
	pageRes, err := query.Paginate(tokenStore, req.Pagination, func(key []byte, value []byte) error {
		tokenID := req.TokenType + string(key)
		legacyToken, err := s.getToken(ctx, req.ContractId, tokenID)
		if err != nil {
			panic(err)
		}

		any, err := codectypes.NewAnyWithValue(legacyToken)
		if err != nil {
			panic(err)
		}

		tokens = append(tokens, *any)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &collection.QueryTokensWithTokenTypeResponse{Tokens: tokens, Pagination: pageRes}, nil
}

func (s queryServer) Tokens(c context.Context, req *collection.QueryTokensRequest) (*collection.QueryTokensResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(s.keeper.storeKey)
	tokenStore := prefix.NewStore(store, legacyTokenKeyPrefixByContractID(req.ContractId))
	var tokens []codectypes.Any
	pageRes, err := query.Paginate(tokenStore, req.Pagination, func(key []byte, value []byte) error {
		tokenID := string(key)
		legacyToken, err := s.getToken(ctx, req.ContractId, tokenID)
		if err != nil {
			panic(err)
		}

		any, err := codectypes.NewAnyWithValue(legacyToken)
		if err != nil {
			panic(err)
		}

		tokens = append(tokens, *any)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &collection.QueryTokensResponse{Tokens: tokens, Pagination: pageRes}, nil
}

func (s queryServer) Root(c context.Context, req *collection.QueryRootRequest) (*collection.QueryRootResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	if err := collection.ValidateNFTID(req.TokenId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	ctx := sdk.UnwrapSDKContext(c)
	if err := s.keeper.hasNFT(ctx, req.ContractId, req.TokenId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	root := s.keeper.GetRoot(ctx, req.ContractId, req.TokenId)
	token, err := s.keeper.GetNFT(ctx, req.ContractId, root)
	if err != nil {
		panic(err)
	}

	return &collection.QueryRootResponse{Root: *token}, nil
}

func (s queryServer) Parent(c context.Context, req *collection.QueryParentRequest) (*collection.QueryParentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	if err := collection.ValidateNFTID(req.TokenId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	ctx := sdk.UnwrapSDKContext(c)
	if err := s.keeper.hasNFT(ctx, req.ContractId, req.TokenId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	parent, err := s.keeper.GetParent(ctx, req.ContractId, req.TokenId)
	if err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	token, err := s.keeper.GetNFT(ctx, req.ContractId, *parent)
	if err != nil {
		panic(err)
	}

	return &collection.QueryParentResponse{Parent: *token}, nil
}

func (s queryServer) Children(c context.Context, req *collection.QueryChildrenRequest) (*collection.QueryChildrenResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	if err := collection.ValidateNFTID(req.TokenId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	ctx := sdk.UnwrapSDKContext(c)
	if err := s.keeper.hasNFT(ctx, req.ContractId, req.TokenId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	store := ctx.KVStore(s.keeper.storeKey)
	childStore := prefix.NewStore(store, childKeyPrefixByTokenID(req.ContractId, req.TokenId))
	var children []collection.NFT
	pageRes, err := query.Paginate(childStore, req.Pagination, func(key []byte, _ []byte) error {
		childID := string(key)
		child, err := s.keeper.GetNFT(ctx, req.ContractId, childID)
		if err != nil {
			panic(err)
		}

		children = append(children, *child)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &collection.QueryChildrenResponse{Children: children, Pagination: pageRes}, nil
}

func (s queryServer) GranteeGrants(c context.Context, req *collection.QueryGranteeGrantsRequest) (*collection.QueryGranteeGrantsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	granteeAddr, err := sdk.AccAddressFromBech32(req.Grantee)
	if err != nil {
		return nil, collection.SDKErrorToGRPCError(sdkerrors.ErrInvalidAddress.Wrapf("invalid grantee address: %s", req.Grantee))
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(s.keeper.storeKey)
	grantStore := prefix.NewStore(store, grantKeyPrefixByGrantee(req.ContractId, granteeAddr))
	var grants []collection.Grant
	pageRes, err := query.Paginate(grantStore, req.Pagination, func(key []byte, _ []byte) error {
		permission := collection.Permission(key[0])
		grants = append(grants, collection.Grant{
			Grantee:    req.Grantee,
			Permission: permission,
		})
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &collection.QueryGranteeGrantsResponse{Grants: grants, Pagination: pageRes}, nil
}

func (s queryServer) Approved(c context.Context, req *collection.QueryApprovedRequest) (*collection.QueryApprovedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, collection.SDKErrorToGRPCError(sdkerrors.ErrInvalidAddress.Wrapf("invalid address address: %s", req.Address))
	}
	approverAddr, err := sdk.AccAddressFromBech32(req.Approver)
	if err != nil {
		return nil, collection.SDKErrorToGRPCError(sdkerrors.ErrInvalidAddress.Wrapf("invalid approver address: %s", req.Approver))
	}

	ctx := sdk.UnwrapSDKContext(c)
	_, err = s.keeper.GetAuthorization(ctx, req.ContractId, approverAddr, addr)
	approved := (err == nil)

	return &collection.QueryApprovedResponse{Approved: approved}, nil
}

func (s queryServer) Approvers(c context.Context, req *collection.QueryApproversRequest) (*collection.QueryApproversResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, collection.SDKErrorToGRPCError(err)
	}

	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, collection.SDKErrorToGRPCError(sdkerrors.ErrInvalidAddress.Wrapf("invalid address address: %s", req.Address))
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
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &collection.QueryApproversResponse{Approvers: approvers, Pagination: pageRes}, nil
}
