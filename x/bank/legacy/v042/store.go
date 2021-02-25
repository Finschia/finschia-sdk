package v042

import (
	"github.com/line/lbm-sdk/store/prefix"
	sdk "github.com/line/lbm-sdk/types"
	v040auth "github.com/line/lbm-sdk/x/auth/legacy/v040"
	v040bank "github.com/line/lbm-sdk/x/bank/legacy/v040"
	"github.com/line/lbm-sdk/x/bank/types"
)

// MigrateStore performs in-place store migrations from v0.40 to v0.42. The
// migration includes:
//
// - Change addresses to be length-prefixed.
// - Change balances prefix to 1 byte
func MigrateStore(ctx sdk.Context, storeKey sdk.StoreKey) error {
	store := ctx.KVStore(storeKey)

	// old key is of format:
	// prefix ("balances") || addrBytes (20 bytes) || denomBytes
	// new key is of format
	// prefix (0x02) || addrLen (1 byte) || addrBytes || denomBytes
	oldStore := prefix.NewStore(store, v040bank.BalancesPrefix)

	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()

	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		addr := v040bank.AddressFromBalancesStore(oldStoreIter.Key())
		denom := oldStoreIter.Key()[v040auth.AddrLen:]
		newStoreKey := append(types.CreateAccountBalancesPrefix(addr.Bytes()), denom...)

		// Set new key on store. Values don't change.
		store.Set(newStoreKey, oldStoreIter.Value())
		oldStore.Delete(oldStoreIter.Key())
	}

	return nil
}
