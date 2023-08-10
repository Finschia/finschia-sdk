package keeper

import (
	"github.com/Finschia/finschia-sdk/store/prefix"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/or/rollup/types"
)

func (k Keeper) GetRollup(ctx sdk.Context, rollupName string) (val types.Rollup, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RollupKeyPrefix)
	b := store.Get(types.RollupKey(rollupName))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) SetRollup(ctx sdk.Context, rollup types.Rollup) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RollupKeyPrefix)
	b := k.cdc.MustMarshal(&rollup)
	store.Set(types.RollupKey(rollup.RollupName), b)
}

func (k Keeper) GetAllRollup(ctx sdk.Context) (list []types.Rollup) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RollupKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Rollup
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
