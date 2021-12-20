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
// isStoredInactiveAddr check if the address is stored or not as blocked address
func (keeper BaseKeeper) isStoredInactiveAddr(ctx sdk.Context, address sdk.AccAddress) bool {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(inactiveAddrKey(address))
	return bz != nil
}

// addToInactiveAddr add a blocked address to the store.
func (keeper BaseKeeper) addToInactiveAddr(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(keeper.storeKey)
	blockedCAddr := types.InactiveAddr{Address: address.String()}
	bz := keeper.cdc.MustMarshalBinaryBare(&blockedCAddr)
	store.Set(inactiveAddrKey(address), bz)
}

// deleteFromInactiveAddr delete blocked address from store
func (keeper BaseKeeper) deleteFromInactiveAddr(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(inactiveAddrKey(address))
}

// loadAllInactiveAddrs load all blocked address and set to `inactiveAddr`.
func (keeper BaseKeeper) loadAllInactiveAddrs(ctx sdk.Context) {
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, inactiveAddrsKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var bAddr types.InactiveAddr
		keeper.cdc.MustUnmarshalBinaryBare(iterator.Value(), &bAddr)

		keeper.inactiveAddrs[bAddr.Address] = true
	}
}
