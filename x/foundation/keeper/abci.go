package keeper

import (
	abci "github.com/line/ostracon/abci/types"

	sdk "github.com/line/lbm-sdk/types"
)

// BeginBlocker withdraws rewards from fee-collector before the distribution
// module's withdraw.
func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k Keeper) {
	if err := k.collectFoundationTax(ctx); err != nil {
		panic(err)
	}
}
