package keeper

import (
	"fmt"
	"strings"

	"github.com/Finschia/finschia-sdk/codec"
	storetypes "github.com/Finschia/finschia-sdk/store/types"
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/x/fswap/config"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey

	config.FswapConfig

	AccountKeeper
	BankKeeper
}

func NewKeeper(cdc codec.BinaryCodec, storeKey storetypes.StoreKey, cfg config.FswapConfig, ak AccountKeeper, bk BankKeeper) Keeper {
	if len(strings.TrimSpace(cfg.NewDenom())) == 0 {
		panic("new denom must be provided")
	}
	// if strings.Compare(strings.TrimSpace(cfg.NewDenom()), cfg.SwapCap().Denom) != 0 {
	//	panic("new denom does not match new denom")
	//}
	return Keeper{
		cdc,
		storeKey,
		cfg,
		ak,
		bk,
	}
}

func (k Keeper) Swap(ctx sdk.Context, addr sdk.AccAddress, oldCoin sdk.Coin) error {
	if ok := k.HasBalance(ctx, addr, oldCoin); !ok {
		return sdkerrors.ErrInsufficientFunds
	}

	newAmount := oldCoin.Amount.Mul(k.SwapMultiple())
	newCoin := sdk.NewCoin(k.NewDenom(), newAmount)
	if err := k.checkSwapCap(ctx, newAmount); err != nil {
		return err
	}

	moduleAddr := k.GetModuleAddress(types.ModuleName)
	if err := k.SendCoins(ctx, addr, moduleAddr, sdk.NewCoins(oldCoin)); err != nil {
		return err
	}

	if err := k.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(oldCoin)); err != nil {
		return err
	}

	if err := k.MintCoins(ctx, types.ModuleName, sdk.NewCoins(newCoin)); err != nil {
		return err
	}

	if err := k.SendCoins(ctx, moduleAddr, addr, sdk.NewCoins(newCoin)); err != nil {
		return err
	}

	if err := k.updateSwapped(ctx, oldCoin.Amount, newAmount); err != nil {
		return err
	}

	if err := ctx.EventManager().EmitTypedEvent(&types.EventSwapCoins{
		Address:       addr.String(),
		OldCoinAmount: oldCoin.Amount,
		NewCoinAmount: newCoin.Amount,
	}); err != nil {
		return err
	}
	return nil
}

func (k Keeper) SwapAll(ctx sdk.Context, addr sdk.AccAddress) error {
	balance := k.GetBalance(ctx, addr, k.OldDenom())
	if balance.IsZero() {
		return sdkerrors.ErrInsufficientFunds
	}

	if err := k.Swap(ctx, addr, balance); err != nil {
		return err
	}
	return nil
}

func (k Keeper) updateSwapped(ctx sdk.Context, oldAmount, newAmount sdk.Int) error {
	prevSwapped, err := k.GetSwapped(ctx)
	if err != nil {
		return err
	}
	updatedSwapped := &types.Swapped{
		OldCoinAmount: oldAmount.Add(prevSwapped.OldCoinAmount),
		NewCoinAmount: newAmount.Add(prevSwapped.NewCoinAmount),
	}

	store := ctx.KVStore(k.storeKey)
	key := swappedKey()
	bz, err := k.cdc.Marshal(updatedSwapped)
	if err != nil {
		return err
	}
	store.Set(key, bz)
	return nil
}

func (k Keeper) GetSwapped(ctx sdk.Context) (types.Swapped, error) {
	store := ctx.KVStore(k.storeKey)
	key := swappedKey()
	bz := store.Get(key)
	swapped := types.Swapped{}
	if err := k.cdc.Unmarshal(bz, &swapped); err != nil {
		return types.Swapped{}, err
	}
	return swapped, nil
}

func (k Keeper) SetSwapped(ctx sdk.Context, swapped types.Swapped) error {
	store := ctx.KVStore(k.storeKey)
	key := swappedKey()
	bz, err := k.cdc.Marshal(&swapped)
	if err != nil {
		return err
	}
	store.Set(key, bz)
	return nil
}

func (k Keeper) checkSwapCap(ctx sdk.Context, newCoinAmount sdk.Int) error {
	swapped, err := k.GetSwapped(ctx)
	if err != nil {
		return err
	}

	if k.SwapCap().LT(swapped.NewCoinAmount.Add(newCoinAmount)) {
		return fmt.Errorf("cann't swap more because of swapCap limit, amount=%s", newCoinAmount.String())
	}

	return nil
}

func (k Keeper) GetSwappableNewCoinAmount(ctx sdk.Context) sdk.Coin {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte{types.SwappableNewCoinAmountKey})
	var swappableNewCoinAmount sdk.Coin
	if bz == nil {
		panic(types.ErrSwappableNewCoinAmountNotFound)
	}
	k.cdc.MustUnmarshal(bz, &swappableNewCoinAmount)
	return swappableNewCoinAmount
}

func (k Keeper) SetSwappableNewCoinAmount(ctx sdk.Context, swappableNewCoinAmount sdk.Coin) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&swappableNewCoinAmount)
	if err != nil {
		return err
	}
	store.Set([]byte{types.SwappableNewCoinAmountKey}, bz)
	return nil
}
