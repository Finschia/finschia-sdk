package keeper

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

// GetSwapped get all parameters as types.Swapped
func (k Keeper) GetSwapped(ctx sdk.Context) types.Swapped {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte{types.SwappedKey})
	var swapped types.Swapped

	if bz == nil {
		panic(types.ErrSwappedNotFound)
	}
	k.cdc.MustUnmarshal(bz, &swapped)
	return swapped
}

// SetSwapped set the types.Swapped
func (k Keeper) SetSwapped(ctx sdk.Context, swapped types.Swapped) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&swapped)
	if err != nil {
		return err
	}
	store.Set([]byte{types.SwappedKey}, bz)
	return nil
}

// GetFswapInit get all parameters as types.FswapInit
func (k Keeper) GetFswapInit(ctx sdk.Context) types.FswapInit {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte{types.FswapInitKey})
	var fswapInit types.FswapInit
	if bz == nil {
		panic(types.ErrFswapInitNotFound)
	}
	k.cdc.MustUnmarshal(bz, &fswapInit)
	return fswapInit
}

// SetfswapInit set the fswapInit
func (k Keeper) SetFswapInit(ctx sdk.Context, fswapInit types.FswapInit) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&fswapInit)
	if err != nil {
		return err
	}
	store.Set([]byte{types.FswapInitKey}, bz)
	return nil
}
