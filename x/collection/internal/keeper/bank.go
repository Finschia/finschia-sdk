package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/collection/internal/types"
)

// For the Token module
type BankKeeper interface {
	GetCoins(ctx sdk.Context, addr sdk.AccAddress) types.Coins
	HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt types.Coins) bool
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt types.Coins) error
	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt types.Coins) (types.Coins, error)
	AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt types.Coins) (types.Coins, error)
	SetCoins(ctx sdk.Context, addr sdk.AccAddress, amt types.Coins) error
}

var _ BankKeeper = (*Keeper)(nil)

func (k Keeper) GetCoins(ctx sdk.Context, addr sdk.AccAddress) types.Coins {
	acc, err := k.GetAccount(ctx, addr)
	if err != nil {
		return types.NewCoins()
	}
	return acc.GetCoins()
}

func (k Keeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt types.Coins) bool {
	return k.GetCoins(ctx, addr).IsAllGTE(amt)
}

func (k Keeper) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt types.Coins) error {
	if !amt.IsValid() {
		return sdkerrors.Wrap(types.ErrInvalidCoin, "send amount must be positive")
	}

	_, err := k.SubtractCoins(ctx, fromAddr, amt)
	if err != nil {
		return err
	}

	_, err = k.AddCoins(ctx, toAddr, amt)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt types.Coins) (types.Coins, error) {
	if !amt.IsValid() {
		return nil, sdkerrors.Wrap(types.ErrInvalidCoin, "amount must be positive")
	}

	acc, err := k.GetAccount(ctx, addr)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInsufficientToken, "insufficient account funds[%s]; account has no coin", k.getContractID(ctx))
	}
	oldCoins := acc.GetCoins()

	newCoins, hasNeg := oldCoins.SafeSub(amt)
	if hasNeg {
		return amt, sdkerrors.Wrapf(types.ErrInsufficientToken, "insufficient account funds[%s]; %s < %s", k.getContractID(ctx), oldCoins, amt)
	}

	err = k.SetCoins(ctx, addr, newCoins)

	return newCoins, err
}

func (k Keeper) AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt types.Coins) (types.Coins, error) {
	if !amt.IsValid() {
		return nil, sdkerrors.Wrap(types.ErrInvalidCoin, "amount must be positive")
	}

	oldCoins := k.GetCoins(ctx, addr)
	newCoins := oldCoins.Add(amt...)

	err := k.SetCoins(ctx, addr, newCoins)
	return newCoins, err
}

func (k Keeper) SetCoins(ctx sdk.Context, addr sdk.AccAddress, amt types.Coins) error {
	if !amt.IsValid() {
		return sdkerrors.Wrapf(types.ErrInvalidCoin, "invalid amount: %s", amt.String())
	}

	acc, err := k.GetOrNewAccount(ctx, addr)
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
