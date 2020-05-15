package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link/x/collection/internal/types"
)

type BurnKeeper interface {
	BurnFT(ctx sdk.Context, from sdk.AccAddress, amount types.Coins) error
	BurnNFT(ctx sdk.Context, from sdk.AccAddress, tokenIDs ...string) error
	BurnFTFrom(ctx sdk.Context, proxy sdk.AccAddress, from sdk.AccAddress, amount types.Coins) error
	BurnNFTFrom(ctx sdk.Context, proxy sdk.AccAddress, from sdk.AccAddress, tokenIDs ...string) error
}

var _ BurnKeeper = (*Keeper)(nil)

func (k Keeper) BurnFT(ctx sdk.Context, from sdk.AccAddress, amount types.Coins) error {
	if err := k.burnFT(ctx, from, from, amount); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnFT,
			sdk.NewAttribute(types.AttributeKeyContractID, k.getContractID(ctx)),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})

	return nil
}

func (k Keeper) BurnFTFrom(ctx sdk.Context, proxy sdk.AccAddress, from sdk.AccAddress, amount types.Coins) error {
	if !k.IsApproved(ctx, proxy, from) {
		return sdkerrors.Wrapf(types.ErrCollectionNotApproved, "Proxy: %s, Approver: %s, ContractID: %s", proxy.String(), from.String(), k.getContractID(ctx))
	}

	if err := k.burnFT(ctx, proxy, from, amount); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnFTFrom,
			sdk.NewAttribute(types.AttributeKeyContractID, k.getContractID(ctx)),
			sdk.NewAttribute(types.AttributeKeyProxy, proxy.String()),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})
	return nil
}

func (k Keeper) burnFT(ctx sdk.Context, permissionOwner, tokenOwner sdk.AccAddress, amount types.Coins) error {
	if err := k.isBurnable(ctx, permissionOwner, tokenOwner, amount); err != nil {
		return err
	}

	if err := k.BurnSupply(ctx, tokenOwner, amount); err != nil {
		return err
	}
	return nil
}

func (k Keeper) isBurnable(ctx sdk.Context, permissionOwner, tokenOwner sdk.AccAddress, amount types.Coins) error {
	perm := types.NewBurnPermission()
	if !k.HasPermission(ctx, permissionOwner, perm) {
		return sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", permissionOwner.String(), perm.String())
	}

	if !k.HasCoins(ctx, tokenOwner, amount) {
		return sdkerrors.Wrapf(types.ErrInsufficientToken, "%v has not enough coins for %v", tokenOwner.String(), amount)
	}
	return nil
}

func (k Keeper) BurnNFT(ctx sdk.Context, from sdk.AccAddress, tokenIDs ...string) error {
	for _, tokenID := range tokenIDs {
		if err := k.burnNFT(ctx, from, from, tokenID); err != nil {
			return err
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnNFT,
			sdk.NewAttribute(types.AttributeKeyContractID, k.getContractID(ctx)),
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

func (k Keeper) BurnNFTFrom(ctx sdk.Context, proxy sdk.AccAddress, from sdk.AccAddress, tokenIDs ...string) error {
	if !k.IsApproved(ctx, proxy, from) {
		return sdkerrors.Wrapf(types.ErrCollectionNotApproved, "Proxy: %s, Approver: %s, ContractID: %s", proxy.String(), from.String(), k.getContractID(ctx))
	}

	for _, tokenID := range tokenIDs {
		if err := k.burnNFT(ctx, proxy, from, tokenID); err != nil {
			return err
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnNFTFrom,
			sdk.NewAttribute(types.AttributeKeyContractID, k.getContractID(ctx)),
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

func (k Keeper) burnNFT(ctx sdk.Context, permissionOwner, tokenOwner sdk.AccAddress, tokenID string) error {
	token, err := k.GetNFT(ctx, tokenID)
	if err != nil {
		return err
	}

	perm := types.NewBurnPermission()
	if !k.HasPermission(ctx, permissionOwner, perm) {
		return sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", permissionOwner.String(), perm.String())
	}

	if !token.GetOwner().Equals(tokenOwner) {
		return sdkerrors.Wrapf(types.ErrTokenNotOwnedBy, "TokenID: %s, Owner: %s", tokenID, tokenOwner.String())
	}

	err = k.burnNFTRecursive(ctx, token, tokenOwner)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) burnNFTRecursive(ctx sdk.Context, token types.NFT, from sdk.AccAddress) error {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeOperationBurnNFT,
			sdk.NewAttribute(types.AttributeKeyTokenID, token.GetTokenID()),
		),
	})

	children, err := k.ChildrenOf(ctx, token.GetTokenID())
	if err != nil {
		return err
	}

	for _, child := range children {
		err = k.burnNFTRecursive(ctx, child.(types.NFT), from)
		if err != nil {
			return err
		}
	}

	parent, err := k.ParentOf(ctx, token.GetTokenID())
	if err != nil {
		return err
	}
	if parent != nil {
		_, err = k.detach(ctx, from, token.GetTokenID())
		if err != nil {
			return err
		}
	}
	err = k.DeleteToken(ctx, token.GetTokenID())
	if err != nil {
		return err
	}
	err = k.BurnSupply(ctx, token.GetOwner(), types.OneCoins(token.GetTokenID()))
	if err != nil {
		return err
	}
	return nil
}
