package keeper

import (
	"context"

	"cosmossdk.io/core/store"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"

	v1 "github.com/Finschia/finschia-sdk/x/bankplus/migrations/v1"
)

// DeprecateBankPlus performs in-place store migrations for bankplus v1
// migration includes:
//
// - Remove all the state(inactive addresses)
func DeprecateBankPlus(ctx context.Context, storeService store.KVStoreService) error {
	kvStore := storeService.OpenKVStore(ctx)
	adapter := runtime.KVStoreAdapter(kvStore)
	iterator := storetypes.KVStorePrefixIterator(adapter, v1.InactiveAddrsKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		err := kvStore.Delete(iterator.Key())
		if err != nil {
			return err
		}
	}
	return nil
}
