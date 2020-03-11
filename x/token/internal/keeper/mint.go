package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link/x/token/internal/types"
)

func (k Keeper) MintToken(ctx sdk.Context, contractID string, amount sdk.Int, from, to sdk.AccAddress) error {
	token, err := k.GetToken(ctx, contractID)
	if err != nil {
		return err
	}
	if err := k.isMintable(ctx, token, from); err != nil {
		return err
	}
	err = k.MintSupply(ctx, contractID, to, amount)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMintToken,
			sdk.NewAttribute(types.AttributeKeyContractID, contractID),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTo, to.String()),
		),
	})
	return nil
}

func (k Keeper) isMintable(ctx sdk.Context, token types.Token, from sdk.AccAddress) error {
	if !token.GetMintable() {
		return sdkerrors.Wrapf(types.ErrTokenNotMintable, "ContractID: %s", token.GetContractID())
	}
	perm := types.NewMintPermission(token.GetContractID())
	if !k.HasPermission(ctx, from, perm) {
		return sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", from.String(), perm.String())
	}
	return nil
}
