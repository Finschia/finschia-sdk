package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/collection/internal/types"
)

type TransferKeeper interface {
	TransferFT(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amount ...types.Coin) error
	TransferNFT(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, tokenID ...string) error
	TransferFTFrom(ctx sdk.Context, proxy sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, amount ...types.Coin) error
	TransferNFTFrom(ctx sdk.Context, proxy sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, tokenID ...string) error
}

var _ TransferKeeper = (*Keeper)(nil)

func (k Keeper) TransferFT(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amount ...types.Coin) error {
	if err := k.transferFT(ctx, from, to, amount); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransferFT,
			sdk.NewAttribute(types.AttributeKeyContractID, k.getContractID(ctx)),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTo, to.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, types.NewCoins(amount...).String()),
		),
	})

	return nil
}

func (k Keeper) transferFT(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amount types.Coins) error {
	return k.SendCoins(ctx, from, to, amount)
}

func (k Keeper) TransferNFT(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, tokenIDs ...string) error {
	for _, tokenID := range tokenIDs {
		if err := k.transferNFT(ctx, from, to, tokenID); err != nil {
			return err
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransferNFT,
			sdk.NewAttribute(types.AttributeKeyContractID, k.getContractID(ctx)),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTo, to.String()),
		),
	})
	for _, tokenID := range tokenIDs {
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeTransferNFT,
				sdk.NewAttribute(types.AttributeKeyTokenID, tokenID),
			),
		})
	}

	return nil
}

func (k Keeper) transferNFT(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, tokenID string) error {
	store := ctx.KVStore(k.storeKey)

	token, err := k.GetToken(ctx, tokenID)
	if err != nil {
		return err
	}

	nft, ok := token.(types.NFT)
	if !ok {
		return sdkerrors.Wrapf(types.ErrTokenNotNFT, "TokenID: %s", token.GetTokenID())
	}
	childToParentKey := types.TokenChildToParentKey(k.getContractID(ctx), nft.GetTokenID())
	if store.Has(childToParentKey) {
		return sdkerrors.Wrapf(types.ErrTokenCannotTransferChildToken, "TokenID: %s", token.GetTokenID())
	}
	if !from.Equals(nft.GetOwner()) {
		return sdkerrors.Wrapf(types.ErrTokenNotOwnedBy, "TokenID: %s, Owner: %s", token.GetTokenID(), from.String())
	}
	if !from.Equals(to) {
		if err := k.moveNFToken(ctx, from, to, nft); err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) TransferFTFrom(ctx sdk.Context, proxy sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, amount ...types.Coin) error {
	if !k.IsApproved(ctx, proxy, from) {
		return sdkerrors.Wrapf(types.ErrCollectionNotApproved, "Proxy: %s, Approver: %s, ContractID: %s", proxy.String(), from.String(), k.getContractID(ctx))
	}

	if err := k.transferFT(ctx, from, to, amount); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransferFTFrom,
			sdk.NewAttribute(types.AttributeKeyContractID, k.getContractID(ctx)),
			sdk.NewAttribute(types.AttributeKeyProxy, proxy.String()),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTo, to.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, types.NewCoins(amount...).String()),
		),
	})

	return nil
}

// nolint:dupl
func (k Keeper) TransferNFTFrom(ctx sdk.Context, proxy sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, tokenIDs ...string) error {
	if !k.IsApproved(ctx, proxy, from) {
		return sdkerrors.Wrapf(types.ErrCollectionNotApproved, "Proxy: %s, Approver: %s, ContractID: %s", proxy.String(), from.String(), k.getContractID(ctx))
	}

	for _, tokenID := range tokenIDs {
		if err := k.transferNFT(ctx, from, to, tokenID); err != nil {
			return err
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransferNFTFrom,
			sdk.NewAttribute(types.AttributeKeyContractID, k.getContractID(ctx)),
			sdk.NewAttribute(types.AttributeKeyProxy, proxy.String()),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTo, to.String()),
		),
	})
	for _, tokenID := range tokenIDs {
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeTransferNFTFrom,
				sdk.NewAttribute(types.AttributeKeyTokenID, tokenID),
			),
		})
	}

	return nil
}

func (k Keeper) moveNFToken(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, token types.NFT) error {
	if from.Equals(to) {
		return nil
	}
	children, err := k.ChildrenOf(ctx, token.GetTokenID())
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeOperationTransferNFT,
			sdk.NewAttribute(types.AttributeKeyTokenID, token.GetTokenID()),
		),
	})

	for _, child := range children {
		err := k.moveNFToken(ctx, from, to, child.(types.NFT))
		if err != nil {
			return err
		}
	}

	if err := k.ChangeNFTOwner(ctx, from, to, token.GetTokenID()); err != nil {
		return err
	}
	token.SetOwner(to)
	return k.UpdateToken(ctx, token)
}
