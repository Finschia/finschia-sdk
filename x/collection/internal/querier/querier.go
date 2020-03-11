package querier

import (
	"github.com/line/link/x/collection/internal/keeper"
	"github.com/line/link/x/collection/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// creates a querier for token REST endpoints
func NewQuerier(keeper keeper.Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryBalance:
			return queryBalance(ctx, req, keeper)
		case types.QueryPerms:
			return queryAccountPermission(ctx, req, keeper)
		case types.QueryTokens:
			return queryTokens(ctx, req, keeper)
		case types.QueryTokenTypes:
			return queryTokenTypes(ctx, req, keeper)
		case types.QueryCollections:
			return queryCollections(ctx, req, keeper)
		case types.QueryNFTCount:
			return queryNFTCount(ctx, req, keeper)
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
		case types.QueryIsApproved:
			return queryIsApproved(ctx, req, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown collection query endpoint")
		}
	}
}

func queryBalance(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, error) {
	var params types.QueryContractIDTokenIDAccAddressParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	supply, err := keeper.GetBalance(ctx, params.ContractID, params.TokenID, params.Addr)
	if err != nil {
		return nil, err
	}

	bz, err2 := keeper.MarshalJSONIndent(supply)
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err2.Error())
	}

	return bz, nil
}

//nolint:dupl
func queryTokenTypes(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, error) {
	var params types.QueryContractIDTokenIDParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	if len(params.TokenID) == 0 {
		var params types.QueryContractIDParams
		if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}
		tokenTypes, err := keeper.GetTokenTypes(ctx, params.ContractID)
		if err != nil {
			return nil, err
		}
		bz, err2 := keeper.MarshalJSONIndent(tokenTypes)
		if err2 != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err2.Error())
		}
		return bz, nil
	}

	tokenType, err := keeper.GetTokenType(ctx, params.ContractID, params.TokenID)
	if err != nil {
		return nil, err
	}

	bz, err2 := keeper.MarshalJSONIndent(tokenType)
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err2.Error())
	}

	return bz, nil
}

//nolint:dupl
func queryTokens(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, error) {
	var params types.QueryContractIDTokenIDParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	if len(params.TokenID) == 0 {
		var params types.QueryContractIDParams
		if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}
		tokens, err := keeper.GetTokens(ctx, params.ContractID)
		if err != nil {
			return nil, err
		}
		bz, err2 := keeper.MarshalJSONIndent(tokens)
		if err2 != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err2.Error())
		}
		return bz, nil
	}

	token, err := keeper.GetToken(ctx, params.ContractID, params.TokenID)
	if err != nil {
		return nil, err
	}

	bz, err2 := keeper.MarshalJSONIndent(token)
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

func queryCollections(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, error) {
	if len(req.Data) == 0 {
		collections := keeper.GetAllCollections(ctx)
		bz, err := keeper.MarshalJSONIndent(collections)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		return bz, nil
	}
	var params types.QueryContractIDParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	collection, err := keeper.GetCollection(ctx, params.ContractID)
	if err != nil {
		return nil, err
	}
	bz, err2 := keeper.MarshalJSONIndent(collection)
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err2.Error())
	}

	return bz, nil
}

func queryNFTCount(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, error) {
	var params types.QueryContractIDTokenIDParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	count, err := keeper.GetNFTCount(ctx, params.ContractID, params.TokenID)
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
	var params types.QueryContractIDTokenIDParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	token, err := keeper.ParentOf(ctx, params.ContractID, params.TokenID)
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
	var params types.QueryContractIDTokenIDParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	token, err := keeper.RootOf(ctx, params.ContractID, params.TokenID)
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
	var params types.QueryContractIDTokenIDParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	tokens, err := keeper.ChildrenOf(ctx, params.ContractID, params.TokenID)
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

func queryIsApproved(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, error) {
	var params types.QueryIsApprovedParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	approved := keeper.IsApproved(ctx, params.ContractID, params.Proxy, params.Approver)

	bz, err := keeper.MarshalJSONIndent(approved)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryTotal(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper, target string) ([]byte, error) {
	var params types.QueryContractIDTokenIDParams
	if err := keeper.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	supply, err := keeper.GetTotalInt(ctx, params.ContractID, params.TokenID, target)
	if err != nil {
		return nil, err
	}

	bz, err2 := keeper.MarshalJSONIndent(supply)
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err2.Error())
	}

	return bz, nil
}
