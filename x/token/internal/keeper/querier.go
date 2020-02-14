package keeper

import (
	"fmt"

	"github.com/line/link/x/token/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// creates a querier for token REST endpoints
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryPerms:
			return queryAccountPermission(ctx, req, keeper)
		case types.QueryTokens:
			return queryTokens(ctx, req, keeper)
		case types.QueryCollections:
			return queryCollections(ctx, req, keeper)
		case types.QuerySupply:
			return querySupply(ctx, req, keeper)
		case types.QueryNFTCount:
			return queryNFTCount(ctx, req, keeper)
		case types.QueryParent:
			return queryParent(ctx, req, keeper)
		case types.QueryRoot:
			return queryRoot(ctx, req, keeper)
		case types.QueryChildren:
			return queryChildren(ctx, req, keeper)
		case types.QueryIsApproved:
			return queryIsApproved(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown token query endpoint")
		}
	}
}

func queryAccountPermission(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	if len(req.Data) == 0 {
		return nil, sdk.ErrUnknownRequest("data is nil")
	}
	var params types.QueryAccAddressParams
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

func queryTokens(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	if len(req.Data) == 0 {
		tokens := keeper.GetAllTokens(ctx)

		bz, err := codec.MarshalJSONIndent(keeper.cdc, tokens)
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}
		return bz, nil
	}
	var params types.QuerySymbolTokenIDParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	token, err := keeper.GetToken(ctx, params.Symbol, params.TokenID)
	if err != nil {
		return nil, err
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, token)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}

	return bz, nil
}

func queryCollections(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	if len(req.Data) == 0 {
		collections := keeper.GetAllCollections(ctx)
		bz, err := codec.MarshalJSONIndent(keeper.cdc, collections)
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}
		return bz, nil
	}
	var params types.QuerySymbolParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	collection, err := keeper.GetCollection(ctx, params.Symbol)
	if err != nil {
		return nil, err
	}
	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, collection)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}

	return bz, nil
}

func querySupply(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QuerySymbolTokenIDParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	supply, err := keeper.GetSupply(ctx, params.Symbol, params.TokenID)
	if err != nil {
		return nil, err
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, supply)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}

	return bz, nil
}

func queryNFTCount(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QuerySymbolTokenIDParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	count, err := keeper.GetNFTCount(ctx, params.Symbol, params.TokenID)
	if err != nil {
		return nil, err
	}
	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, count)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}

	return bz, nil
}

func queryParent(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QuerySymbolTokenIDParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	token, err := keeper.ParentOf(ctx, params.Symbol, params.TokenID)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, nil
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, token)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}

	return bz, nil
}

func queryRoot(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QuerySymbolTokenIDParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	token, err := keeper.RootOf(ctx, params.Symbol, params.TokenID)
	if err != nil {
		return nil, err
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, token)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}

	return bz, nil
}

func queryChildren(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QuerySymbolTokenIDParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	tokens, err := keeper.ChildrenOf(ctx, params.Symbol, params.TokenID)
	if err != nil {
		return nil, err
	}
	if tokens == nil {
		return nil, nil
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, tokens)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}

	return bz, nil
}

func queryIsApproved(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryIsApprovedParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params[Proxy=%s, Approver=%s, Symbol=%s]: %s", params.Proxy.String(), params.Approver.String(), params.Symbol, err))
	}

	approved := keeper.IsApproved(ctx, params.Proxy, params.Approver, params.Symbol)

	bz, err := codec.MarshalJSONIndent(keeper.cdc, approved)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}
