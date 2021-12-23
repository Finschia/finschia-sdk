package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/consortium/types"
)

func (k Keeper) GetParams(ctx sdk.Context) *types.Params {
	store := ctx.KVStore(k.storeKey)
	key := types.ParamsKey
	bz := store.Get(key)

	var params types.Params
	k.cdc.MustUnmarshalBinaryBare(bz, &params)

	return &params
}

func (k Keeper) SetParams(ctx sdk.Context, params *types.Params) {
	bz := k.cdc.MustMarshalBinaryBare(params)

	store := ctx.KVStore(k.storeKey)
	key := types.ParamsKey
	store.Set(key, bz)
}

// aliases
func (k Keeper) GetEnabled(ctx sdk.Context) bool {
	params := k.GetParams(ctx)

	return params.Enabled
}
