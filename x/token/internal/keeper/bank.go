package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
)

//For the Token module
type BankKeeper interface {
	GetBalance(ctx sdk.Context, contractID string, addr sdk.AccAddress) sdk.Int
	SetBalance(ctx sdk.Context, contractID string, addr sdk.AccAddress, amt sdk.Int) sdk.Error
	HasBalance(ctx sdk.Context, contractID string, addr sdk.AccAddress, amt sdk.Int) bool

	SubtractBalance(ctx sdk.Context, contractID string, addr sdk.AccAddress, amt sdk.Int) (sdk.Int, sdk.Error)
	AddBalance(ctx sdk.Context, contractID string, addr sdk.AccAddress, amt sdk.Int) (sdk.Int, sdk.Error)
	Send(ctx sdk.Context, contractID string, from, to sdk.AccAddress, amt sdk.Int) sdk.Error
}

var _ BankKeeper = (*Keeper)(nil)

func (k Keeper) GetBalance(ctx sdk.Context, contractID string, addr sdk.AccAddress) sdk.Int {
	acc, err := k.GetAccount(ctx, contractID, addr)
	if err != nil {
		return sdk.ZeroInt()
	}
	return acc.GetBalance()
}

func (k Keeper) SetBalance(ctx sdk.Context, contractID string, addr sdk.AccAddress, amt sdk.Int) sdk.Error {
	acc, err := k.GetAccount(ctx, contractID, addr)
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

func (k Keeper) HasBalance(ctx sdk.Context, contractID string, addr sdk.AccAddress, amt sdk.Int) bool {
	return k.GetBalance(ctx, contractID, addr).GTE(amt)
}

func (k Keeper) Send(ctx sdk.Context, contractID string, from, to sdk.AccAddress, amt sdk.Int) sdk.Error {
	if amt.IsNegative() {
		return types.ErrInvalidAmount(types.DefaultCodespace, "send amount must be positive")
	}

	_, err := k.SubtractBalance(ctx, contractID, from, amt)
	if err != nil {
		return err
	}

	_, err = k.AddBalance(ctx, contractID, to, amt)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) SubtractBalance(ctx sdk.Context, contractID string, addr sdk.AccAddress, amt sdk.Int) (sdk.Int, sdk.Error) {
	return k.subtractBalance(ctx, contractID, addr, amt)
}
func (k Keeper) subtractBalance(ctx sdk.Context, contractID string, addr sdk.AccAddress, amt sdk.Int) (sdk.Int, sdk.Error) {
	acc, err := k.GetAccount(ctx, contractID, addr)
	if err != nil {
		return sdk.ZeroInt(), types.ErrInsufficientBalance(types.DefaultCodespace, fmt.Sprintf("insufficient account funds for token [%s]; 0 < %s", contractID, amt))
	}
	oldBalance := acc.GetBalance()
	newBalance := oldBalance.Sub(amt)
	if newBalance.IsNegative() {
		return amt, types.ErrInsufficientBalance(types.DefaultCodespace, fmt.Sprintf("insufficient account funds for token [%s]; %s < %s", contractID, oldBalance, amt))
	}
	acc = acc.SetBalance(newBalance)
	err = k.UpdateAccount(ctx, acc)
	if err != nil {
		return amt, err
	}
	return newBalance, nil
}

func (k Keeper) AddBalance(ctx sdk.Context, contractID string, addr sdk.AccAddress, amt sdk.Int) (sdk.Int, sdk.Error) {
	return k.addBalance(ctx, contractID, addr, amt)
}

func (k Keeper) addBalance(ctx sdk.Context, contractID string, addr sdk.AccAddress, amt sdk.Int) (sdk.Int, sdk.Error) {
	acc, err := k.GetOrNewAccount(ctx, contractID, addr)
	if err != nil {
		return amt, err
	}
	oldBalance := acc.GetBalance()
	newBalance := oldBalance.Add(amt)
	if newBalance.IsNegative() {
		return amt, types.ErrInsufficientBalance(types.DefaultCodespace, fmt.Sprintf("insufficient account funds for token [%s]; %s < %s", contractID, oldBalance, amt))
	}
	acc = acc.SetBalance(newBalance)
	err = k.UpdateAccount(ctx, acc)
	if err != nil {
		return amt, err
	}
	return newBalance, nil
}
