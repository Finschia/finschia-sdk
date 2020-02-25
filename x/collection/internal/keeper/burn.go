package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	linktype "github.com/line/link/types"
	"github.com/line/link/x/collection/internal/types"
)

type BurnKeeper interface {
	BurnFT(ctx sdk.Context, from sdk.AccAddress, amount linktype.CoinWithTokenIDs) sdk.Error
	BurnNFT(ctx sdk.Context, from sdk.AccAddress, symbol, tokenID string) sdk.Error
	BurnFTFrom(ctx sdk.Context, proxy sdk.AccAddress, from sdk.AccAddress, amount linktype.CoinWithTokenIDs) sdk.Error
	BurnNFTFrom(ctx sdk.Context, proxy sdk.AccAddress, from sdk.AccAddress, symbol, tokenID string) sdk.Error
}

var _ BurnKeeper = (*Keeper)(nil)

func (k Keeper) BurnFT(ctx sdk.Context, from sdk.AccAddress, amount linktype.CoinWithTokenIDs) sdk.Error {
	if err := k.burnFT(ctx, from, from, amount); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnFT,
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.ToCoins().String()),
		),
	})

	return nil
}

func (k Keeper) BurnFTFrom(ctx sdk.Context, proxy sdk.AccAddress, from sdk.AccAddress, amount linktype.CoinWithTokenIDs) sdk.Error {
	for _, coin := range amount {
		if !k.IsApproved(ctx, proxy, from, coin.Symbol) {
			return types.ErrCollectionNotApproved(types.DefaultCodespace, proxy.String(), from.String(), coin.Symbol)
		}
	}

	if err := k.burnFT(ctx, proxy, from, amount); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnFTFrom,
			sdk.NewAttribute(types.AttributeKeyProxy, proxy.String()),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.ToCoins().String()),
		),
	})
	return nil
}

func (k Keeper) burnFT(ctx sdk.Context, permissionOwner, tokenOwner sdk.AccAddress, amount linktype.CoinWithTokenIDs) sdk.Error {
	coins := amount.ToCoins()

	if err := k.isBurnable(ctx, permissionOwner, tokenOwner, coins); err != nil {
		return err
	}

	if err := k.burnTokens(ctx, tokenOwner, coins); err != nil {
		return err
	}
	return nil
}

func (k Keeper) isBurnable(ctx sdk.Context, permissionOwner, tokenOwner sdk.AccAddress, amount sdk.Coins) sdk.Error {
	if !k.hasEnoughCoins(ctx, tokenOwner, amount) {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("%v has not enough coins for %v", tokenOwner, amount))
	}

	for _, coin := range amount {
		perm := types.NewBurnPermission(coin.Denom)
		if !k.HasPermission(ctx, permissionOwner, perm) {
			return types.ErrTokenNoPermission(types.DefaultCodespace, permissionOwner, perm)
		}
	}
	return nil
}

func (k Keeper) hasEnoughCoins(ctx sdk.Context, from sdk.AccAddress, amount sdk.Coins) bool {
	return k.bankKeeper.GetCoins(ctx, from).IsAllGTE(amount)
}

func (k Keeper) BurnNFT(ctx sdk.Context, from sdk.AccAddress, symbol, tokenID string) sdk.Error {
	if err := k.burnNFT(ctx, from, from, symbol, tokenID); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnNFT,
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, symbol),
			sdk.NewAttribute(types.AttributeKeyTokenID, tokenID),
		),
	})
	return nil
}

func (k Keeper) BurnNFTFrom(ctx sdk.Context, proxy sdk.AccAddress, from sdk.AccAddress, symbol, tokenID string) sdk.Error {
	if !k.IsApproved(ctx, proxy, from, symbol) {
		return types.ErrCollectionNotApproved(types.DefaultCodespace, proxy.String(), from.String(), symbol)
	}

	if err := k.burnNFT(ctx, proxy, from, symbol, tokenID); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnNFTFrom,
			sdk.NewAttribute(types.AttributeKeyProxy, proxy.String()),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, symbol),
			sdk.NewAttribute(types.AttributeKeyTokenID, tokenID),
		),
	})
	return nil
}

func (k Keeper) burnNFT(ctx sdk.Context, permissionOwner, tokenOwner sdk.AccAddress, symbol, tokenID string) sdk.Error {
	token, err := k.GetNFT(ctx, symbol, tokenID)
	if err != nil {
		return err
	}

	perm := types.NewBurnPermission(symbol + tokenID[:types.TokenTypeLength])
	if !k.HasPermission(ctx, permissionOwner, perm) {
		return types.ErrTokenNoPermission(types.DefaultCodespace, permissionOwner, perm)
	}

	if !token.GetOwner().Equals(tokenOwner) {
		return types.ErrTokenNotOwnedBy(types.DefaultCodespace, symbol+tokenID, tokenOwner)
	}

	err = k.burnNFTrecursive(ctx, token, tokenOwner)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) burnNFTrecursive(ctx sdk.Context, token types.NFT, from sdk.AccAddress) sdk.Error {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeOperationBurnNFT,
			sdk.NewAttribute(types.AttributeKeyTokenID, token.GetTokenID()),
		),
	})

	children, err := k.ChildrenOf(ctx, token.GetSymbol(), token.GetTokenID())
	if err != nil {
		return err
	}

	for _, child := range children {
		err = k.burnNFTrecursive(ctx, child.(types.NFT), from)
		if err != nil {
			return err
		}
	}

	parent, err := k.ParentOf(ctx, token.GetSymbol(), token.GetTokenID())
	if err != nil {
		return err
	}
	if parent != nil {
		_, err = k.detach(ctx, from, token.GetSymbol(), token.GetTokenID())
		if err != nil {
			return err
		}
	}
	collection, err := k.GetCollection(ctx, token.GetSymbol())
	if err != nil {
		return err
	}
	collection, err = collection.DeleteToken(token)
	if err != nil {
		return err
	}
	err = k.UpdateCollection(ctx, collection)
	if err != nil {
		return err
	}
	err = k.burnTokens(ctx, token.GetOwner(), sdk.NewCoins(sdk.NewCoin(token.GetDenom(), sdk.NewInt(1))))
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) burnTokens(ctx sdk.Context, from sdk.AccAddress, amount sdk.Coins) sdk.Error {
	moduleAddr := k.supplyKeeper.GetModuleAddress(types.ModuleName)
	if moduleAddr == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("module account %s does not exist", types.ModuleName))
	}

	_, err := k.bankKeeper.SubtractCoins(ctx, from, amount)
	if err != nil {
		return err
	}

	_, err = k.bankKeeper.AddCoins(ctx, moduleAddr, amount)
	if err != nil {
		return err
	}

	err = k.supplyKeeper.BurnCoins(ctx, types.ModuleName, amount)
	if err != nil {
		return err
	}
	return nil
}
