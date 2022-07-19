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
	coin := collection.NewCoin(req.TokenId, balance)

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

// Supply queries the number of tokens from the given contract id.
func (s queryServer) Supply(c context.Context, req *collection.QuerySupplyRequest) (*collection.QuerySupplyResponse, error) {
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
	supply := s.keeper.GetSupply(ctx, req.ContractId, req.ClassId)

	return &collection.QuerySupplyResponse{Supply: supply}, nil
}

// Minted queries the number of tokens from the given contract id.
func (s queryServer) Minted(c context.Context, req *collection.QueryMintedRequest) (*collection.QueryMintedResponse, error) {
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
	minted := s.keeper.GetMinted(ctx, req.ContractId, req.ClassId)

	return &collection.QueryMintedResponse{Minted: minted}, nil
}

// Burnt queries the number of tokens from the given contract id.
func (s queryServer) Burnt(c context.Context, req *collection.QueryBurntRequest) (*collection.QueryBurntResponse, error) {
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
	burnt := s.keeper.GetBurnt(ctx, req.ContractId, req.ClassId)

	return &collection.QueryBurntResponse{Burnt: burnt}, nil
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
	if _, err := s.keeper.GetTokenClass(ctx, req.ContractId, classID); err != nil {
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
	if _, err := s.keeper.GetTokenClass(ctx, req.ContractId, classID); err != nil {
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
	if _, err := s.keeper.GetTokenClass(ctx, req.ContractId, classID); err != nil {
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
	if _, err := s.keeper.GetTokenClass(ctx, req.ContractId, classID); err != nil {
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
	if _, err := s.keeper.GetTokenClass(ctx, req.ContractId, classID); err != nil {
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
	if _, err := s.keeper.GetTokenClass(ctx, req.ContractId, classID); err != nil {
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
		return nil, err
	}

	return &collection.QueryContractResponse{Contract: *contract}, nil
}

func (s queryServer) Contracts(c context.Context, req *collection.QueryContractsRequest) (*collection.QueryContractsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(s.keeper.storeKey)
	contractStore := prefix.NewStore(store, contractKeyPrefix)
	var contracts []collection.Contract
	pageRes, err := query.Paginate(contractStore, req.Pagination, func(key []byte, value []byte) error {
		var contract collection.Contract
		s.keeper.cdc.MustUnmarshal(value, &contract)

		contracts = append(contracts, contract)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &collection.QueryContractsResponse{Contracts: contracts, Pagination: pageRes}, nil
}

// FTClass queries a fungible token class based on its class id.
func (s queryServer) FTClass(c context.Context, req *collection.QueryFTClassRequest) (*collection.QueryFTClassResponse, error) {
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
	class, err := s.keeper.GetTokenClass(ctx, req.ContractId, req.ClassId)
	if err != nil {
		return nil, err
	}
	ftClass, ok := class.(*collection.FTClass)
	if !ok {
		return nil, sdkerrors.ErrInvalidType.Wrapf("not a class of fungible token: %s", req.ClassId)
	}

	return &collection.QueryFTClassResponse{Class: *ftClass}, nil
}

// NFTClass queries a non-fungible token class based on its class id.
func (s queryServer) NFTClass(c context.Context, req *collection.QueryNFTClassRequest) (*collection.QueryNFTClassResponse, error) {
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
	class, err := s.keeper.GetTokenClass(ctx, req.ContractId, req.ClassId)
	if err != nil {
		return nil, err
	}
	nftClass, ok := class.(*collection.NFTClass)
	if !ok {
		return nil, sdkerrors.ErrInvalidType.Wrapf("not a class of non-fungible token: %s", req.ClassId)
	}

	return &collection.QueryNFTClassResponse{Class: *nftClass}, nil
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
	class, err := s.keeper.GetTokenClass(ctx, req.ContractId, req.ClassId)
	if err != nil {
		return nil, err
	}
	name := proto.MessageName(class)

	return &collection.QueryTokenClassTypeNameResponse{Name: name}, nil
}

// TokenClasses queries all token class metadata.
// func (s queryServer) TokenClasses(c context.Context, req *collection.QueryTokenClassesRequest) (*collection.QueryTokenClassesResponse, error) {
// 	if req == nil {
// 		return nil, status.Error(codes.InvalidArgument, "empty request")
// 	}

// 	if err := collection.ValidateContractID(req.ContractId); err != nil {
// 		return nil, err
// 	}

// 	ctx := sdk.UnwrapSDKContext(c)
// 	store := ctx.KVStore(s.keeper.storeKey)
// 	classStore := prefix.NewStore(store, classKeyPrefix)
// 	var classes []codectypes.Any
// 	pageRes, err := query.Paginate(classStore, req.Pagination, func(key []byte, value []byte) error {
// 		var class collection.TokenClass
// 		if err := s.keeper.cdc.UnmarshalInterface(value, &class); err != nil {
// 			panic(err)
// 		}
// 		classes = append(classes, *collection.TokenClassToAny(class))
// 		return nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &collection.QueryTokenClassesResponse{Classes: classes, Pagination: pageRes}, nil
// }

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
		return nil, err
	}

	nftClass, ok := class.(*collection.NFTClass)
	if !ok {
		return nil, sdkerrors.ErrInvalidType.Wrapf("not a class of non-fungible token: %s", classID)
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
		return nil, err
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
		return nil, err
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
			TokenId:    ftClass.Id,
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
		return nil, err
	}

	any, err := codectypes.NewAnyWithValue(legacyToken)
	if err != nil {
		panic(err)
	}

	return &collection.QueryTokenResponse{Token: *any}, nil
}

func (s queryServer) Tokens(c context.Context, req *collection.QueryTokensRequest) (*collection.QueryTokensResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
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
		return nil, err
	}

	return &collection.QueryTokensResponse{Tokens: tokens, Pagination: pageRes}, nil
}

func (s queryServer) NFT(c context.Context, req *collection.QueryNFTRequest) (*collection.QueryNFTResponse, error) {
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
	token, err := s.keeper.GetNFT(ctx, req.ContractId, req.TokenId)
	if err != nil {
		return nil, err
	}

	return &collection.QueryNFTResponse{Token: *token}, nil
}

func (s queryServer) Owner(c context.Context, req *collection.QueryOwnerRequest) (*collection.QueryOwnerResponse, error) {
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
	if err := s.keeper.hasNFT(ctx, req.ContractId, req.TokenId); err != nil {
		return nil, err
	}

	owner := s.keeper.GetRootOwner(ctx, req.ContractId, req.TokenId)

	return &collection.QueryOwnerResponse{Owner: owner.String()}, nil
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
	if err := s.keeper.hasNFT(ctx, req.ContractId, req.TokenId); err != nil {
		return nil, err
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
	if err := s.keeper.hasNFT(ctx, req.ContractId, req.TokenId); err != nil {
		return nil, err
	}

	parent, err := s.keeper.GetParent(ctx, req.ContractId, req.TokenId)
	if err != nil {
		return nil, err
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
	if err := s.keeper.hasNFT(ctx, req.ContractId, req.TokenId); err != nil {
		return nil, err
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

func (s queryServer) Grant(c context.Context, req *collection.QueryGrantRequest) (*collection.QueryGrantResponse, error) {
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

	if err := collection.ValidatePermission(req.Permission); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	grant, err := s.keeper.GetGrant(ctx, req.ContractId, granteeAddr, req.Permission)
	if err != nil {
		return nil, err
	}

	return &collection.QueryGrantResponse{Grant: *grant}, nil
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

func (s queryServer) Authorization(c context.Context, req *collection.QueryAuthorizationRequest) (*collection.QueryAuthorizationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	operatorAddr, err := sdk.AccAddressFromBech32(req.Operator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", req.Operator)
	}
	holderAddr, err := sdk.AccAddressFromBech32(req.Holder)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid holder address: %s", req.Holder)
	}

	ctx := sdk.UnwrapSDKContext(c)
	authorization, err := s.keeper.GetAuthorization(ctx, req.ContractId, holderAddr, operatorAddr)
	if err != nil {
		return nil, err
	}

	return &collection.QueryAuthorizationResponse{Authorization: *authorization}, nil
}

func (s queryServer) OperatorAuthorizations(c context.Context, req *collection.QueryOperatorAuthorizationsRequest) (*collection.QueryOperatorAuthorizationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	operatorAddr, err := sdk.AccAddressFromBech32(req.Operator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", req.Operator)
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(s.keeper.storeKey)
	authorizationStore := prefix.NewStore(store, authorizationKeyPrefixByOperator(req.ContractId, operatorAddr))
	var authorizations []collection.Authorization
	pageRes, err := query.Paginate(authorizationStore, req.Pagination, func(key []byte, value []byte) error {
		holder := sdk.AccAddress(key)
		authorizations = append(authorizations, collection.Authorization{
			Holder:   holder.String(),
			Operator: req.Operator,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &collection.QueryOperatorAuthorizationsResponse{Authorizations: authorizations, Pagination: pageRes}, nil
}

func (s queryServer) Approved(c context.Context, req *collection.QueryApprovedRequest) (*collection.QueryApprovedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}
	approverAddr, err := sdk.AccAddressFromBech32(req.Approver)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid address address: %s", req.Address)
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(s.keeper.storeKey)
	authorizationStore := prefix.NewStore(store, authorizationKeyPrefixByOperator(req.ContractId, addr))
	var approvers []string
	pageRes, err := query.Paginate(authorizationStore, req.Pagination, func(key []byte, value []byte) error {
		holder := string(key)
		approvers = append(approvers, holder)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &collection.QueryApproversResponse{Approvers: approvers, Pagination: pageRes}, nil
}
