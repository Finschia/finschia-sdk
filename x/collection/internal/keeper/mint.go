package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	linktype "github.com/line/link/types"
	"github.com/line/link/x/collection/internal/types"
)

type MintKeeper interface {
	MintCFT(ctx sdk.Context, from, to sdk.AccAddress, amount linktype.CoinWithTokenIDs) sdk.Error
	MintCNFT(ctx sdk.Context, from sdk.AccAddress, token types.NFT) sdk.Error
}

func (k Keeper) MintCFT(ctx sdk.Context, from, to sdk.AccAddress, amount linktype.CoinWithTokenIDs) sdk.Error {
	for _, coin := range amount {
		symbol, tokenID := coin.Symbol, coin.TokenID
		token, err := k.GetToken(ctx, symbol, tokenID)
		if err != nil {
			return err
		}
		if err := k.isMintable(ctx, token, from); err != nil {
			return err
		}
	}
	err := k.mintTokens(ctx, amount.ToCoins(), to)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMintCFT,
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTo, to.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.ToCoins().String()),
		),
	})
	return nil
}

func (k Keeper) MintCNFT(ctx sdk.Context, from sdk.AccAddress, token types.NFT) sdk.Error {
	if !types.ValidTokenURI(token.GetTokenURI()) {
		return types.ErrInvalidTokenURILength(types.DefaultCodespace, token.GetTokenURI())
	}
	if !k.hasTokenType(ctx, token.GetSymbol(), token.GetTokenType()) {
		return types.ErrCollectionTokenTypeNotExist(types.DefaultCodespace, token.GetSymbol(), token.GetTokenType())
	}

	perm := types.NewMintPermission(token.GetSymbol() + token.GetTokenType())
	if !k.HasPermission(ctx, from, perm) {
		return types.ErrTokenNoPermission(types.DefaultCodespace, from, perm)
	}

	err := k.SetToken(ctx, token)
	if err != nil {
		return err
	}

	err = k.mintTokens(ctx, sdk.NewCoins(sdk.NewCoin(token.GetDenom(), sdk.NewInt(1))), token.GetOwner())
	if err != nil {
		return err
	}

	tokenURIModifyPerm := types.NewModifyTokenURIPermission(token.GetDenom())
	k.AddPermission(ctx, token.GetOwner(), tokenURIModifyPerm)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMintCNFT,
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
func (k Keeper) isMintable(ctx sdk.Context, token types.Token, from sdk.AccAddress) sdk.Error {
	ft, ok := token.(types.FT)
	if !ok {
		return types.ErrTokenNotMintable(types.DefaultCodespace, token.GetSymbol(), token.GetTokenID())
	}

	if !ft.GetMintable() {
		return types.ErrTokenNotMintable(types.DefaultCodespace, token.GetSymbol(), token.GetTokenID())
	}
	perm := types.NewMintPermission(ft.GetDenom())
	if !k.HasPermission(ctx, from, perm) {
		return types.ErrTokenNoPermission(types.DefaultCodespace, from, perm)
	}
	return nil
}

func (k Keeper) mintTokens(ctx sdk.Context, amount sdk.Coins, to sdk.AccAddress) sdk.Error {
	err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, amount)
	if err != nil {
		return err
	}

	moduleAddr := k.supplyKeeper.GetModuleAddress(types.ModuleName)
	if moduleAddr == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("module account %s does not exist", types.ModuleName))
	}

	_, err = k.bankKeeper.SubtractCoins(ctx, moduleAddr, amount)
	if err != nil {
		return err
	}

	_, err = k.bankKeeper.AddCoins(ctx, to, amount)
	if err != nil {
		return err
	}

	return nil
}
