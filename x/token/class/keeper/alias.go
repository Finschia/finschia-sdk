package keeper

import (
	sdk "github.com/Finschia/finschia-sdk/types"
)

// iterate through the classes and perform the provided function
func (k Keeper) iterateIDs(ctx sdk.Context, fn func(id string) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, idKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		id := splitIDKey(iterator.Key())

		stop := fn(id)
		if stop {
			break
		}
	}
}
