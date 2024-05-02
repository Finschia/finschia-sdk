package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/Finschia/finschia-sdk/codec"
	"github.com/Finschia/finschia-sdk/store/prefix"
	storetypes "github.com/Finschia/finschia-sdk/store/types"
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey

	config types.Config

	AccountKeeper
	BankKeeper

	swaps        []types.Swap
	swappedSlice []types.Swapped
}

func NewKeeper(cdc codec.BinaryCodec, storeKey storetypes.StoreKey, config types.Config, ak AccountKeeper, bk BankKeeper) Keeper {
	return Keeper{
		cdc,
		storeKey,
		config,
		ak,
		bk,
		make([]types.Swap, config.MaxSwaps),
		make([]types.Swapped, config.MaxSwaps),
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) MakeSwap(ctx sdk.Context, swap types.Swap) error {
	if err := swap.ValidateBasic(); err != nil {
		return err
	}

	isNewSwap := true
	if _, err := k.getSwap(ctx, swap.ToDenom, swap.FromDenom); err == nil {
		return sdkerrors.ErrInvalidRequest.Wrap("there is already a swap in reverse way, only one way swap allowed")
	}

	if _, err := k.getSwap(ctx, swap.FromDenom, swap.ToDenom); err == nil {
		isNewSwap = false
	}

	if !isNewSwap && !k.config.UpdateAllowed {
		return sdkerrors.ErrInvalidRequest.Wrap("update existing swap not allowed")
	}

	if err := k.increaseSwapCount(ctx); err != nil {
		return err
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

func (k Keeper) Swap(ctx sdk.Context, addr sdk.AccAddress, fromCoinAmount sdk.Coin, toDenom string) error {
	swap, err := k.getSwap(ctx, fromCoinAmount.Denom, toDenom)
	if err != nil {
		return err
	}

	if swap.GetFromDenom() != fromCoinAmount.GetDenom() {
		return sdkerrors.ErrInvalidRequest.Wrapf("denom mismatch, expected %s, got %s", swap.GetFromDenom(), fromCoinAmount.Denom)
	}

	if swap.GetToDenom() != toDenom {
		return sdkerrors.ErrInvalidRequest.Wrapf("denom mismatch, expected %s, got %s", swap.GetToDenom(), toDenom)
	}

	newAmount := fromCoinAmount.Amount.Mul(swap.SwapMultiple)
	newCoinAmount := sdk.NewCoin(toDenom, newAmount)
	if err := k.checkSwapCap(ctx, swap, newCoinAmount); err != nil {
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

	if err := k.updateSwapped(ctx, fromCoinAmount, newCoinAmount); err != nil {
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

func (k Keeper) SwapAll(ctx sdk.Context, addr sdk.AccAddress, fromDenom, toDenom string) error {
	balance := k.GetBalance(ctx, addr, fromDenom)
	if balance.IsZero() {
		return sdkerrors.ErrInsufficientFunds
	}

	if err := k.Swap(ctx, addr, balance, toDenom); err != nil {
		return err
	}
	return nil
}

func (k Keeper) getAllSwapped(ctx sdk.Context) []types.Swapped {
	idx := 0
	k.iterateAllSwapped(ctx, func(swapped types.Swapped) bool {
		k.swappedSlice[idx] = swapped
		idx++
		return false
	})
	if idx == 0 {
		return []types.Swapped{}
	}
	return k.swappedSlice
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

func (k Keeper) setSwapped(ctx sdk.Context, fromDenom, toDenom string, swapped types.Swapped) error {
	key := swappedKey(fromDenom, toDenom)
	bz, err := k.cdc.Marshal(&swapped)
	if err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(key, bz)
	return nil
}

func (k Keeper) getAllSwaps(ctx sdk.Context) []types.Swap {
	idx := 0
	k.iterateAllSwaps(ctx, func(swap types.Swap) bool {
		k.swaps[idx] = swap
		idx++
		return false
	})
	if idx == 0 {
		return []types.Swap{}
	}
	return k.swaps
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

func (k Keeper) updateSwapped(ctx sdk.Context, fromAmount, toAmount sdk.Coin) error {
	prevSwapped, err := k.getSwapped(ctx, fromAmount.Denom, toAmount.Denom)
	if err != nil {
		return err
	}

	updatedSwapped := &types.Swapped{
		FromCoinAmount: fromAmount.Add(prevSwapped.FromCoinAmount),
		ToCoinAmount:   toAmount.Add(prevSwapped.ToCoinAmount),
	}

	key := swappedKey(fromAmount.Denom, toAmount.Denom)
	bz, err := k.cdc.Marshal(updatedSwapped)
	if err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(key, bz)
	return nil
}

func (k Keeper) checkSwapCap(ctx sdk.Context, swap types.Swap, newCoinAmount sdk.Coin) error {
	swapped, err := k.getSwapped(ctx, swap.GetFromDenom(), swap.GetToDenom())
	if err != nil {
		return err
	}

	swapCap := swap.AmountCapForToDenom
	if swapCap.LT(swapped.ToCoinAmount.Add(newCoinAmount).Amount) {
		return types.ErrExceedSwappableToCoinAmount
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
