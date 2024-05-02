package keeper

import (
	"github.com/Finschia/finschia-sdk/store/prefix"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

func (k Keeper) MakeSwap(ctx sdk.Context, swap types.Swap) error {
	if err := swap.ValidateBasic(); err != nil {
		return err
	}

	isNewSwap := true
	if _, err := k.getSwap(ctx, swap.FromDenom, swap.ToDenom); err == nil {
		isNewSwap = false
	}

	if !isNewSwap && !k.config.UpdateAllowed {
		return errors.ErrInvalidRequest.Wrap("update existing swap not allowed")
	}

	if isNewSwap {
		if err := k.increaseSwapCount(ctx); err != nil {
			return err
		}
	}

	key := swapKey(swap.FromDenom, swap.ToDenom)
	bz, err := k.cdc.Marshal(&swap)
	if err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(key, bz)

	if isNewSwap {
		swapped := types.Swapped{
			FromCoinAmount: sdk.Coin{
				Denom:  swap.GetFromDenom(),
				Amount: sdk.ZeroInt(),
			},
			ToCoinAmount: sdk.Coin{
				Denom:  swap.GetToDenom(),
				Amount: sdk.ZeroInt(),
			},
		}
		if err := k.setSwapped(ctx, swap.GetFromDenom(), swap.GetToDenom(), swapped); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) increaseSwapCount(ctx sdk.Context) error {
	stats, err := k.getSwapStats(ctx)
	if err != nil {
		return err
	}

	prev := stats.SwapCount
	stats.SwapCount += 1
	if stats.SwapCount < prev {
		return types.ErrInvalidState.Wrap("overflow detected")
	}

	if err := k.setSwapStats(ctx, stats); err != nil {
		return err
	}
	return nil
}

func (k Keeper) getSwap(ctx sdk.Context, fromDenom, toDenom string) (types.Swap, error) {
	store := ctx.KVStore(k.storeKey)
	key := swapKey(fromDenom, toDenom)
	bz := store.Get(key)
	if bz == nil {
		return types.Swap{}, errors.ErrNotFound.Wrap("swap not found")
	}

	swap := types.Swap{}
	if err := k.cdc.Unmarshal(bz, &swap); err != nil {
		return types.Swap{}, err
	}

	return swap, nil
}

func (k Keeper) getAllSwaps(ctx sdk.Context) []types.Swap {
	swaps := []types.Swap{}
	k.iterateAllSwaps(ctx, func(swap types.Swap) bool {
		swaps = append(swaps, swap)
		return false
	})
	return swaps
}

func (k Keeper) iterateAllSwaps(ctx sdk.Context, cb func(swapped types.Swap) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	swapDataStore := prefix.NewStore(store, swapPrefix)

	iterator := swapDataStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		swap := types.Swap{}
		k.cdc.MustUnmarshal(iterator.Value(), &swap)
		if cb(swap) {
			break
		}
	}
}
