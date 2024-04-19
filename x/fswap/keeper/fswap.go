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

func (k Keeper) GetTotalSupply(ctx sdk.Context) sdk.DecCoin {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte{types.TotalSupplyKey})
	var totalSupply sdk.DecCoin
	if bz == nil {
		panic(types.ErrParamsNotFound)
	}
	k.cdc.MustUnmarshal(bz, &totalSupply)
	return totalSupply
}

func (k Keeper) SetTotalSupply(ctx sdk.Context, totalSupply sdk.DecCoin) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&totalSupply)
	if err != nil {
		return err
	}
	store.Set([]byte{types.TotalSupplyKey}, bz)
	return nil
}
