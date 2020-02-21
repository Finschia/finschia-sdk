package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
)

func (k Keeper) IssueToken(ctx sdk.Context, token types.Token, amount sdk.Int, owner sdk.AccAddress) sdk.Error {
	if !types.ValidTokenURI(token.GetTokenURI()) {
		return types.ErrInvalidTokenURILength(types.DefaultCodespace, token.GetTokenURI())
	}
	err := k.SetToken(ctx, token)
	if err != nil {
		return err
	}

	err = k.mintTokens(ctx, sdk.NewCoins(sdk.NewCoin(token.GetSymbol(), amount)), owner)
	if err != nil {
		return err
	}

	modifyTokenURIPermission := types.NewModifyTokenURIPermission(token.GetSymbol())
	k.AddPermission(ctx, owner, modifyTokenURIPermission)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeIssueToken,
			sdk.NewAttribute(types.AttributeKeyName, token.GetName()),
			sdk.NewAttribute(types.AttributeKeySymbol, token.GetSymbol()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyMintable, strconv.FormatBool(token.GetMintable())),
			sdk.NewAttribute(types.AttributeKeyDecimals, token.GetDecimals().String()),
		),
		sdk.NewEvent(
			types.EventTypeGrantPermToken,
			sdk.NewAttribute(types.AttributeKeyTo, owner.String()),
			sdk.NewAttribute(types.AttributeKeyResource, modifyTokenURIPermission.GetResource()),
			sdk.NewAttribute(types.AttributeKeyAction, modifyTokenURIPermission.GetAction()),
		),
	})

	if token.GetMintable() {
		mintPerm := types.NewMintPermission(token.GetSymbol())
		k.AddPermission(ctx, owner, mintPerm)
		burnPerm := types.NewBurnPermission(token.GetSymbol())
		k.AddPermission(ctx, owner, burnPerm)
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeGrantPermToken,
				sdk.NewAttribute(types.AttributeKeyTo, owner.String()),
				sdk.NewAttribute(types.AttributeKeyResource, mintPerm.GetResource()),
				sdk.NewAttribute(types.AttributeKeyAction, mintPerm.GetAction()),
			),
			sdk.NewEvent(
				types.EventTypeGrantPermToken,
				sdk.NewAttribute(types.AttributeKeyTo, owner.String()),
				sdk.NewAttribute(types.AttributeKeyResource, burnPerm.GetResource()),
				sdk.NewAttribute(types.AttributeKeyAction, burnPerm.GetAction()),
			),
		})
	}

	return nil
}
