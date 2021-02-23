package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link-modules/x/token/internal/types"
)

// For the Token module
type BankKeeper interface {
	GetBalance(ctx sdk.Context, addr sdk.AccAddress) sdk.Int
	SetBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Int) error
	HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Int) bool

	SubtractBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Int) (sdk.Int, error)
	AddBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Int) (sdk.Int, error)
	Send(ctx sdk.Context, from, to sdk.AccAddress, amt sdk.Int) error
}

var _ BankKeeper = (*Keeper)(nil)

func (k Keeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress) sdk.Int {
	acc, err := k.GetAccount(ctx, addr)
	if err != nil {
		return sdk.ZeroInt()
	}
	return acc.GetBalance()
}

func (k Keeper) SetBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Int) error {
	acc, err := k.GetAccount(ctx, addr)
	if err != nil {
		return err
	}
	acc = acc.SetBalance(amt)
	err = k.UpdateAccount(ctx, acc)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Int) bool {
	return k.GetBalance(ctx, addr).GTE(amt)
}

func (k Keeper) Send(ctx sdk.Context, from, to sdk.AccAddress, amt sdk.Int) error {
	if amt.IsNegative() {
		return sdkerrors.Wrap(types.ErrInvalidAmount, "send amount must be positive")
	}

	_, err := k.SubtractBalance(ctx, from, amt)
	if err != nil {
		return err
	}

	_, err = k.AddBalance(ctx, to, amt)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) SubtractBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Int) (sdk.Int, error) {
	return k.subtractBalance(ctx, addr, amt)
}
func (k Keeper) subtractBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Int) (sdk.Int, error) {
	acc, err := k.GetAccount(ctx, addr)
	if err != nil {
		return sdk.ZeroInt(), sdkerrors.Wrapf(types.ErrInsufficientBalance, fmt.Sprintf("insufficient account funds for token [%s]; 0 < %s", k.getContractID(ctx), amt))
	}
	oldBalance := acc.GetBalance()
	newBalance := oldBalance.Sub(amt)
	if newBalance.IsNegative() {
		return amt, sdkerrors.Wrapf(types.ErrInsufficientBalance, "insufficient account funds for token [%s]; %s < %s", k.getContractID(ctx), oldBalance, amt)
	}
	acc = acc.SetBalance(newBalance)
	err = k.UpdateAccount(ctx, acc)
	if err != nil {
		return amt, err
	}
	return newBalance, nil
}

func (k Keeper) AddBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Int) (sdk.Int, error) {
	return k.addBalance(ctx, addr, amt)
}

func (k Keeper) addBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Int) (sdk.Int, error) {
	acc, err := k.GetOrNewAccount(ctx, addr)
	if err != nil {
		return amt, err
	}
	oldBalance := acc.GetBalance()
	newBalance := oldBalance.Add(amt)
	if newBalance.IsNegative() {
		return amt, sdkerrors.Wrapf(types.ErrInsufficientBalance, "insufficient account funds for token [%s]; %s < %s", k.getContractID(ctx), oldBalance, amt)
	}
	acc = acc.SetBalance(newBalance)
	err = k.UpdateAccount(ctx, acc)
	if err != nil {
		return amt, err
	}
	return newBalance, nil
}
