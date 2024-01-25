package keeper

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// iterate through the classes and perform the provided function
func (k Keeper) iterateIDs(ctx sdk.Context, fn func(id string) (stop bool)) {
	store := k.storeService.OpenKVStore(ctx)
	adapter := runtime.KVStoreAdapter(store)
	iterator := storetypes.KVStorePrefixIterator(adapter, idKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		id := splitIDKey(iterator.Key())

		stop := fn(id)
		if stop {
			break
		}
	}
}
