package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/bankplus/types"
)

// Keys for bankplus store but this prefix must not be overlap with bank key prefix.
var blockedAddrsKeyPrefix = []byte{0xa0}

// blockedAddrKey key of a specific blockedAddr from store
func blockedAddrKey(addr sdk.AccAddress) []byte {
	return append(blockedAddrsKeyPrefix, addr.Bytes()...)
}

//nolint:deadcode,unused
// isAddedBlockedAddr check if the address is stored or not as blocked address
func (keeper BaseKeeper) isAddedBlockedAddr(ctx sdk.Context, address sdk.AccAddress) bool {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(blockedAddrKey(address))
	return bz != nil
}

// addBlockedAddr add a blocked address to the store.
func (keeper BaseKeeper) addBlockedAddr(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(keeper.storeKey)
	blockedCAddr := types.BlockedAddr{Address: address.String()}
	bz := keeper.cdc.MustMarshalBinaryBare(&blockedCAddr)
	store.Set(blockedAddrKey(address), bz)
}

// deleteBlockedAddr delete blocked address from store
func (keeper BaseKeeper) deleteBlockedAddr(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(blockedAddrKey(address))
}

// loadAllBlockedAddrs load all blocked address and set to `blockedAddr`.
func (keeper BaseKeeper) loadAllBlockedAddrs(ctx sdk.Context) {
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, blockedAddrsKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var bAddr types.BlockedAddr
		keeper.cdc.MustUnmarshalBinaryBare(iterator.Value(), &bAddr)

		keeper.blockedAddrs[bAddr.Address] = true
	}
}
