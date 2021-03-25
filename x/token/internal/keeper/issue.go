package keeper

import (
	"strconv"
	"unicode/utf8"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link-modules/x/token/internal/types"
)

func (k Keeper) IssueToken(ctx sdk.Context, token types.Token, amount sdk.Int, owner, to sdk.AccAddress) error {
	if !types.ValidateImageURI(token.GetImageURI()) {
		return sdkerrors.Wrapf(types.ErrInvalidImageURILength, "[%s] should be shorter than [%d] UTF-8 characters, current length: [%d]", token.GetImageURI(), types.MaxImageURILength, utf8.RuneCountInString(token.GetImageURI()))
	}
	err := k.SetToken(ctx, token)
	if err != nil {
		return err
	}

	err = k.MintSupply(ctx, to, amount)
	if err != nil {
		return err
	}

	modifyTokenURIPermission := types.NewModifyPermission()
	k.AddPermission(ctx, owner, modifyTokenURIPermission)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeIssueToken,
			sdk.NewAttribute(types.AttributeKeyContractID, token.GetContractID()),
			sdk.NewAttribute(types.AttributeKeyName, token.GetName()),
			sdk.NewAttribute(types.AttributeKeySymbol, token.GetSymbol()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyTo, to.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyMintable, strconv.FormatBool(token.GetMintable())),
			sdk.NewAttribute(types.AttributeKeyDecimals, token.GetDecimals().String()),
		),
		sdk.NewEvent(
			types.EventTypeGrantPermToken,
			sdk.NewAttribute(types.AttributeKeyTo, owner.String()),
			sdk.NewAttribute(types.AttributeKeyContractID, token.GetContractID()),
			sdk.NewAttribute(types.AttributeKeyPerm, modifyTokenURIPermission.String()),
		),
	})

	if token.GetMintable() {
		mintPerm := types.NewMintPermission()
		k.AddPermission(ctx, owner, mintPerm)
		burnPerm := types.NewBurnPermission()
		k.AddPermission(ctx, owner, burnPerm)
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeGrantPermToken,
				sdk.NewAttribute(types.AttributeKeyTo, owner.String()),
				sdk.NewAttribute(types.AttributeKeyContractID, token.GetContractID()),
				sdk.NewAttribute(types.AttributeKeyPerm, mintPerm.String()),
			),
			sdk.NewEvent(
				types.EventTypeGrantPermToken,
				sdk.NewAttribute(types.AttributeKeyTo, owner.String()),
				sdk.NewAttribute(types.AttributeKeyContractID, token.GetContractID()),
				sdk.NewAttribute(types.AttributeKeyPerm, burnPerm.String()),
			),
		})
	}

	return nil
}
