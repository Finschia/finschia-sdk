package token

import (
	"fmt"

	"github.com/link-chain/link/x/token/types"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// creates a querier for auth REST endpoints
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryToken:
			return queryToken(ctx, req, keeper)
		case types.QueryAllTokens:
			return queryAllTokens(ctx, req, keeper)
		case types.QueryPerm:
			return queryAccountPermission(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown token query endpoint")
		}
	}
}

func queryToken(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryTokenParams
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
	var params QueryAccountPermissionParams
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
