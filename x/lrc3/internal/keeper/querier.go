package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/link-chain/link/x/lrc3/internal/types"
)

// query endpoints supported by the NFT Querier
const (
	QueryLRC3             = "lrc3"
	QueryAllLRC3          = "allLRC3"
	QueryOwnerOf          = "ownerOf"
	QueryBalanceOf        = "balanceOf"
	QueryGetApproved      = "getApproved"
	QueryIsApprovedForAll = "isApprovedForAll"
)

// NewQuerier is the module level router for state queries
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryLRC3:
			return queryLRC3(ctx, path[1:], req, k)
		case QueryAllLRC3:
			return queryAllLRC3(ctx, path[1:], req, k)
		case QueryGetApproved:
			return queryGetApproved(ctx, path[1:], req, k)
		case QueryIsApprovedForAll:
			return queryIsApprovedForAll(ctx, path[1:], req, k)
		case QueryBalanceOf:
			return queryBalanceOf(ctx, path[1:], req, k)
		case QueryOwnerOf:
			return queryOwnerOf(ctx, path[1:], req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nft query endpoint")
		}
	}
}

func queryLRC3(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryLRC3Params
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)

	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}
	lrc3, found := k.NFTKeeper.GetCollection(ctx, params.Denom)
	if !found {
		return nil, types.ErrNotExistLRC3(types.DefaultCodespace)
	}

	bz, err := types.ModuleCdc.MarshalJSON(lrc3)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return bz, nil
}

func queryAllLRC3(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {

	allLRC3 := k.NFTKeeper.GetCollections(ctx)

	bz, err := types.ModuleCdc.MarshalJSON(allLRC3)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return bz, nil
}

func queryGetApproved(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryApproveParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	tokenApprove, _ := k.GetApproval(ctx, params.Denom, params.TokenID)
	bz, err := types.ModuleCdc.MarshalJSON(tokenApprove)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return bz, nil
}

func queryIsApprovedForAll(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryOperatorApproveParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	operators, _ := k.GetOperatorApprovals(ctx, params.Denom, params.OwnerAddress)
	operatorApprovals := types.NewOperatorApprovals(params.Denom, params.OwnerAddress, operators)

	bz, err := types.ModuleCdc.MarshalJSON(operatorApprovals)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return bz, nil
}

func queryBalanceOf(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryBalanceParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	idCollection, _ := k.NFTKeeper.GetOwnerByDenom(ctx, params.Owner, params.Denom)
	balance := types.NewTokenBalance(params.Owner, params.Denom, len(idCollection.IDs))

	bz, err := types.ModuleCdc.MarshalJSON(balance)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return bz, nil
}

func queryOwnerOf(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryOwnerOfParams

	err1 := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err1 != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err1.Error()))
	}

	nft, err2 := k.NFTKeeper.GetNFT(ctx, params.Denom, params.TokenID)
	if err2 != nil {
		return nil, err2
	}

	tokenOwner := types.NewTokenOwner(params.Denom, params.TokenID, nft.GetOwner())

	bz, err := types.ModuleCdc.MarshalJSON(tokenOwner)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return bz, nil
}
