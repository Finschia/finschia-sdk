package querier

import (
	"fmt"

	"github.com/line/link/x/collection/internal/keeper"
	"github.com/line/link/x/collection/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// creates a querier for token REST endpoints
func NewQuerier(keeper keeper.Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryBalance:
			return queryBalance(ctx, req, keeper)
		case types.QueryPerms:
			return queryAccountPermission(ctx, req, keeper)
		case types.QueryTokens:
			return queryTokens(ctx, req, keeper)
		case types.QueryCollections:
			return queryCollections(ctx, req, keeper)
		case types.QueryNFTCount:
			return queryNFTCount(ctx, req, keeper)
		case types.QuerySupply:
			return querySupply(ctx, req, keeper)
		case types.QueryParent:
			return queryParent(ctx, req, keeper)
		case types.QueryRoot:
			return queryRoot(ctx, req, keeper)
		case types.QueryChildren:
			return queryChildren(ctx, req, keeper)
		case types.QueryIsApproved:
			return queryIsApproved(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown collection query endpoint")
		}
	}
}

func queryBalance(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, sdk.Error) {
	var params types.QuerySymbolTokenIDAccAddressParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	supply, err := keeper.GetBalance(ctx, params.Symbol, params.TokenID, params.Addr)
	if err != nil {
		return nil, err
	}

	bz, err2 := keeper.MarshalJSONIndent(supply)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}

	return bz, nil
}

func queryTokens(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, sdk.Error) {
	var params types.QuerySymbolTokenIDParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	if len(params.TokenID) == 0 {
		var params types.QuerySymbolParams
		if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
			return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
		}
		tokens, err := keeper.GetTokens(ctx, params.Symbol)
		if err != nil {
			return nil, err
		}
		bz, err2 := keeper.MarshalJSONIndent(tokens)
		if err2 != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
		}
		return bz, nil
	}

	token, err := keeper.GetToken(ctx, params.Symbol, params.TokenID)
	if err != nil {
		return nil, err
	}

	bz, err2 := keeper.MarshalJSONIndent(token)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}

	return bz, nil
}

func queryAccountPermission(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, sdk.Error) {
	if len(req.Data) == 0 {
		return nil, sdk.ErrUnknownRequest("data is nil")
	}
	var params types.QueryAccAddressParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	pms := keeper.GetPermissions(ctx, params.Addr)

	bz, err := keeper.MarshalJSONIndent(pms)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryCollections(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, sdk.Error) {
	if len(req.Data) == 0 {
		collections := keeper.GetAllCollections(ctx)
		bz, err := keeper.MarshalJSONIndent(collections)
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}
		return bz, nil
	}
	var params types.QuerySymbolParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	collection, err := keeper.GetCollection(ctx, params.Symbol)
	if err != nil {
		return nil, err
	}
	bz, err2 := keeper.MarshalJSONIndent(collection)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}

	return bz, nil
}

func queryNFTCount(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, sdk.Error) {
	var params types.QuerySymbolTokenIDParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	count, err := keeper.GetNFTCount(ctx, params.Symbol, params.TokenID)
	if err != nil {
		return nil, err
	}
	bz, err2 := keeper.MarshalJSONIndent(count)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}

	return bz, nil
}

func queryParent(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, sdk.Error) {
	var params types.QuerySymbolTokenIDParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	token, err := keeper.ParentOf(ctx, params.Symbol, params.TokenID)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, nil
	}

	bz, err2 := keeper.MarshalJSONIndent(token)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}

	return bz, nil
}

func queryRoot(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, sdk.Error) {
	var params types.QuerySymbolTokenIDParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	token, err := keeper.RootOf(ctx, params.Symbol, params.TokenID)
	if err != nil {
		return nil, err
	}

	bz, err2 := keeper.MarshalJSONIndent(token)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}

	return bz, nil
}

func queryChildren(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, sdk.Error) {
	var params types.QuerySymbolTokenIDParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	tokens, err := keeper.ChildrenOf(ctx, params.Symbol, params.TokenID)
	if err != nil {
		return nil, err
	}
	if tokens == nil {
		return nil, nil
	}

	bz, err2 := keeper.MarshalJSONIndent(tokens)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}

	return bz, nil
}

func queryIsApproved(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, sdk.Error) {
	var params types.QueryIsApprovedParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params[Proxy=%s, Approver=%s, Symbol=%s]: %s", params.Proxy.String(), params.Approver.String(), params.Symbol, err))
	}

	approved := keeper.IsApproved(ctx, params.Symbol, params.Proxy, params.Approver)

	bz, err := keeper.MarshalJSONIndent(approved)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}
func querySupply(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, sdk.Error) {
	var params types.QuerySymbolTokenIDParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	supply, err := keeper.GetSupplyInt(ctx, params.Symbol, params.TokenID)
	if err != nil {
		return nil, err
	}

	bz, err2 := keeper.MarshalJSONIndent(supply)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}

	return bz, nil
}
