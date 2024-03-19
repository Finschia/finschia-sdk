package internal

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

// DeprecateBankPlus performs remove logic for bankplus v1.
// This will remove all the state(inactive addresses)
func DeprecateBankPlus(ctx context.Context, bankKey *storetypes.KVStoreKey) error {
	kss := runtime.NewKVStoreService(bankKey)
	ks := kss.OpenKVStore(ctx)
	adapter := runtime.KVStoreAdapter(ks)
	iter := storetypes.KVStorePrefixIterator(adapter, inactiveAddrsKeyPrefix)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		err := ks.Delete(iter.Key())
		if err != nil {
			return err
		}
	}
	return nil
}
