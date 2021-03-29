package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/lbm-sdk/v2/x/token/internal/types"
)

func (k Keeper) MintToken(ctx sdk.Context, amount sdk.Int, from, to sdk.AccAddress) error {
	token, err := k.GetToken(ctx)
	if err != nil {
		return err
	}
	if err := k.isMintable(ctx, token, from, amount); err != nil {
		return err
	}
	err = k.MintSupply(ctx, to, amount)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMintToken,
			sdk.NewAttribute(types.AttributeKeyContractID, k.getContractID(ctx)),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTo, to.String()),
		),
	})
	return nil
}

func (k Keeper) isMintable(ctx sdk.Context, token types.Token, from sdk.AccAddress, amount sdk.Int) error {
	if !token.GetMintable() {
		return sdkerrors.Wrapf(types.ErrTokenNotMintable, "ContractID: %s", token.GetContractID())
	}
	if !amount.IsPositive() {
		return sdkerrors.Wrap(types.ErrInvalidAmount, amount.String())
	}
	perm := types.NewMintPermission()
	if !k.HasPermission(ctx, from, perm) {
		return sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", from.String(), perm.String())
	}
	return nil
}
