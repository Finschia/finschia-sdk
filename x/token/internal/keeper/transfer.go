package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/bank"
	"github.com/line/link/x/token/internal/types"
)

func (k Keeper) Transfer(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, symbol string, amount sdk.Int) sdk.Error {
	coin := sdk.NewCoins(sdk.NewCoin(symbol, amount))

	// reject if to address is blacklisted (safety box addresses)
	if k.IsBlacklisted(ctx, to, bank.ActionTransferTo) {
		return bank.ErrCanNotTransferToBlacklisted(types.DefaultCodespace, to.String())
	}

	_, err := k.bankKeeper.SubtractCoins(ctx, from, coin)
	if err != nil {
		return err
	}

	_, err = k.bankKeeper.AddCoins(ctx, to, coin)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransfer,
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTo, to.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, symbol),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})

	return nil
}
