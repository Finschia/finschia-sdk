package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link-modules/x/coin"
	"github.com/line/link-modules/x/token/internal/types"
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
