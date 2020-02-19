package querier

import (
	"fmt"

	"github.com/line/link/x/token/internal/keeper"
	"github.com/line/link/x/token/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// creates a querier for token REST endpoints
func NewQuerier(keeper keeper.Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryPerms:
			return queryAccountPermission(ctx, req, keeper)
		case types.QueryTokens:
			return queryTokens(ctx, req, keeper)
		case types.QuerySupply:
			return querySupply(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown token query endpoint")
		}
	}
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

func queryTokens(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, sdk.Error) {
	if len(req.Data) == 0 {
		tokens := keeper.GetAllTokens(ctx)

		bz, err := keeper.MarshalJSONIndent(tokens)
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}
		return bz, nil
	}
	var params types.QuerySymbolParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	token, err := keeper.GetToken(ctx, params.Symbol)
	if err != nil {
		return nil, err
	}

	bz, err2 := keeper.MarshalJSONIndent(token)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}

	return bz, nil
}

func querySupply(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, sdk.Error) {
	var params types.QuerySymbolParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	supply, err := keeper.GetSupply(ctx, params.Symbol)
	if err != nil {
		return nil, err
	}

	bz, err2 := keeper.MarshalJSONIndent(supply)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}

	return bz, nil
}
