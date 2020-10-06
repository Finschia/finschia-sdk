package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link-modules/x/token/internal/types"
)

func (k Keeper) BurnToken(ctx sdk.Context, amount sdk.Int, from sdk.AccAddress) error {
	err := k.isBurnable(ctx, from, from, amount)
	if err != nil {
		return err
	}

	err = k.BurnSupply(ctx, from, amount)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnToken,
			sdk.NewAttribute(types.AttributeKeyContractID, k.getContractID(ctx)),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
		),
	})
	return nil
}

func (k Keeper) BurnTokenFrom(ctx sdk.Context, proxy sdk.AccAddress, from sdk.AccAddress, amount sdk.Int) error {
	if !k.IsApproved(ctx, proxy, from) {
		return sdkerrors.Wrapf(types.ErrTokenNotApproved, "Proxy: %s, Approver: %s, ContractID: %s", proxy.String(), from.String(), k.getContractID(ctx))
	}

	err := k.isBurnable(ctx, proxy, from, amount)
	if err != nil {
		return err
	}

	err = k.BurnSupply(ctx, from, amount)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnTokenFrom,
			sdk.NewAttribute(types.AttributeKeyContractID, k.getContractID(ctx)),
			sdk.NewAttribute(types.AttributeKeyProxy, proxy.String()),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})
	return nil
}

func (k Keeper) isBurnable(ctx sdk.Context, permissionOwner, tokenOwner sdk.AccAddress, amount sdk.Int) error {
	if !k.HasBalance(ctx, tokenOwner, amount) {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "%v has not enough coins for %v", tokenOwner.String(), amount)
	}
	if !amount.IsPositive() {
		return sdkerrors.Wrap(types.ErrInvalidAmount, amount.String())
	}

	perm := types.NewBurnPermission()
	if !k.HasPermission(ctx, permissionOwner, perm) {
		return sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", permissionOwner.String(), perm.String())
	}
	return nil
}
