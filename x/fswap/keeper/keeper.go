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

	AccountKeeper
	BankKeeper

	fromDenom    string
	toDenom      string
	swapInit     types.SwapInit
	swapMultiple sdk.Int
	swapCap      sdk.Int
}

func NewKeeper(cdc codec.BinaryCodec, storeKey storetypes.StoreKey, ak AccountKeeper, bk BankKeeper) Keeper {
	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		AccountKeeper: ak,
		BankKeeper:    bk,
		fromDenom:     "",
		toDenom:       "",
		swapInit:      types.SwapInit{},
		swapMultiple:  sdk.ZeroInt(),
		swapCap:       sdk.ZeroInt(),
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) SwapInit(ctx sdk.Context, swapInit types.SwapInit) error {
	if err := swapInit.ValidateBasic(); err != nil {
		return err
	}
	if k.hasBeenInitialized(ctx) {
		return types.ErrSwapCanNotBeInitializedTwice
	}
	if err := k.setSwapInit(ctx, swapInit); err != nil {
		return err
	}
	swapped := types.Swapped{
		FromCoinAmount: sdk.Coin{
			Denom:  swapInit.GetFromDenom(),
			Amount: sdk.ZeroInt(),
		},
		ToCoinAmount: sdk.Coin{
			Denom:  swapInit.GetToDenom(),
			Amount: sdk.ZeroInt(),
		},
	}
	if err := k.setSwapped(ctx, swapped); err != nil {
		return err
	}
	return nil
}

func (k Keeper) Swap(ctx sdk.Context, addr sdk.AccAddress, fromCoinAmount sdk.Coin, toDenom string) error {
	swapInit, err := k.getSwapInit(ctx)
	if err != nil {
		return err
	}

	if swapInit.GetFromDenom() != fromCoinAmount.GetDenom() {
		return sdkerrors.ErrInvalidRequest.Wrapf("denom mismatch, expected %s, got %s", swapInit.GetFromDenom(), fromCoinAmount.Denom)
	}

	if swapInit.GetToDenom() != toDenom {
		return sdkerrors.ErrInvalidRequest.Wrapf("denom mismatch, expected %s, got %s", swapInit.GetToDenom(), toDenom)
	}

	newAmount := fromCoinAmount.Amount.Mul(swapInit.SwapMultiple)
	newCoinAmount := sdk.NewCoin(toDenom, newAmount)
	if err := k.checkSwapCap(ctx, newCoinAmount); err != nil {
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

func (k Keeper) setSwapInit(ctx sdk.Context, swapInit types.SwapInit) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&swapInit)
	if err != nil {
		return err
	}
	store.Set(swapInitKey(swapInit.ToDenom), bz)
	return nil
}

func (k Keeper) getAllSwapped(ctx sdk.Context) []types.Swapped {
	swappedSlice := make([]types.Swapped, 0) // TODO(bjs)
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

func (k Keeper) getSwapped(ctx sdk.Context) (types.Swapped, error) {
	toDenom, err := k.getToDenom(ctx)
	if err != nil {
		return types.Swapped{}, err
	}

	store := ctx.KVStore(k.storeKey)
	key := swappedKey(toDenom)
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
	store := ctx.KVStore(k.storeKey)
	key := swappedKey(swapped.ToCoinAmount.Denom)
	bz, err := k.cdc.Marshal(&swapped)
	if err != nil {
		return err
	}

	store.Set(key, bz)
	return nil
}

func (k Keeper) getAllSwapInits(ctx sdk.Context) []types.SwapInit {
	swapInits := make([]types.SwapInit, 0)
	k.iterateAllSwapInits(ctx, func(swapInit types.SwapInit) bool {
		swapInits = append(swapInits, swapInit)
		return false
	})
	return swapInits
}

func (k Keeper) iterateAllSwapInits(ctx sdk.Context, cb func(swapped types.SwapInit) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	swapInitDataStore := prefix.NewStore(store, swapInitPrefix)

	iterator := swapInitDataStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		swapInit := types.SwapInit{}
		k.cdc.MustUnmarshal(iterator.Value(), &swapInit)
		if cb(swapInit) {
			break
		}
	}
}

func (k Keeper) getSwappableNewCoinAmount(ctx sdk.Context) (sdk.Coin, error) {
	swapCap, err := k.getSwapCap(ctx)
	if err != nil {
		return sdk.Coin{}, err
	}
	swapped, err := k.getSwapped(ctx)
	if err != nil {
		return sdk.Coin{}, err
	}
	denom, err := k.getToDenom(ctx)
	if err != nil {
		return sdk.Coin{}, err
	}

	remainingAmount := swapCap.Sub(swapped.GetToCoinAmount().Amount)

	return sdk.NewCoin(denom, remainingAmount), nil
}

func (k Keeper) getFromDenom(ctx sdk.Context) (string, error) {
	if len(k.fromDenom) > 0 {
		return k.fromDenom, nil
	}

	swapInit, err := k.getSwapInit(ctx)
	if err != nil {
		return "", err
	}

	k.fromDenom = swapInit.GetFromDenom()
	return k.fromDenom, nil
}

func (k Keeper) getToDenom(ctx sdk.Context) (string, error) {
	if len(k.toDenom) > 0 {
		return k.toDenom, nil
	}

	swapInit, err := k.getSwapInit(ctx)
	if err != nil {
		return "", err
	}

	k.toDenom = swapInit.GetToDenom()
	return k.toDenom, nil
}

func (k Keeper) getSwapMultiple(ctx sdk.Context) (sdk.Int, error) {
	if k.swapMultiple.IsPositive() {
		return k.swapMultiple, nil
	}

	swapInit, err := k.getSwapInit(ctx)
	if err != nil {
		return sdk.Int{}, err
	}

	k.swapMultiple = swapInit.SwapMultiple
	return k.swapMultiple, nil
}

func (k Keeper) getSwapCap(ctx sdk.Context) (sdk.Int, error) {
	if k.swapCap.IsPositive() {
		return k.swapCap, nil
	}

	swapInit, err := k.getSwapInit(ctx)
	if err != nil {
		return sdk.Int{}, err
	}

	k.swapCap = swapInit.AmountCapForToDenom
	return k.swapCap, nil
}

func (k Keeper) getSwapInit(ctx sdk.Context) (types.SwapInit, error) {
	if !k.swapInit.IsEmpty() {
		return k.swapInit, nil
	}

	swapInits := k.getAllSwapInits(ctx)
	if len(swapInits) == 0 {
		return types.SwapInit{}, types.ErrSwapNotInitialized
	}

	k.swapInit = swapInits[0]
	return k.swapInit, nil
}

func (k Keeper) updateSwapped(ctx sdk.Context, fromAmount, toAmount sdk.Coin) error {
	prevSwapped, err := k.getSwapped(ctx)
	if err != nil {
		return err
	}

	updatedSwapped := &types.Swapped{
		FromCoinAmount: fromAmount.Add(prevSwapped.FromCoinAmount),
		ToCoinAmount:   toAmount.Add(prevSwapped.ToCoinAmount),
	}
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(updatedSwapped)
	if err != nil {
		return err
	}

	key := swappedKey(toAmount.Denom)
	store.Set(key, bz)
	return nil
}

func (k Keeper) checkSwapCap(ctx sdk.Context, newCoinAmount sdk.Coin) error {
	swapped, err := k.getSwapped(ctx)
	if err != nil {
		return err
	}

	swapCap, err := k.getSwapCap(ctx)
	if err != nil {
		return err
	}

	if swapCap.LT(swapped.ToCoinAmount.Add(newCoinAmount).Amount) {
		return types.ErrExceedSwappableToCoinAmount
	}
	return nil
}

func (k Keeper) hasBeenInitialized(ctx sdk.Context) bool {
	inits := k.getAllSwapInits(ctx)
	return len(inits) > 0
}
