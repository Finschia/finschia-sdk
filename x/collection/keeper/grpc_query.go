package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/gogoproto/proto"

	"github.com/Finschia/finschia-sdk/x/collection"
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

func (s queryServer) addressFromBech32GRPC(bech32, context string) (sdk.AccAddress, error) {
	addr, err := sdk.AccAddressFromBech32(bech32)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, errorsmod.Wrap(sdkerrors.ErrInvalidAddress.Wrap(bech32), context).Error())
	}

	return addr, nil
}

func (s queryServer) assertTokenIsFungible(ctx sdk.Context, contractID, classID string) error {
	class, err := s.keeper.GetTokenClass(ctx, contractID, classID)
	if err != nil {
		return err
	}

	if _, ok := class.(*collection.FTClass); !ok {
		return collection.ErrTokenNotExist.Wrap(collection.NewFTID(classID))
	}

	return nil
}

func (s queryServer) assertTokenTypeIsNonFungible(ctx sdk.Context, contractID, classID string) error {
	class, err := s.keeper.GetTokenClass(ctx, contractID, classID)
	if err != nil {
		return err
	}

	if _, ok := class.(*collection.NFTClass); !ok {
		return collection.ErrTokenTypeNotExist.Wrap(classID)
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
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	addr, err := s.addressFromBech32GRPC(req.Address, "address")
	if err != nil {
		return nil, err
	}

	if err := collection.ValidateTokenID(req.TokenId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
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
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	addr, err := s.addressFromBech32GRPC(req.Address, "address")
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := s.keeper.storeService.OpenKVStore(ctx)
	adapter := runtime.KVStoreAdapter(store)
	balanceStore := prefix.NewStore(adapter, balanceKeyPrefixByAddress(req.ContractId, addr))
	var balances []collection.Coin
	pageRes, err := query.Paginate(balanceStore, req.Pagination, func(key, value []byte) error {
		tokenID := string(key)

		var balance math.Int
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
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := collection.ValidateTokenID(req.TokenId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	classID := collection.SplitTokenID(req.TokenId)

	ctx := sdk.UnwrapSDKContext(c)

	if err := s.assertTokenIsFungible(ctx, req.ContractId, classID); err != nil {
		return &collection.QueryFTSupplyResponse{Supply: math.ZeroInt()}, nil
	}

	supply := s.keeper.GetSupply(ctx, req.ContractId, classID)

	return &collection.QueryFTSupplyResponse{Supply: supply}, nil
}

func (s queryServer) FTMinted(c context.Context, req *collection.QueryFTMintedRequest) (*collection.QueryFTMintedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := collection.ValidateTokenID(req.TokenId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	classID := collection.SplitTokenID(req.TokenId)

	ctx := sdk.UnwrapSDKContext(c)

	if err := s.assertTokenIsFungible(ctx, req.ContractId, classID); err != nil {
		return &collection.QueryFTMintedResponse{Minted: math.ZeroInt()}, nil
	}

	minted := s.keeper.GetMinted(ctx, req.ContractId, classID)

	return &collection.QueryFTMintedResponse{Minted: minted}, nil
}

func (s queryServer) FTBurnt(c context.Context, req *collection.QueryFTBurntRequest) (*collection.QueryFTBurntResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := collection.ValidateTokenID(req.TokenId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	classID := collection.SplitTokenID(req.TokenId)

	ctx := sdk.UnwrapSDKContext(c)

	if err := s.assertTokenIsFungible(ctx, req.ContractId, classID); err != nil {
		return &collection.QueryFTBurntResponse{Burnt: math.ZeroInt()}, nil
	}

	burnt := s.keeper.GetBurnt(ctx, req.ContractId, classID)

	return &collection.QueryFTBurntResponse{Burnt: burnt}, nil
}

func (s queryServer) NFTSupply(c context.Context, req *collection.QueryNFTSupplyRequest) (*collection.QueryNFTSupplyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	classID := req.TokenType
	if err := collection.ValidateClassID(classID); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := s.assertTokenTypeIsNonFungible(ctx, req.ContractId, classID); err != nil {
		return &collection.QueryNFTSupplyResponse{Supply: math.ZeroInt()}, nil
	}

	supply := s.keeper.GetSupply(ctx, req.ContractId, classID)

	return &collection.QueryNFTSupplyResponse{Supply: supply}, nil
}

func (s queryServer) NFTMinted(c context.Context, req *collection.QueryNFTMintedRequest) (*collection.QueryNFTMintedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	classID := req.TokenType
	if err := collection.ValidateClassID(classID); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := s.assertTokenTypeIsNonFungible(ctx, req.ContractId, classID); err != nil {
		return &collection.QueryNFTMintedResponse{Minted: math.ZeroInt()}, nil
	}

	minted := s.keeper.GetMinted(ctx, req.ContractId, classID)

	return &collection.QueryNFTMintedResponse{Minted: minted}, nil
}

func (s queryServer) NFTBurnt(c context.Context, req *collection.QueryNFTBurntRequest) (*collection.QueryNFTBurntResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	classID := req.TokenType
	if err := collection.ValidateClassID(classID); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := s.assertTokenTypeIsNonFungible(ctx, req.ContractId, classID); err != nil {
		return &collection.QueryNFTBurntResponse{Burnt: math.ZeroInt()}, nil
	}

	burnt := s.keeper.GetBurnt(ctx, req.ContractId, classID)

	return &collection.QueryNFTBurntResponse{Burnt: burnt}, nil
}

func (s queryServer) Contract(c context.Context, req *collection.QueryContractRequest) (*collection.QueryContractResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
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
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := collection.ValidateClassID(req.ClassId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)
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
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	classID := req.TokenType
	if err := collection.ValidateClassID(classID); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
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

func (s queryServer) getToken(ctx sdk.Context, contractID, tokenID string) (collection.Token, error) {
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
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := collection.ValidateTokenID(req.TokenId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
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
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := collection.ValidateNFTID(req.TokenId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)
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

func (s queryServer) HasParent(c context.Context, req *collection.QueryHasParentRequest) (*collection.QueryHasParentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.GetContractId()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := collection.ValidateNFTID(req.GetTokenId()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)
	_, err := s.keeper.GetParent(ctx, req.ContractId, req.TokenId)
	return &collection.QueryHasParentResponse{HasParent: (err == nil)}, nil
}

func (s queryServer) Parent(c context.Context, req *collection.QueryParentRequest) (*collection.QueryParentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.GetContractId()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := collection.ValidateNFTID(req.GetTokenId()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)
	parent, err := s.keeper.GetParent(ctx, req.ContractId, req.TokenId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
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
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := collection.ValidateNFTID(req.TokenId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)

	store := s.keeper.storeService.OpenKVStore(ctx)
	adapter := runtime.KVStoreAdapter(store)
	childStore := prefix.NewStore(adapter, childKeyPrefixByTokenID(req.ContractId, req.TokenId))
	var children []collection.NFT
	pageRes, err := query.Paginate(childStore, req.Pagination, func(key, _ []byte) error {
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
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	granteeAddr, err := s.addressFromBech32GRPC(req.Grantee, "grantee")
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := s.keeper.storeService.OpenKVStore(ctx)
	adapter := runtime.KVStoreAdapter(store)
	grantStore := prefix.NewStore(adapter, grantKeyPrefixByGrantee(req.ContractId, granteeAddr))
	var grants []collection.Grant
	pageRes, err := query.Paginate(grantStore, req.Pagination, func(key, _ []byte) error {
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
	authorized := (err == nil)

	return &collection.QueryIsOperatorForResponse{Authorized: authorized}, nil
}

func (s queryServer) HoldersByOperator(c context.Context, req *collection.QueryHoldersByOperatorRequest) (*collection.QueryHoldersByOperatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	operator, err := s.addressFromBech32GRPC(req.Operator, "operator")
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := s.keeper.storeService.OpenKVStore(ctx)
	adapter := runtime.KVStoreAdapter(store)
	authorizationStore := prefix.NewStore(adapter, authorizationKeyPrefixByOperator(req.ContractId, operator))
	var holders []string
	pageRes, err := query.Paginate(authorizationStore, req.Pagination, func(key, value []byte) error {
		holder := sdk.AccAddress(key)
		holders = append(holders, holder.String())
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &collection.QueryHoldersByOperatorResponse{Holders: holders, Pagination: pageRes}, nil
}
