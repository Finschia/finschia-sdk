package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/coin"
	"github.com/line/lbm-sdk/x/token/internal/types"
)

func (k Keeper) Transfer(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amount sdk.Int) error {
	// reject if to address is blacklisted (safety box addresses)
	if k.IsBlacklisted(ctx, to, coin.ActionTransferTo) {
		return sdkerrors.Wrapf(coin.ErrCanNotTransferToBlacklisted, "Addr: %s", to.String())
	}

	err := k.Send(ctx, from, to, amount)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransfer,
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTo, to.String()),
			sdk.NewAttribute(types.AttributeKeyContractID, k.getContractID(ctx)),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})

	return nil
}

func (k Keeper) TransferFrom(ctx sdk.Context, proxy sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, amount sdk.Int) error {
	if !k.IsApproved(ctx, proxy, from) {
		return sdkerrors.Wrapf(types.ErrTokenNotApproved, "Proxy: %s, Approver: %s, ContractID: %s", proxy.String(), from.String(), k.getContractID(ctx))
	}

	if err := k.Send(ctx, from, to, amount); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransferFrom,
			sdk.NewAttribute(types.AttributeKeyContractID, k.getContractID(ctx)),
			sdk.NewAttribute(types.AttributeKeyProxy, proxy.String()),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTo, to.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})

	return nil
}
