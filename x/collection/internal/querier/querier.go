package querier

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link-modules/x/collection/internal/keeper"
	"github.com/line/link-modules/x/collection/internal/types"
	"github.com/line/link-modules/x/contract"
	abci "github.com/tendermint/tendermint/abci/types"
)

// creates a querier for token REST endpoints
func NewQuerier(keeper keeper.Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		if len(path) >= 2 {
			ctx = ctx.WithContext(context.WithValue(ctx.Context(), contract.CtxKey{}, path[1]))
		}
		switch path[0] {
		case types.QueryBalance:
			return queryBalance(ctx, req, keeper)
		case types.QueryBalances:
			return queryBalances(ctx, req, keeper)
		case types.QueryPerms:
			return queryAccountPermission(ctx, req, keeper)
		case types.QueryTokens:
			return queryTokens(ctx, req, keeper)
		case types.QueryTokensWithTokenType:
			return queryTokensWithTokenType(ctx, req, keeper)
		case types.QueryTokenTypes:
			return queryTokenTypes(ctx, req, keeper)
		case types.QueryCollections:
			return queryCollections(ctx, req, keeper)
		case types.QueryNFTCount:
			return queryNFTCount(ctx, req, keeper, types.QueryNFTCount)
		case types.QueryNFTMint:
			return queryNFTCount(ctx, req, keeper, types.QueryNFTMint)
		case types.QueryNFTBurn:
			return queryNFTCount(ctx, req, keeper, types.QueryNFTBurn)
		case types.QuerySupply:
			return queryTotal(ctx, req, keeper, types.QuerySupply)
		case types.QueryMint:
			return queryTotal(ctx, req, keeper, types.QueryMint)
		case types.QueryBurn:
			return queryTotal(ctx, req, keeper, types.QueryBurn)
		case types.QueryParent:
			return queryParent(ctx, req, keeper)
		case types.QueryRoot:
			return queryRoot(ctx, req, keeper)
		case types.QueryChildren:
			return queryChildren(ctx, req, keeper)
		case types.QueryApprovers:
			return queryApprovers(ctx, req, keeper)
		case types.QueryIsApproved:
			return queryIsApproved(ctx, req, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown collection query endpoint")
		}
	}
}

func queryBalance(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, error) {
	var params types.QueryTokenIDAccAddressParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	if !keeper.HasContractID(ctx) {
		return nil, sdkerrors.Wrap(types.ErrCollectionNotExist, ctx.Context().Value(contract.CtxKey{}).(string))
	}

	if !keeper.HasToken(ctx, params.TokenID) {
		return nil, sdkerrors.Wrapf(types.ErrTokenNotExist, "%s %s", ctx.Context().Value(contract.CtxKey{}).(string), params.TokenID)
	}

	var balance sdk.Int
	if params.TokenID[0] == types.FungibleFlag[0] {
		var err error
		balance, err = keeper.GetBalance(ctx, params.TokenID, params.Addr)
		if err != nil {
			if _, err2 := keeper.GetAccount(ctx, params.Addr); err2 != nil {
				balance = sdk.ZeroInt()
			} else {
				return nil, err
			}
		}
	} else {
		if keeper.HasNFTOwner(ctx, params.Addr, params.TokenID) {
			balance = sdk.NewInt(1)
		} else {
			balance = sdk.NewInt(0)
		}
	}

	bz, err2 := keeper.MarshalJSONIndent(balance)
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err2.Error())
	}

	return bz, nil
}

func queryBalances(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, error) {
	var params types.QueryTokenIDAccAddressParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	coins := make([]types.Coin, 0)
	if !keeper.HasContractID(ctx) {
		return nil, sdkerrors.Wrap(types.ErrCollectionNotExist, ctx.Context().Value(contract.CtxKey{}).(string))
	}
	// FT
	acc, err := keeper.GetAccount(ctx, params.Addr)
	if err == nil {
		coins = acc.GetCoins()
	}

	// NFT
	tokenIds := keeper.GetNFTsOwner(ctx, params.Addr)
	for _, tokenID := range tokenIds {
		var coin types.Coin
		coin.Amount = sdk.NewInt(1)
		coin.Denom = tokenID
		coins = append(coins, coin)
	}

	bz, err2 := keeper.MarshalJSONIndent(coins)
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err2.Error())
	}

	return bz, nil
}

// nolint:dupl
func queryTokenTypes(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, error) {
	var params types.QueryTokenIDParams
	if len(req.Data) != 0 {
		if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}
	}
	if len(params.TokenID) == 0 {
		tokenTypes, err := keeper.GetTokenTypes(ctx)
		if err != nil {
			return nil, err
		}
		bz, err2 := keeper.MarshalJSONIndent(tokenTypes)
		if err2 != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err2.Error())
		}
		return bz, nil
	}

	tokenType, err := keeper.GetTokenType(ctx, params.TokenID)
	if err != nil {
		return nil, err
	}

	bz, err2 := keeper.MarshalJSONIndent(tokenType)
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err2.Error())
	}

	return bz, nil
}

// nolint:dupl
func queryTokens(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, error) {
	var params types.QueryTokenIDParams
	if len(req.Data) != 0 {
		if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}
	}
	if len(params.TokenID) == 0 {
		tokens, err := keeper.GetTokens(ctx)
		if err != nil {
			return nil, err
		}
		bz, err2 := keeper.MarshalJSONIndent(tokens)
		if err2 != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err2.Error())
		}
		return bz, nil
	}

	token, err := keeper.GetToken(ctx, params.TokenID)
	if err != nil {
		return nil, err
	}

	bz, err2 := keeper.MarshalJSONIndent(token)
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err2.Error())
	}

	return bz, nil
}

func queryTokensWithTokenType(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, error) {
	var params types.QueryTokenTypeParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	tokens, err := keeper.GetNFTs(ctx, params.TokenType)
	if err != nil {
		return nil, err
	}
	bz, err2 := keeper.MarshalJSONIndent(tokens)
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err2.Error())
	}
	return bz, nil
}

func queryAccountPermission(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, error) {
	if len(req.Data) == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "data is nil")
	}
	var params types.QueryAccAddressParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	pms := keeper.GetPermissions(ctx, params.Addr)

	bz, err := keeper.MarshalJSONIndent(pms)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryCollections(ctx sdk.Context, _ abci.RequestQuery, keeper keeper.Keeper) ([]byte, error) {
	collection, err := keeper.GetCollection(ctx)
	if err != nil {
		return nil, err
	}
	bz, err2 := keeper.MarshalJSONIndent(collection)
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err2.Error())
	}

	return bz, nil
}

func queryNFTCount(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper, target string) ([]byte, error) {
	var params types.QueryTokenIDParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	count, err := keeper.GetNFTCountInt(ctx, params.TokenID, target)
	if err != nil {
		return nil, err
	}
	bz, err2 := keeper.MarshalJSONIndent(count)
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err2.Error())
	}

	return bz, nil
}

func queryParent(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, error) {
	var params types.QueryTokenIDParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	token, err := keeper.ParentOf(ctx, params.TokenID)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, nil
	}

	bz, err2 := keeper.MarshalJSONIndent(token)
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err2.Error())
	}

	return bz, nil
}

func queryRoot(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, error) {
	var params types.QueryTokenIDParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	token, err := keeper.RootOf(ctx, params.TokenID)
	if err != nil {
		return nil, err
	}

	bz, err2 := keeper.MarshalJSONIndent(token)
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err2.Error())
	}

	return bz, nil
}

func queryChildren(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, error) {
	var params types.QueryTokenIDParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	tokens, err := keeper.ChildrenOf(ctx, params.TokenID)
	if err != nil {
		return nil, err
	}
	if tokens == nil {
		return nil, nil
	}

	bz, err2 := keeper.MarshalJSONIndent(tokens)
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err2.Error())
	}

	return bz, nil
}

func queryApprovers(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, error) {
	var params types.QueryProxyParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	approvers, err := keeper.GetApprovers(ctx, params.Proxy)
	if err != nil {
		return nil, err
	}

	bz, err2 := keeper.MarshalJSONIndent(approvers)
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err2.Error())
	}

	return bz, nil
}

func queryIsApproved(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, error) {
	var params types.QueryIsApprovedParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	approved := keeper.IsApproved(ctx, params.Proxy, params.Approver)

	bz, err := keeper.MarshalJSONIndent(approved)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryTotal(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper, target string) ([]byte, error) {
	var params types.QueryTokenIDParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	supply, err := keeper.GetTotalInt(ctx, params.TokenID, target)
	if err != nil {
		return nil, err
	}

	bz, err2 := keeper.MarshalJSONIndent(supply)
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err2.Error())
	}

	return bz, nil
}
