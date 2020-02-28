package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/types"
)

//For the Token module
type BankKeeper interface {
	GetCoins(ctx sdk.Context, contractID string, addr sdk.AccAddress) types.Coins
	HasCoins(ctx sdk.Context, contractID string, addr sdk.AccAddress, amt types.Coins) bool
	SendCoins(ctx sdk.Context, contractID string, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt types.Coins) sdk.Error
	SubtractCoins(ctx sdk.Context, contractID string, addr sdk.AccAddress, amt types.Coins) (types.Coins, sdk.Error)
	AddCoins(ctx sdk.Context, contractID string, addr sdk.AccAddress, amt types.Coins) (types.Coins, sdk.Error)
	SetCoins(ctx sdk.Context, contractID string, addr sdk.AccAddress, amt types.Coins) sdk.Error
}

var _ BankKeeper = (*Keeper)(nil)

func (k Keeper) GetCoins(ctx sdk.Context, contractID string, addr sdk.AccAddress) types.Coins {
	acc, err := k.GetAccount(ctx, contractID, addr)
	if err != nil {
		return types.NewCoins()
	}
	return acc.GetCoins()
}

func (k Keeper) HasCoins(ctx sdk.Context, contractID string, addr sdk.AccAddress, amt types.Coins) bool {
	return k.GetCoins(ctx, contractID, addr).IsAllGTE(amt)
}

func (k Keeper) SendCoins(ctx sdk.Context, contractID string, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt types.Coins) sdk.Error {
	if !amt.IsValid() {
		return types.ErrInvalidCoin(types.DefaultCodespace, "send amount must be positive")
	}

	_, err := k.SubtractCoins(ctx, contractID, fromAddr, amt)
	if err != nil {
		return err
	}

	_, err = k.AddCoins(ctx, contractID, toAddr, amt)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) SubtractCoins(ctx sdk.Context, contractID string, addr sdk.AccAddress, amt types.Coins) (types.Coins, sdk.Error) {
	if !amt.IsValid() {
		return nil, types.ErrInvalidCoin(types.DefaultCodespace, "amount must be positive")
	}

	acc, err := k.GetAccount(ctx, contractID, addr)
	if err != nil {
		return types.ZeroCoins(contractID), err
	}
	oldCoins := acc.GetCoins()

	newCoins, hasNeg := oldCoins.SafeSub(amt)
	if hasNeg {
		return amt, sdk.ErrInsufficientCoins(
			fmt.Sprintf("insufficient account funds; %s < %s", oldCoins, amt),
		)
	}

	err = k.SetCoins(ctx, contractID, addr, newCoins)

	return newCoins, err
}

func (k Keeper) AddCoins(ctx sdk.Context, contractID string, addr sdk.AccAddress, amt types.Coins) (types.Coins, sdk.Error) {
	if !amt.IsValid() {
		return nil, types.ErrInvalidCoin(types.DefaultCodespace, "amount must be positive")
	}

	oldCoins := k.GetCoins(ctx, contractID, addr)
	newCoins := oldCoins.Add(amt...)

	if newCoins.IsAnyNegative() {
		return amt, sdk.ErrInsufficientCoins(
			fmt.Sprintf("insufficient account funds; %s < %s", oldCoins, amt),
		)
	}

	err := k.SetCoins(ctx, contractID, addr, newCoins)
	return newCoins, err
}

func (k Keeper) SetCoins(ctx sdk.Context, contractID string, addr sdk.AccAddress, amt types.Coins) sdk.Error {
	if !amt.IsValid() {
		return sdk.ErrInvalidCoins(amt.String())
	}

	acc, err := k.GetOrNewAccount(ctx, contractID, addr)
	if err != nil {
		return err
	}

	acc = acc.SetCoins(amt)
	err = k.UpdateAccount(ctx, acc)
	if err != nil {
		return err
	}
	return nil
}
