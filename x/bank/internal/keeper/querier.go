package keeper

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cbank "github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/link-chain/link/x/bank/internal/types"
)

const (
	QueryBalanceOf = "balance_of"
)

func NewQuerier(k cbank.Keeper, fallbackQuerier sdk.Querier) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryBalanceOf:
			return queryBalanceOf(ctx, req, k)

		default:
			return fallbackQuerier(ctx, path, req)
		}
	}
}

// queryBalanceOf fetch an account's balance of given denomination for the supplied height.
// Height and account address are passed as first and second path components respectively.
func queryBalanceOf(ctx sdk.Context, req abci.RequestQuery, k cbank.Keeper) ([]byte, sdk.Error) {
	var params types.QueryBalanceOfParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, k.GetCoins(ctx, params.Address).AmountOf(params.Denom))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}
