package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link/x/collection/internal/types"
)

type BurnKeeper interface {
	BurnFT(ctx sdk.Context, contractID string, from sdk.AccAddress, amount types.Coins) error
	BurnNFT(ctx sdk.Context, contractID string, from sdk.AccAddress, tokenIDs ...string) error
	BurnFTFrom(ctx sdk.Context, contractID string, proxy sdk.AccAddress, from sdk.AccAddress, amount types.Coins) error
	BurnNFTFrom(ctx sdk.Context, contractID string, proxy sdk.AccAddress, from sdk.AccAddress, tokenIDs ...string) error
}

var _ BurnKeeper = (*Keeper)(nil)

func (k Keeper) BurnFT(ctx sdk.Context, contractID string, from sdk.AccAddress, amount types.Coins) error {
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

func (k Keeper) BurnFTFrom(ctx sdk.Context, contractID string, proxy sdk.AccAddress, from sdk.AccAddress, amount types.Coins) error {
	if !k.IsApproved(ctx, contractID, proxy, from) {
		return sdkerrors.Wrapf(types.ErrCollectionNotApproved, "Proxy: %s, Approver: %s, ContractID: %s", proxy.String(), from.String(), contractID)
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

func (k Keeper) burnFT(ctx sdk.Context, contractID string, permissionOwner, tokenOwner sdk.AccAddress, amount types.Coins) error {
	if err := k.isBurnable(ctx, contractID, permissionOwner, tokenOwner, amount); err != nil {
		return err
	}

	if err := k.BurnSupply(ctx, contractID, tokenOwner, amount); err != nil {
		return err
	}
	return nil
}

func (k Keeper) isBurnable(ctx sdk.Context, contractID string, permissionOwner, tokenOwner sdk.AccAddress, amount types.Coins) error {
	perm := types.NewBurnPermission()
	if !k.HasPermission(ctx, contractID, permissionOwner, perm) {
		return sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", permissionOwner.String(), perm.String())
	}

	if !k.HasCoins(ctx, contractID, tokenOwner, amount) {
		return sdkerrors.Wrapf(types.ErrInsufficientToken, "%v has not enough coins for %v", tokenOwner.String(), amount)
	}
	return nil
}

func (k Keeper) BurnNFT(ctx sdk.Context, contractID string, from sdk.AccAddress, tokenIDs ...string) error {
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

func (k Keeper) BurnNFTFrom(ctx sdk.Context, contractID string, proxy sdk.AccAddress, from sdk.AccAddress, tokenIDs ...string) error {
	if !k.IsApproved(ctx, contractID, proxy, from) {
		return sdkerrors.Wrapf(types.ErrCollectionNotApproved, "Proxy: %s, Approver: %s, ContractID: %s", proxy.String(), from.String(), contractID)
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

func (k Keeper) burnNFT(ctx sdk.Context, contractID string, permissionOwner, tokenOwner sdk.AccAddress, tokenID string) error {
	token, err := k.GetNFT(ctx, contractID, tokenID)
	if err != nil {
		return err
	}

	perm := types.NewBurnPermission()
	if !k.HasPermission(ctx, contractID, permissionOwner, perm) {
		return sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", permissionOwner.String(), perm.String())
	}

	if !token.GetOwner().Equals(tokenOwner) {
		return sdkerrors.Wrapf(types.ErrTokenNotOwnedBy, "TokenID: %s, Owner: %s", tokenID, tokenOwner.String())
	}

	err = k.burnNFTRecursive(ctx, contractID, token, tokenOwner)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) burnNFTRecursive(ctx sdk.Context, contractID string, token types.NFT, from sdk.AccAddress) error {
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
