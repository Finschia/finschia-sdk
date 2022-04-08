package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

func (k Keeper) GetParams(ctx sdk.Context) *foundation.Params {
	store := ctx.KVStore(k.storeKey)
	key := paramsKey
	bz := store.Get(key)

	var params foundation.Params
	k.cdc.MustUnmarshal(bz, &params)

	return &params
}

func (k Keeper) SetParams(ctx sdk.Context, params *foundation.Params) {
	bz := k.cdc.MustMarshal(params)

	store := ctx.KVStore(k.storeKey)
	key := paramsKey
	store.Set(key, bz)
}

// aliases
func (k Keeper) GetEnabled(ctx sdk.Context) bool {
	params := k.GetParams(ctx)

	return params.Enabled
}

func (k Keeper) GetFoundationTax(ctx sdk.Context) sdk.Dec {
	params := k.GetParams(ctx)

	return params.FoundationTax
}
