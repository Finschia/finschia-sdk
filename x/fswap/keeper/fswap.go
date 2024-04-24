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

func (k Keeper) GetSwappableNewCoinAmount(ctx sdk.Context) sdk.Coin {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte{types.SwappableNewCoinAmountKey})
	var SwappableNewCoinAmount sdk.Coin
	if bz == nil {
		panic(types.ErrSwappableNewCoinAmountNotFound)
	}
	k.cdc.MustUnmarshal(bz, &SwappableNewCoinAmount)
	return SwappableNewCoinAmount
}

func (k Keeper) SetSwappableNewCoinAmount(ctx sdk.Context, SwappableNewCoinAmount sdk.Coin) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&SwappableNewCoinAmount)
	if err != nil {
		return err
	}
	store.Set([]byte{types.SwappableNewCoinAmountKey}, bz)
	return nil
}
