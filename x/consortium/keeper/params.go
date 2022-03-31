package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/consortium"
)

func (k Keeper) GetParams(ctx sdk.Context) *consortium.Params {
	store := ctx.KVStore(k.storeKey)
	key := paramsKey
	bz := store.Get(key)

	var params consortium.Params
	k.cdc.MustUnmarshalBinaryBare(bz, &params)

	return &params
}

func (k Keeper) SetParams(ctx sdk.Context, params *consortium.Params) {
	bz := k.cdc.MustMarshalBinaryBare(params)

	store := ctx.KVStore(k.storeKey)
	key := paramsKey
	store.Set(key, bz)
}

// aliases
func (k Keeper) GetEnabled(ctx sdk.Context) bool {
	params := k.GetParams(ctx)

	return params.Enabled
}
