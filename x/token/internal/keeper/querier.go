package keeper

import (
	"fmt"
	"github.com/line/link/x/token/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// creates a querier for auth REST endpoints
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryTokens:
			return queryTokens(ctx, req, keeper)
		case types.QueryPerms:
			return queryAccountPermission(ctx, req, keeper)
		case types.QueryCollections:
			return queryCollections(ctx, req, keeper)
		case types.QuerySupply:
			return querySupply(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown token query endpoint")
		}
	}
}

func queryTokens(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	if len(req.Data) == 0 {
		return queryAllTokens(ctx, req, keeper)
	}
	return queryToken(ctx, req, keeper)
}

func queryToken(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryTokenParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	token, err := keeper.GetToken(ctx, params.String())
	if err != nil {
		return nil, err
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, token)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}

	return bz, nil
}

func queryAllTokens(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	tokens := keeper.GetAllTokens(ctx)

	bz, err := codec.MarshalJSONIndent(keeper.cdc, tokens)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryAccountPermission(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	if len(req.Data) == 0 {
		return nil, sdk.ErrUnknownRequest("data is nil")
	}
	var params types.QueryAccountPermissionParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	pms := keeper.iamKeeper.GetPermissions(ctx, params.Addr)

	bz, err := codec.MarshalJSONIndent(keeper.cdc, pms)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func querySupply(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QuerySupplyParams

	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	supply, err := keeper.GetSupply(ctx, params.String())
	if err != nil {
		return nil, err
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, supply)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}

	return bz, nil
}

func queryCollections(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	if len(req.Data) == 0 {
		return queryAllCollections(ctx, req, keeper)
	}
	return queryCollection(ctx, req, keeper)
}

func queryCollection(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryCollectionParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	var err sdk.Error
	var collection types.Collection
	collection, err = keeper.GetCollection(ctx, params.String())
	if err != nil {
		return nil, err
	}
	var tokens types.Tokens
	tokens = keeper.GetPrefixedTokens(ctx, params.String())
	collectionWithTokens := types.CollectionWithTokens{
		Collection: collection,
		Tokens:     tokens,
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, collectionWithTokens)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}

	return bz, nil
}

func queryAllCollections(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	collections := keeper.GetAllCollections(ctx)

	var collectionsWithTokens types.CollectionsWithTokens
	for _, collection := range collections {
		var tokens types.Tokens
		tokens = keeper.GetPrefixedTokens(ctx, collection.Symbol)
		collectionWithTokens := types.CollectionWithTokens{
			Collection: collection,
			Tokens:     tokens,
		}
		collectionsWithTokens = append(collectionsWithTokens, collectionWithTokens)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, collectionsWithTokens)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
