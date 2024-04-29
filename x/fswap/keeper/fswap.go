package keeper

import (
	"github.com/Finschia/finschia-sdk/store/prefix"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

func (k Keeper) GetSwapped(ctx sdk.Context, toDenom string) (types.Swapped, error) {
	store := ctx.KVStore(k.storeKey)
	key := swappedKey(toDenom)
	bz := store.Get(key)
	if bz == nil {
		return types.Swapped{
			OldCoinAmount: sdk.Coin{Denom: "", Amount: sdk.ZeroInt()},
			NewCoinAmount: sdk.Coin{Denom: "", Amount: sdk.ZeroInt()},
		}, nil
	}
	swapped := types.Swapped{}
	if err := k.cdc.Unmarshal(bz, &swapped); err != nil {
		return types.Swapped{}, err
	}
	return swapped, nil
}

func (k Keeper) GetAllSwapped(ctx sdk.Context) []types.Swapped {
	swappedSlice := make([]types.Swapped, 0)
	k.IterateAllSwapped(ctx, func(fswapInit types.Swapped) bool {
		swappedSlice = append(swappedSlice, fswapInit)
		return false
	})
	return swappedSlice
}

func (k Keeper) IterateAllSwapped(ctx sdk.Context, cb func(swapped types.Swapped) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	swappedDataStore := prefix.NewStore(store, swappedKeyPrefix)

	iterator := swappedDataStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		swapped := types.Swapped{}
		k.cdc.MustUnmarshal(iterator.Value(), &swapped)
		if cb(swapped) {
			break
		}
	}
}

func (k Keeper) SetSwapped(ctx sdk.Context, swapped types.Swapped) error {
	store := ctx.KVStore(k.storeKey)
	key := swappedKey(swapped.NewCoinAmount.Denom)
	bz, err := k.cdc.Marshal(&swapped)
	if err != nil {
		return err
	}
	store.Set(key, bz)
	return nil
}

func (k Keeper) GetFswapInit(ctx sdk.Context, toDenom string) (types.FswapInit, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(fswapInitKey(toDenom))
	if bz == nil {
		return types.FswapInit{}, types.ErrFswapInitNotFound
	}
	fswapInit := types.FswapInit{}
	err := k.cdc.Unmarshal(bz, &fswapInit)
	if err != nil {
		return types.FswapInit{}, err
	}
	return fswapInit, nil
}

func (k Keeper) GetAllFswapInits(ctx sdk.Context) []types.FswapInit {
	fswapInits := make([]types.FswapInit, 0)
	k.IterateAllFswapInits(ctx, func(fswapInit types.FswapInit) bool {
		fswapInits = append(fswapInits, fswapInit)
		return false
	})
	return fswapInits
}

func (k Keeper) IterateAllFswapInits(ctx sdk.Context, cb func(swapped types.FswapInit) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	fswapInitDataStore := prefix.NewStore(store, fswapInitPrefix)

	iterator := fswapInitDataStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		fswapInit := types.FswapInit{}
		k.cdc.MustUnmarshal(iterator.Value(), &fswapInit)
		if cb(fswapInit) {
			break
		}
	}
}

func (k Keeper) SetFswapInit(ctx sdk.Context, fswapInit types.FswapInit) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&fswapInit)
	if err != nil {
		return err
	}
	store.Set(fswapInitKey(fswapInit.ToDenom), bz)
	return nil
}
