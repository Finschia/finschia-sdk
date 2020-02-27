package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/types"
)

type IssueKeeper interface {
	IssueFT(ctx sdk.Context, symbol string, owner sdk.AccAddress, token types.FT, amount sdk.Int) sdk.Error
	IssueNFT(ctx sdk.Context, symbol string, owner sdk.AccAddress, tokenType string) sdk.Error
}

func (k Keeper) IssueFT(ctx sdk.Context, symbol string, owner sdk.AccAddress, token types.FT, amount sdk.Int) sdk.Error {
	if !types.ValidTokenURI(token.GetTokenURI()) {
		return types.ErrInvalidTokenURILength(types.DefaultCodespace, token.GetTokenURI())
	}
	err := k.SetToken(ctx, symbol, token)
	if err != nil {
		return err
	}

	err = k.MintSupply(ctx, symbol, owner, types.NewCoins(types.NewCoin(token.GetTokenID(), amount)))
	if err != nil {
		return err
	}

	tokenURIModifyPerm := types.NewModifyTokenURIPermission(symbol)
	k.AddPermission(ctx, owner, tokenURIModifyPerm)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeIssueFT,
			sdk.NewAttribute(types.AttributeKeyName, token.GetName()),
			sdk.NewAttribute(types.AttributeKeySymbol, token.GetSymbol()),
			sdk.NewAttribute(types.AttributeKeyTokenID, token.GetTokenID()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyMintable, strconv.FormatBool(token.GetMintable())),
			sdk.NewAttribute(types.AttributeKeyDecimals, token.GetDecimals().String()),
			sdk.NewAttribute(types.AttributeKeyTokenURI, token.GetTokenURI()),
		),
	})

	return nil
}

func (k Keeper) IssueNFT(ctx sdk.Context, symbol string, tokenType types.TokenType, owner sdk.AccAddress) sdk.Error {
	err := k.SetTokenType(ctx, symbol, tokenType)
	if err != nil {
		return err
	}

	mintPerm := types.NewMintPermission(symbol)
	k.AddPermission(ctx, owner, mintPerm)
	burnPerm := types.NewBurnPermission(symbol)
	k.AddPermission(ctx, owner, burnPerm)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeIssueNFT,
			sdk.NewAttribute(types.AttributeKeySymbol, symbol),
			sdk.NewAttribute(types.AttributeKeyTokenType, tokenType.GetTokenType()),
		),
	})

	return nil
}
