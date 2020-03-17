package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/types"
)

type BurnKeeper interface {
	BurnFT(ctx sdk.Context, contractID string, from sdk.AccAddress, amount types.Coins) sdk.Error
	BurnNFT(ctx sdk.Context, contractID string, from sdk.AccAddress, tokenIDs ...string) sdk.Error
	BurnFTFrom(ctx sdk.Context, contractID string, proxy sdk.AccAddress, from sdk.AccAddress, amount types.Coins) sdk.Error
	BurnNFTFrom(ctx sdk.Context, contractID string, proxy sdk.AccAddress, from sdk.AccAddress, tokenIDs ...string) sdk.Error
}

var _ BurnKeeper = (*Keeper)(nil)

func (k Keeper) BurnFT(ctx sdk.Context, contractID string, from sdk.AccAddress, amount types.Coins) sdk.Error {
	if err := k.burnFT(ctx, contractID, from, from, amount); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnFT,
			sdk.NewAttribute(types.AttributeKeyContractID, contractID),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})

	return nil
}

func (k Keeper) BurnFTFrom(ctx sdk.Context, contractID string, proxy sdk.AccAddress, from sdk.AccAddress, amount types.Coins) sdk.Error {
	if !k.IsApproved(ctx, contractID, proxy, from) {
		return types.ErrCollectionNotApproved(types.DefaultCodespace, proxy.String(), from.String(), contractID)
	}

	if err := k.burnFT(ctx, contractID, proxy, from, amount); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnFTFrom,
			sdk.NewAttribute(types.AttributeKeyContractID, contractID),
			sdk.NewAttribute(types.AttributeKeyProxy, proxy.String()),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})
	return nil
}

func (k Keeper) burnFT(ctx sdk.Context, contractID string, permissionOwner, tokenOwner sdk.AccAddress, amount types.Coins) sdk.Error {
	if err := k.isBurnable(ctx, contractID, permissionOwner, tokenOwner, amount); err != nil {
		return err
	}

	if err := k.BurnSupply(ctx, contractID, tokenOwner, amount); err != nil {
		return err
	}
	return nil
}

func (k Keeper) isBurnable(ctx sdk.Context, contractID string, permissionOwner, tokenOwner sdk.AccAddress, amount types.Coins) sdk.Error {
	perm := types.NewBurnPermission(contractID)
	if !k.HasPermission(ctx, permissionOwner, perm) {
		return types.ErrTokenNoPermission(types.DefaultCodespace, permissionOwner, perm)
	}

	if !k.HasCoins(ctx, contractID, tokenOwner, amount) {
		return types.ErrInsufficientToken(types.DefaultCodespace, fmt.Sprintf("%v has not enough coins for %v", tokenOwner.String(), amount))
	}

	return nil
}

func (k Keeper) BurnNFT(ctx sdk.Context, contractID string, from sdk.AccAddress, tokenIDs ...string) sdk.Error {
	for _, tokenID := range tokenIDs {
		if err := k.burnNFT(ctx, contractID, from, from, tokenID); err != nil {
			return err
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnNFT,
			sdk.NewAttribute(types.AttributeKeyContractID, contractID),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
		),
	})

	for _, tokenID := range tokenIDs {
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeBurnNFT,
				sdk.NewAttribute(types.AttributeKeyTokenID, tokenID),
			),
		})
	}

	return nil
}

func (k Keeper) BurnNFTFrom(ctx sdk.Context, contractID string, proxy sdk.AccAddress, from sdk.AccAddress, tokenIDs ...string) sdk.Error {
	if !k.IsApproved(ctx, contractID, proxy, from) {
		return types.ErrCollectionNotApproved(types.DefaultCodespace, proxy.String(), from.String(), contractID)
	}

	for _, tokenID := range tokenIDs {
		if err := k.burnNFT(ctx, contractID, proxy, from, tokenID); err != nil {
			return err
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnNFTFrom,
			sdk.NewAttribute(types.AttributeKeyContractID, contractID),
			sdk.NewAttribute(types.AttributeKeyProxy, proxy.String()),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
		),
	})
	for _, tokenID := range tokenIDs {
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeBurnNFTFrom,
				sdk.NewAttribute(types.AttributeKeyTokenID, tokenID),
			),
		})
	}
	return nil
}

func (k Keeper) burnNFT(ctx sdk.Context, contractID string, permissionOwner, tokenOwner sdk.AccAddress, tokenID string) sdk.Error {
	token, err := k.GetNFT(ctx, contractID, tokenID)
	if err != nil {
		return err
	}

	perm := types.NewBurnPermission(contractID)
	if !k.HasPermission(ctx, permissionOwner, perm) {
		return types.ErrTokenNoPermission(types.DefaultCodespace, permissionOwner, perm)
	}

	if !token.GetOwner().Equals(tokenOwner) {
		return types.ErrTokenNotOwnedBy(types.DefaultCodespace, tokenID, tokenOwner)
	}

	err = k.burnNFTRecursive(ctx, contractID, token, tokenOwner)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) burnNFTRecursive(ctx sdk.Context, contractID string, token types.NFT, from sdk.AccAddress) sdk.Error {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeOperationBurnNFT,
			sdk.NewAttribute(types.AttributeKeyTokenID, token.GetTokenID()),
		),
	})

	children, err := k.ChildrenOf(ctx, contractID, token.GetTokenID())
	if err != nil {
		return err
	}

	for _, child := range children {
		err = k.burnNFTRecursive(ctx, contractID, child.(types.NFT), from)
		if err != nil {
			return err
		}
	}

	parent, err := k.ParentOf(ctx, contractID, token.GetTokenID())
	if err != nil {
		return err
	}
	if parent != nil {
		_, err = k.detach(ctx, contractID, from, token.GetTokenID())
		if err != nil {
			return err
		}
	}
	err = k.DeleteToken(ctx, contractID, token.GetTokenID())
	if err != nil {
		return err
	}
	err = k.BurnSupply(ctx, contractID, token.GetOwner(), types.OneCoins(token.GetTokenID()))
	if err != nil {
		return err
	}
	return nil
}
