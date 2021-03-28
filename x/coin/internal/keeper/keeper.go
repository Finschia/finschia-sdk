package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/lbm-sdk/v2/x/coin/internal/types"
)

type Keeper struct {
	bk       types.BankKeeper
	storeKey sdk.StoreKey
}

func NewKeeper(bk types.BankKeeper, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		bk:       bk,
		storeKey: storeKey,
	}
}

// SendCoins moves coins from one account to another
func (keeper Keeper) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	// reject if to address is blacklisted (safety box addresses)
	if keeper.IsBlacklistedAccountAction(ctx, toAddr, types.ActionTransferTo) {
		return sdkerrors.Wrapf(types.ErrCanNotTransferToBlacklisted, "Addr: %s", toAddr.String())
	}

	_, err := keeper.bk.SubtractCoins(ctx, fromAddr, amt)
	if err != nil {
		return err
	}

	_, err = keeper.bk.AddCoins(ctx, toAddr, amt)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransfer,
			sdk.NewAttribute(types.AttributeKeySender, fromAddr.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, toAddr.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, amt.String()),
		),
	})

	return nil
}

// InputOutputCoins handles a list of inputs and outputs
func (keeper Keeper) InputOutputCoins(ctx sdk.Context, inputs []types.Input, outputs []types.Output) error {
	if err := types.ValidateInputsOutputs(inputs, outputs); err != nil {
		return err
	}

	for _, in := range inputs {
		_, err := keeper.bk.SubtractCoins(ctx, in.Address, in.Coins)
		if err != nil {
			return err
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeTransfer,
				sdk.NewAttribute(types.AttributeKeySender, in.Address.String()),
			),
		)
	}

	for _, out := range outputs {
		// reject if to address is blacklisted (safety box addresses)
		if keeper.IsBlacklistedAccountAction(ctx, out.Address, types.ActionTransferTo) {
			return sdkerrors.Wrapf(types.ErrCanNotTransferToBlacklisted, "Addr: %s", out.Address.String())
		}

		_, err := keeper.bk.AddCoins(ctx, out.Address, out.Coins)
		if err != nil {
			return err
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeTransfer,
				sdk.NewAttribute(types.AttributeKeyRecipient, out.Address.String()),
			),
		)
	}

	return nil
}

// GetCoins returns the coins at the addr.
func (keeper Keeper) GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return keeper.bk.GetCoins(ctx, addr)
}

// HasCoins returns whether or not an account has at least amt coins.
func (keeper Keeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool {
	return keeper.GetCoins(ctx, addr).IsAllGTE(amt)
}

func (keeper Keeper) BlacklistedAddr(creator sdk.AccAddress) bool {
	return keeper.bk.BlacklistedAddr(creator)
}
