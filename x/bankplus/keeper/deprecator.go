package keeper

import (
	"context"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// inactiveAddrsKeyPrefix Keys for bankplus store but this prefix must not be overlap with bank key prefix.
var inactiveAddrsKeyPrefix = []byte{0xa0}

// inactiveAddrKey key of a specific inactiveAddr from store
func inactiveAddrKey(addr sdk.AccAddress) []byte {
	return append(inactiveAddrsKeyPrefix, addr.Bytes()...)
}

// DeprecateBankPlus performs in-place store migrations for bankplus v1
// migration includes:
//
// - Remove all the state(inactive addresses)
func DeprecateBankPlus(ctx context.Context, keeper BaseKeeper) error {
	kvStore := keeper.storeService.OpenKVStore(ctx)
	adapter := runtime.KVStoreAdapter(kvStore)
	iterator := storetypes.KVStorePrefixIterator(adapter, inactiveAddrsKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		err := kvStore.Delete(iterator.Key())
		if err != nil {
			return err
		}
	}
	return nil
}
