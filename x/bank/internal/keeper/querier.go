package keeper

import (
	"fmt"

	"github.com/line/link/x/bank/internal/types"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// query balance path
	QueryBalance      = "balances"
	QueryBulkBalances = "bulk_balances"
)

// NewQuerier returns a new sdk.Keeper instance.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryBalance:
			return queryBalance(ctx, req, k)
		case QueryBulkBalances:
			return queryBulkBalances(ctx, req, k)

		default:
			return nil, sdk.ErrUnknownRequest("unknown bank query endpoint")
		}
	}
}

// queryBalance fetch an account's balance for the supplied height.
// Height and account address are passed as first and second path components respectively.
func queryBalance(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryBalanceParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	if len(params.Denom) == 0 {
		bz, err := codec.MarshalJSONIndent(types.ModuleCdc, k.GetCoins(ctx, params.Address))
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}

		return bz, nil
	}
	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, k.GetCoins(ctx, params.Address).AmountOf(params.Denom))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryBulkBalances(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryBulkBalancesParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	if len(params.Addresses) > types.RequestGetsLimit {
		return nil, types.ErrRequestGetsLimit(types.DefaultCodespace, types.RequestGetsLimit).TraceSDK("")
	}

	res := make([]types.QueryBulkBalancesResult, len(params.Addresses))
	for idx, addr := range params.Addresses {
		res[idx] = types.NewQueryBulkBalancesResult(addr, k.GetCoins(ctx, addr))
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, res)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}
