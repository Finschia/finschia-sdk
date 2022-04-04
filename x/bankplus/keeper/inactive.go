package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/bankplus/types"
)

// Keys for bankplus store but this prefix must not be overlap with bank key prefix.
var inactiveAddrsKeyPrefix = []byte{0xa0}

// inactiveAddrKey key of a specific inactiveAddr from store
func inactiveAddrKey(addr sdk.AccAddress) []byte {
	return append(inactiveAddrsKeyPrefix, addr.Bytes()...)
}

//nolint:deadcode,unused
// isStoredInactiveAddr checks if the address is stored or not as blocked address
func (keeper BaseKeeper) isStoredInactiveAddr(ctx sdk.Context, address sdk.AccAddress) bool {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(inactiveAddrKey(address))
	return bz != nil
}

// addToInactiveAddr adds a blocked address to the store.
func (keeper BaseKeeper) addToInactiveAddr(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(keeper.storeKey)
	blockedCAddr := types.InactiveAddr{Address: address.String()}
	bz := keeper.cdc.MustMarshal(&blockedCAddr)
	store.Set(inactiveAddrKey(address), bz)
}

// deleteFromInactiveAddr deletes blocked address from store
func (keeper BaseKeeper) deleteFromInactiveAddr(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(inactiveAddrKey(address))
}

// loadAllInactiveAddrs loads all blocked address and set to `inactiveAddr`.
// This function is executed when the app is initiated and save all inactive address in caches
// in order to prevent to query to store in every time to send
func (keeper BaseKeeper) loadAllInactiveAddrs(ctx sdk.Context) {
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, inactiveAddrsKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var bAddr types.InactiveAddr
		keeper.cdc.MustUnmarshal(iterator.Value(), &bAddr)

		keeper.inactiveAddrs[bAddr.Address] = true
	}
}
