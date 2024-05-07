package keeper

import (
	"github.com/Finschia/finschia-sdk/store/prefix"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/types/errors"
	bank "github.com/Finschia/finschia-sdk/x/bank/types"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

func (k Keeper) MakeSwap(ctx sdk.Context, swap types.Swap, toDenomMetadata bank.Metadata) error {
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

	stats, err := k.getSwapStats(ctx)
	if err != nil {
		return err
	}

	if int(stats.SwapCount) > k.config.MaxSwaps && !k.isUnlimited() {
		return types.ErrCanNotHaveMoreSwap.Wrapf("cannot make more swaps, max swaps is %d", k.config.MaxSwaps)
	}

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
		if err := k.setSwapped(ctx, swapped); err != nil {
			return err
		}
	}

	if err := k.setSwap(ctx, swap); err != nil {
		return err
	}

	existingMetadata, ok := k.GetDenomMetaData(ctx, swap.ToDenom)
	if !ok {
		k.SetDenomMetaData(ctx, toDenomMetadata)
		return nil
	}
	if !denomMetadataEqual(existingMetadata, toDenomMetadata) {
		return errors.ErrInvalidRequest.Wrap("changing existing metadata not allowed")
	}

	return nil
}

func denomMetadataEqual(metadata, otherMetadata bank.Metadata) bool {
	if metadata.Description != otherMetadata.Description {
		return false
	}
	if len(metadata.DenomUnits) != len(otherMetadata.DenomUnits) {
		return false
	}
	for i, unit := range metadata.DenomUnits {
		if unit.Denom != otherMetadata.DenomUnits[i].Denom {
			return false
		}
	}
	if metadata.Base != otherMetadata.Base {
		return false
	}
	if metadata.Display != otherMetadata.Display {
		return false
	}
	if metadata.Name != otherMetadata.Name {
		return false
	}
	if metadata.Symbol != otherMetadata.Symbol {
		return false
	}
	return true
}

func (k Keeper) increaseSwapCount(ctx sdk.Context) error {
	stats, err := k.getSwapStats(ctx)
	if err != nil {
		return err
	}

	prev := stats.SwapCount
	stats.SwapCount++
	if stats.SwapCount < prev {
		return types.ErrInvalidState.Wrap("overflow detected")
	}

	if err := k.setSwapStats(ctx, stats); err != nil {
		return err
	}
	return nil
}

func (k Keeper) setSwap(ctx sdk.Context, swap types.Swap) error {
	key := swapKey(swap.FromDenom, swap.ToDenom)
	bz, err := k.cdc.Marshal(&swap)
	if err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(key, bz)
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
