package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/Finschia/finschia-sdk/codec"
	"github.com/Finschia/finschia-sdk/store/prefix"
	storetypes "github.com/Finschia/finschia-sdk/store/types"
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	bank "github.com/Finschia/finschia-sdk/x/bank/types"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

type Keeper struct {
	cdc       codec.BinaryCodec
	storeKey  storetypes.StoreKey
	config    types.Config
	authority string
	BankKeeper
}

func NewKeeper(cdc codec.BinaryCodec, storeKey storetypes.StoreKey, config types.Config, authority string, bk BankKeeper) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic("authority is not a valid acc address")
	}

	// authority is x/foundation module account for now.
	if authority != types.DefaultAuthority().String() {
		panic("x/foundation authority must be the module account")
	}

	return Keeper{
		cdc,
		storeKey,
		config,
		authority,
		bk,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Swap(ctx sdk.Context, addr sdk.AccAddress, fromCoinAmount sdk.Coin, toDenom string) error {
	swap, err := k.getSwap(ctx, fromCoinAmount.Denom, toDenom)
	if err != nil {
		return err
	}

	newCoinAmountInt := CalcSwap(swap.SwapRate, fromCoinAmount.Amount)
	newCoinAmount := sdk.NewCoin(toDenom, newCoinAmountInt)
	swapped, err := k.getSwapped(ctx, swap.GetFromDenom(), swap.GetToDenom())
	if err != nil {
		return err
	}

	updateSwapped, err := k.updateSwapped(ctx, swapped, fromCoinAmount, newCoinAmount)
	if err != nil {
		return err
	}

	if err := k.checkSwapCap(swap, updateSwapped); err != nil {
		return err
	}

	if err := k.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, sdk.NewCoins(fromCoinAmount)); err != nil {
		return err
	}

	if err := k.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(fromCoinAmount)); err != nil {
		return err
	}

	if err := k.MintCoins(ctx, types.ModuleName, sdk.NewCoins(newCoinAmount)); err != nil {
		return err
	}

	if err := k.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.NewCoins(newCoinAmount)); err != nil {
		return err
	}

	if err := ctx.EventManager().EmitTypedEvent(&types.EventSwapCoins{
		Address:        addr.String(),
		FromCoinAmount: fromCoinAmount,
		ToCoinAmount:   newCoinAmount,
	}); err != nil {
		return err
	}
	return nil
}

func (k Keeper) SetSwap(ctx sdk.Context, swap types.Swap, toDenomMetadata bank.Metadata) error {
	isNewSwap := true
	if _, err := k.getSwap(ctx, swap.FromDenom, swap.ToDenom); err == nil {
		isNewSwap = false
	}

	if !isNewSwap && !k.config.UpdateAllowed {
		return sdkerrors.ErrInvalidRequest.Wrap("update existing swap not allowed")
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
		return sdkerrors.ErrInvalidRequest.Wrap("changing existing metadata not allowed")
	}

	return nil
}

func (k Keeper) getAllSwapped(ctx sdk.Context) []types.Swapped {
	swappedSlice := []types.Swapped{}
	k.iterateAllSwapped(ctx, func(swapped types.Swapped) bool {
		swappedSlice = append(swappedSlice, swapped)
		return false
	})
	return swappedSlice
}

func (k Keeper) iterateAllSwapped(ctx sdk.Context, cb func(swapped types.Swapped) (stop bool)) {
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

func (k Keeper) getSwapped(ctx sdk.Context, fromDenom, toDenom string) (types.Swapped, error) {
	store := ctx.KVStore(k.storeKey)
	key := swappedKey(fromDenom, toDenom)
	bz := store.Get(key)
	if bz == nil {
		return types.Swapped{}, types.ErrSwappedNotFound
	}

	swapped := types.Swapped{}
	if err := k.cdc.Unmarshal(bz, &swapped); err != nil {
		return types.Swapped{}, err
	}
	return swapped, nil
}

func (k Keeper) setSwapped(ctx sdk.Context, swapped types.Swapped) error {
	key := swappedKey(swapped.FromCoinAmount.Denom, swapped.ToCoinAmount.Denom)
	bz, err := k.cdc.Marshal(&swapped)
	if err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(key, bz)
	return nil
}

func (k Keeper) getSwappableNewCoinAmount(ctx sdk.Context, fromDenom, toDenom string) (sdk.Coin, error) {
	swap, err := k.getSwap(ctx, fromDenom, toDenom)
	if err != nil {
		return sdk.Coin{}, err
	}

	swapped, err := k.getSwapped(ctx, fromDenom, toDenom)
	if err != nil {
		return sdk.Coin{}, err
	}

	swapCap := swap.AmountCapForToDenom
	remainingAmount := swapCap.Sub(swapped.GetToCoinAmount().Amount)

	return sdk.NewCoin(toDenom, remainingAmount), nil
}

func (k Keeper) updateSwapped(ctx sdk.Context, curSwapped types.Swapped, fromAmount, toAmount sdk.Coin) (types.Swapped, error) {
	updatedSwapped := types.Swapped{
		FromCoinAmount: fromAmount.Add(curSwapped.FromCoinAmount),
		ToCoinAmount:   toAmount.Add(curSwapped.ToCoinAmount),
	}

	key := swappedKey(fromAmount.Denom, toAmount.Denom)
	bz, err := k.cdc.Marshal(&updatedSwapped)
	if err != nil {
		return types.Swapped{}, err
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(key, bz)
	return updatedSwapped, nil
}

func (k Keeper) checkSwapCap(swap types.Swap, swapped types.Swapped) error {
	swapCap := swap.AmountCapForToDenom
	if swapCap.LT(swapped.ToCoinAmount.Amount) {
		return types.ErrExceedSwappableToCoinAmount
	}
	return nil
}

func (k Keeper) getSwapStats(ctx sdk.Context) (types.SwapStats, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(swapStatsKey)
	if bz == nil {
		return types.SwapStats{}, nil
	}

	stats := types.SwapStats{}
	err := k.cdc.Unmarshal(bz, &stats)
	if err != nil {
		return types.SwapStats{}, err
	}
	return stats, nil
}

func (k Keeper) setSwapStats(ctx sdk.Context, stats types.SwapStats) error {
	bz, err := k.cdc.Marshal(&stats)
	if err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(swapStatsKey, bz)
	return nil
}

func (k Keeper) validateAuthority(authority string) error {
	if authority != k.authority {
		return sdkerrors.ErrUnauthorized.Wrapf("invalid authority; expected %s, got %s", k.authority, authority)
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
	stats.SwapCount += 1
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
		return types.Swap{}, sdkerrors.ErrNotFound.Wrap("swap not found")
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
