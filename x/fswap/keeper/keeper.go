package keeper

import (
	"errors"

	"github.com/Finschia/finschia-sdk/codec"
	storetypes "github.com/Finschia/finschia-sdk/store/types"
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

// TODO: move const to proper place
const (
	OldDenom     = "cony"
	NewDenom     = "TBD"
	SwapMultiple = 123
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey

	AccountKeeper
	BankKeeper
}

func NewKeeper(cdc codec.BinaryCodec, storeKey storetypes.StoreKey, ak AccountKeeper, bk BankKeeper) *Keeper {
	return &Keeper{
		cdc,
		storeKey,
		ak,
		bk,
	}
}

func (k Keeper) Swap(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coin) error {
	// 1. Check Balance
	// 2. Check swapCap
	// 3. Send coin to module
	// 4. Burn & Mint TBD coins(with multiple)
	// 5. Send back coins to addr
	// 6. Update swapped state
	if ok := k.HasBalance(ctx, addr, amount); !ok {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "not enough amount: %s", amount)
	}

	multipleAmount := amount.Amount.Mul(sdk.NewInt(SwapMultiple))
	newAmount := sdk.NewCoin(NewDenom, multipleAmount)
	if err := k.checkSwapCap(ctx, newAmount); err != nil {
		return err
	}

	moduleAddr := k.GetModuleAddress(types.ModuleName)
	if err := k.SendCoins(ctx, addr, moduleAddr, sdk.NewCoins(amount)); err != nil {
		return err
	}

	if err := k.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(amount)); err != nil {
		return err
	}

	if err := k.MintCoins(ctx, types.ModuleName, sdk.NewCoins(newAmount)); err != nil {
		return err
	}

	if err := k.SendCoins(ctx, moduleAddr, addr, sdk.NewCoins(newAmount)); err != nil {
		return err
	}

	if err := k.updateSwapped(ctx, amount, newAmount); err != nil {
		return err
	}

	if err := ctx.EventManager().EmitTypedEvent(&types.EventSwapCoins{
		Address:       addr.String(),
		OldCoinAmount: amount,
		NewCoinAmount: newAmount,
	}); err != nil {
		panic(err)
	}
	return nil
}

func (k Keeper) SwapAll(ctx sdk.Context, addr sdk.AccAddress) error {
	balance := k.GetBalance(ctx, addr, OldDenom)
	if balance.IsZero() {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "zero balance for %s", OldDenom)
	}

	if err := k.Swap(ctx, addr, balance); err != nil {
		return err
	}
	return nil
}

func (k Keeper) updateSwapped(ctx sdk.Context, oldAmount sdk.Coin, newAmount sdk.Coin) error {
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

func (k Keeper) setSwapCap(ctx sdk.Context) error {
	// TODO(bjs): how to set it only once
	//store := ctx.KVStore(k.storeKey)
	return nil
}

func (k Keeper) getSwapCap(ctx sdk.Context) (sdk.Coin, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(swapCapKey())
	if bz == nil {
		return sdk.Coin{}, errors.New("swap cap not found")
	}

	swapCap := sdk.Coin{}
	if err := k.cdc.Unmarshal(bz, &swapCap); err != nil {
		return sdk.Coin{}, err
	}
	return swapCap, nil
}

func (k Keeper) checkSwapCap(ctx sdk.Context, amountNewCoin sdk.Coin) error {
	swapCap, err := k.getSwapCap(ctx)
	if err != nil {
		return err
	}
	swapped, err := k.GetSwapped(ctx)
	if err != nil {
		return err
	}
	if swapCap.IsLT(swapped.GetNewCoinAmount().Add(amountNewCoin)) {
		return errors.New("not enough swap coin amount")
	}
	return nil
}
