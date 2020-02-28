package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
)

func (k Keeper) MintToken(ctx sdk.Context, contractID string, amount sdk.Int, from, to sdk.AccAddress) sdk.Error {
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

func (k Keeper) isMintable(ctx sdk.Context, token types.Token, from sdk.AccAddress) sdk.Error {
	if !token.GetMintable() {
		return types.ErrTokenNotMintable(types.DefaultCodespace, token.GetContractID())
	}
	perm := types.NewMintPermission(token.GetContractID())
	if !k.HasPermission(ctx, from, perm) {
		return types.ErrTokenNoPermission(types.DefaultCodespace, from, perm)
	}
	return nil
}
