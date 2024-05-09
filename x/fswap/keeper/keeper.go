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
