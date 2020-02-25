package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	linktype "github.com/line/link/types"
	"github.com/line/link/x/collection/internal/types"
)

type MintKeeper interface {
	MintFT(ctx sdk.Context, symbol string, from, to sdk.AccAddress, amount linktype.CoinWithTokenIDs) sdk.Error
	MintNFT(ctx sdk.Context, symbol string, from sdk.AccAddress, token types.NFT) sdk.Error
}

func (k Keeper) MintFT(ctx sdk.Context, symbol string, from, to sdk.AccAddress, amount types.Coins) sdk.Error {
	for _, coin := range amount {
		token, err := k.GetToken(ctx, symbol, coin.Denom)
		if err != nil {
			return err
		}
		if err := k.isMintable(ctx, symbol, token, from); err != nil {
			return err
		}
	}
	err := k.MintSupply(ctx, symbol, to, amount)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMintFT,
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTo, to.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})
	return nil
}

func (k Keeper) MintNFT(ctx sdk.Context, symbol string, from sdk.AccAddress, token types.NFT) sdk.Error {
	if !types.ValidTokenURI(token.GetTokenURI()) {
		return types.ErrInvalidTokenURILength(types.DefaultCodespace, token.GetTokenURI())
	}
	if !k.hasTokenType(ctx, symbol, token.GetTokenType()) {
		return types.ErrCollectionTokenTypeNotExist(types.DefaultCodespace, token.GetSymbol(), token.GetTokenType())
	}

	perm := types.NewMintPermission(token.GetSymbol(), token.GetTokenType())
	if !k.HasPermission(ctx, from, perm) {
		return types.ErrTokenNoPermission(types.DefaultCodespace, from, perm)
	}

	err := k.SetToken(ctx, token)
	if err != nil {
		return err
	}

	err = k.MintSupply(ctx, symbol, token.GetOwner(), types.OneCoins(token.GetTokenID()))
	if err != nil {
		return err
	}

	tokenURIModifyPerm := types.NewModifyTokenURIPermission(symbol, token.GetTokenID())
	k.AddPermission(ctx, token.GetOwner(), tokenURIModifyPerm)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMintNFT,
			sdk.NewAttribute(types.AttributeKeyName, token.GetName()),
			sdk.NewAttribute(types.AttributeKeySymbol, token.GetSymbol()),
			sdk.NewAttribute(types.AttributeKeyTokenID, token.GetTokenID()),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTo, token.GetOwner().String()),
			sdk.NewAttribute(types.AttributeKeyTokenURI, token.GetTokenURI()),
		),
		sdk.NewEvent(
			types.EventTypeGrantPermToken,
			sdk.NewAttribute(types.AttributeKeyTo, token.GetOwner().String()),
			sdk.NewAttribute(types.AttributeKeyResource, tokenURIModifyPerm.GetResource()),
			sdk.NewAttribute(types.AttributeKeyAction, tokenURIModifyPerm.GetAction()),
		),
	})

	return nil
}
func (k Keeper) isMintable(ctx sdk.Context, symbol string, token types.Token, from sdk.AccAddress) sdk.Error {
	ft, ok := token.(types.FT)
	if !ok {
		return types.ErrTokenNotMintable(types.DefaultCodespace, symbol, token.GetTokenID())
	}

	if !ft.GetMintable() {
		return types.ErrTokenNotMintable(types.DefaultCodespace, symbol, token.GetTokenID())
	}
	perm := types.NewMintPermission(symbol, ft.GetTokenID())
	if !k.HasPermission(ctx, from, perm) {
		return types.ErrTokenNoPermission(types.DefaultCodespace, from, perm)
	}
	return nil
}
