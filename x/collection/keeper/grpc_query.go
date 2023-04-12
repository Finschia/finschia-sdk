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

func (s queryServer) validateExistenceOfCollectionGRPC(ctx sdk.Context, id string) error {
	if _, err := s.keeper.GetContract(ctx, id); err != nil {
		return status.Error(codes.NotFound, err.Error())
	}

	return nil
}

func (s queryServer) validateExistenceOfFTClassGRPC(ctx sdk.Context, contractID, classID string) error {
	class, err := s.keeper.GetTokenClass(ctx, contractID, classID)
	if err != nil {
		return status.Error(codes.NotFound, err.Error())
	}

	_, ok := class.(*collection.FTClass)
	if !ok {
		return status.Error(codes.NotFound, sdkerrors.ErrInvalidType.Wrapf("not a class of fungible token: %s", classID).Error())
	}
	return nil
}

func (s queryServer) validateExistenceOfNFTClassGRPC(ctx sdk.Context, contractID, classID string) error {
	class, err := s.keeper.GetTokenClass(ctx, contractID, classID)
	if err != nil {
		return status.Error(codes.NotFound, err.Error())
	}

	_, ok := class.(*collection.NFTClass)
	if !ok {
		return status.Error(codes.NotFound, sdkerrors.ErrInvalidType.Wrapf("not a class of non-fungible token: %s", classID).Error())
	}
	return nil
}

var _ collection.QueryServer = queryServer{}

// Balance queries the number of tokens of a given token id owned by the owner.
func (s queryServer) Balance(c context.Context, req *collection.QueryBalanceRequest) (*collection.QueryBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid address: %s", req.Address)
	}

	if err := collection.ValidateTokenID(req.TokenId); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	balance := s.keeper.GetBalance(ctx, req.ContractId, addr, req.TokenId)
	coin := collection.Coin{
		TokenId: req.TokenId,
		Amount:  balance,
	}

	return &collection.QueryBalanceResponse{Balance: coin}, nil
}

// AllBalances queries all tokens owned by owner.
func (s queryServer) AllBalances(c context.Context, req *collection.QueryAllBalancesRequest) (*collection.QueryAllBalancesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid address: %s", req.Address)
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
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &collection.QueryAllBalancesResponse{Balances: balances, Pagination: pageRes}, nil
}

func (s queryServer) FTSupply(c context.Context, req *collection.QueryFTSupplyRequest) (*collection.QueryFTSupplyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	if err := collection.ValidateTokenID(req.TokenId); err != nil {
		return nil, err
	}

	classID := collection.SplitTokenID(req.TokenId)

	ctx := sdk.UnwrapSDKContext(c)

	if err := s.validateExistenceOfCollectionGRPC(ctx, req.ContractId); err != nil {
		return nil, err
	}

	if err := s.validateExistenceOfFTClassGRPC(ctx, req.ContractId, classID); err != nil {
		return nil, err
	}

	supply := s.keeper.GetSupply(ctx, req.ContractId, classID)

	return &collection.QueryFTSupplyResponse{Supply: supply}, nil
}

func (s queryServer) FTMinted(c context.Context, req *collection.QueryFTMintedRequest) (*collection.QueryFTMintedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	if err := collection.ValidateTokenID(req.TokenId); err != nil {
		return nil, err
	}

	classID := collection.SplitTokenID(req.TokenId)

	ctx := sdk.UnwrapSDKContext(c)

	if err := s.validateExistenceOfCollectionGRPC(ctx, req.ContractId); err != nil {
		return nil, err
	}

	if err := s.validateExistenceOfFTClassGRPC(ctx, req.ContractId, classID); err != nil {
		return nil, err
	}

	minted := s.keeper.GetMinted(ctx, req.ContractId, classID)

	return &collection.QueryFTMintedResponse{Minted: minted}, nil
}

func (s queryServer) FTBurnt(c context.Context, req *collection.QueryFTBurntRequest) (*collection.QueryFTBurntResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	if err := collection.ValidateTokenID(req.TokenId); err != nil {
		return nil, err
	}

	classID := collection.SplitTokenID(req.TokenId)

	ctx := sdk.UnwrapSDKContext(c)

	if err := s.validateExistenceOfCollectionGRPC(ctx, req.ContractId); err != nil {
		return nil, err
	}

	if err := s.validateExistenceOfFTClassGRPC(ctx, req.ContractId, classID); err != nil {
		return nil, err
	}

	burnt := s.keeper.GetBurnt(ctx, req.ContractId, classID)

	return &collection.QueryFTBurntResponse{Burnt: burnt}, nil
}

func (s queryServer) NFTSupply(c context.Context, req *collection.QueryNFTSupplyRequest) (*collection.QueryNFTSupplyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	classID := req.TokenType
	if err := collection.ValidateClassID(classID); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := s.validateExistenceOfCollectionGRPC(ctx, req.ContractId); err != nil {
		return nil, err
	}

	if err := s.validateExistenceOfNFTClassGRPC(ctx, req.ContractId, classID); err != nil {
		return nil, err
	}

	supply := s.keeper.GetSupply(ctx, req.ContractId, classID)

	return &collection.QueryNFTSupplyResponse{Supply: supply}, nil
}

func (s queryServer) NFTMinted(c context.Context, req *collection.QueryNFTMintedRequest) (*collection.QueryNFTMintedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	classID := req.TokenType
	if err := collection.ValidateClassID(classID); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := s.validateExistenceOfCollectionGRPC(ctx, req.ContractId); err != nil {
		return nil, err
	}

	if err := s.validateExistenceOfNFTClassGRPC(ctx, req.ContractId, classID); err != nil {
		return nil, err
	}

	minted := s.keeper.GetMinted(ctx, req.ContractId, classID)

	return &collection.QueryNFTMintedResponse{Minted: minted}, nil
}

func (s queryServer) NFTBurnt(c context.Context, req *collection.QueryNFTBurntRequest) (*collection.QueryNFTBurntResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	classID := req.TokenType
	if err := collection.ValidateClassID(classID); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := s.validateExistenceOfCollectionGRPC(ctx, req.ContractId); err != nil {
		return nil, err
	}

	if err := s.validateExistenceOfNFTClassGRPC(ctx, req.ContractId, classID); err != nil {
		return nil, err
	}

	burnt := s.keeper.GetBurnt(ctx, req.ContractId, classID)

	return &collection.QueryNFTBurntResponse{Burnt: burnt}, nil
}

func (s queryServer) Contract(c context.Context, req *collection.QueryContractRequest) (*collection.QueryContractResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	contract, err := s.keeper.GetContract(ctx, req.ContractId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &collection.QueryContractResponse{Contract: *contract}, nil
}

// TokenClassTypeName queries the fully qualified message type name of a token class based on its class id.
func (s queryServer) TokenClassTypeName(c context.Context, req *collection.QueryTokenClassTypeNameRequest) (*collection.QueryTokenClassTypeNameResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	if err := collection.ValidateClassID(req.ClassId); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := s.validateExistenceOfCollectionGRPC(ctx, req.ContractId); err != nil {
		return nil, err
	}

	class, err := s.keeper.GetTokenClass(ctx, req.ContractId, req.ClassId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	name := proto.MessageName(class)

	return &collection.QueryTokenClassTypeNameResponse{Name: name}, nil
}

func (s queryServer) TokenType(c context.Context, req *collection.QueryTokenTypeRequest) (*collection.QueryTokenTypeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	classID := req.TokenType
	if err := collection.ValidateClassID(classID); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	class, err := s.keeper.GetTokenClass(ctx, req.ContractId, classID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	nftClass, ok := class.(*collection.NFTClass)
	if !ok {
		return nil, status.Error(codes.NotFound, sdkerrors.ErrInvalidType.Wrapf("not a class of non-fungible token: %s", classID).Error())
	}

	tokenType := collection.TokenType{
		ContractId: req.ContractId,
		TokenType:  nftClass.Id,
		Name:       nftClass.Name,
		Meta:       nftClass.Meta,
	}

	return &collection.QueryTokenTypeResponse{TokenType: tokenType}, nil
}

func (s queryServer) getToken(ctx sdk.Context, contractID string, tokenID string) (collection.Token, error) {
	switch {
	case collection.ValidateNFTID(tokenID) == nil:
		token, err := s.keeper.GetNFT(ctx, contractID, tokenID)
		if err != nil {
			return nil, err
		}

		owner := s.keeper.GetRootOwner(ctx, contractID, token.TokenId)
		return &collection.OwnerNFT{
			ContractId: contractID,
			TokenId:    token.TokenId,
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
		return nil, err
	}

	if err := collection.ValidateTokenID(req.TokenId); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	legacyToken, err := s.getToken(ctx, req.ContractId, req.TokenId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	any, err := codectypes.NewAnyWithValue(legacyToken)
	if err != nil {
		panic(err)
	}

	return &collection.QueryTokenResponse{Token: *any}, nil
}

func (s queryServer) Root(c context.Context, req *collection.QueryRootRequest) (*collection.QueryRootResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	if err := collection.ValidateNFTID(req.TokenId); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := s.validateExistenceOfCollectionGRPC(ctx, req.ContractId); err != nil {
		return nil, err
	}

	if err := s.keeper.hasNFT(ctx, req.ContractId, req.TokenId); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
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
		return nil, err
	}

	if err := collection.ValidateNFTID(req.TokenId); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := s.validateExistenceOfCollectionGRPC(ctx, req.ContractId); err != nil {
		return nil, err
	}

	if err := s.keeper.hasNFT(ctx, req.ContractId, req.TokenId); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	parent, err := s.keeper.GetParent(ctx, req.ContractId, req.TokenId)
	if err != nil {
		return nil, nil
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
		return nil, err
	}

	if err := collection.ValidateNFTID(req.TokenId); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := s.validateExistenceOfCollectionGRPC(ctx, req.ContractId); err != nil {
		return nil, err
	}

	if err := s.keeper.hasNFT(ctx, req.ContractId, req.TokenId); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
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
		return nil, err
	}

	return &collection.QueryChildrenResponse{Children: children, Pagination: pageRes}, nil
}

func (s queryServer) GranteeGrants(c context.Context, req *collection.QueryGranteeGrantsRequest) (*collection.QueryGranteeGrantsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	granteeAddr, err := sdk.AccAddressFromBech32(req.Grantee)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid grantee address: %s", req.Grantee)
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := s.validateExistenceOfCollectionGRPC(ctx, req.ContractId); err != nil {
		return nil, err
	}

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
		return nil, err
	}

	return &collection.QueryGranteeGrantsResponse{Grants: grants, Pagination: pageRes}, nil
}

func (s queryServer) IsOperatorFor(c context.Context, req *collection.QueryIsOperatorForRequest) (*collection.QueryIsOperatorForResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	operator, err := sdk.AccAddressFromBech32(req.Operator)
	if err != nil {
		return nil, err
	}
	holder, err := sdk.AccAddressFromBech32(req.Holder)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := s.validateExistenceOfCollectionGRPC(ctx, req.ContractId); err != nil {
		return nil, err
	}

	_, err = s.keeper.GetAuthorization(ctx, req.ContractId, holder, operator)
	authorized := (err == nil)

	return &collection.QueryIsOperatorForResponse{Authorized: authorized}, nil
}

func (s queryServer) HoldersByOperator(c context.Context, req *collection.QueryHoldersByOperatorRequest) (*collection.QueryHoldersByOperatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	operator, err := sdk.AccAddressFromBech32(req.Operator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", req.Operator)
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(s.keeper.storeKey)
	authorizationStore := prefix.NewStore(store, authorizationKeyPrefixByOperator(req.ContractId, operator))
	var holders []string
	pageRes, err := query.Paginate(authorizationStore, req.Pagination, func(key []byte, value []byte) error {
		holder := sdk.AccAddress(key)
		holders = append(holders, holder.String())
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &collection.QueryHoldersByOperatorResponse{Holders: holders, Pagination: pageRes}, nil
}
