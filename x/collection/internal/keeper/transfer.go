package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/types"
)

type TransferKeeper interface {
	TransferFT(ctx sdk.Context, contractID string, from sdk.AccAddress, to sdk.AccAddress, amount ...types.Coin) sdk.Error
	TransferNFT(ctx sdk.Context, contractID string, from sdk.AccAddress, to sdk.AccAddress, tokenID ...string) sdk.Error
	TransferFTFrom(ctx sdk.Context, contractID string, proxy sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, amount ...types.Coin) sdk.Error
	TransferNFTFrom(ctx sdk.Context, contractID string, proxy sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, tokenID ...string) sdk.Error
}

var _ TransferKeeper = (*Keeper)(nil)

func (k Keeper) TransferFT(ctx sdk.Context, contractID string, from sdk.AccAddress, to sdk.AccAddress, amount ...types.Coin) sdk.Error {
	if err := k.transferFT(ctx, contractID, from, to, amount); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransferFT,
			sdk.NewAttribute(types.AttributeKeyContractID, contractID),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTo, to.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, types.NewCoins(amount...).String()),
		),
	})

	return nil
}

func (k Keeper) transferFT(ctx sdk.Context, contractID string, from sdk.AccAddress, to sdk.AccAddress, amount types.Coins) sdk.Error {
	return k.SendCoins(ctx, contractID, from, to, amount)
}

func (k Keeper) TransferNFT(ctx sdk.Context, contractID string, from sdk.AccAddress, to sdk.AccAddress, tokenIDs ...string) sdk.Error {
	for _, tokenID := range tokenIDs {
		if err := k.transferNFT(ctx, contractID, from, to, tokenID); err != nil {
			return err
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransferNFT,
			sdk.NewAttribute(types.AttributeKeyContractID, contractID),
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

func (k Keeper) transferNFT(ctx sdk.Context, contractID string, from sdk.AccAddress, to sdk.AccAddress, tokenID string) sdk.Error {
	store := ctx.KVStore(k.storeKey)

	token, err := k.GetToken(ctx, contractID, tokenID)
	if err != nil {
		return err
	}

	nft, ok := token.(types.NFT)
	if !ok {
		return types.ErrTokenNotNFT(types.DefaultCodespace, token.GetTokenID())
	}
	childToParentKey := types.TokenChildToParentKey(contractID, nft.GetTokenID())
	if store.Has(childToParentKey) {
		return types.ErrTokenCannotTransferChildToken(types.DefaultCodespace, token.GetTokenID())
	}
	if !from.Equals(nft.GetOwner()) {
		return types.ErrTokenNotOwnedBy(types.DefaultCodespace, token.GetTokenID(), from)
	}
	if !from.Equals(to) {
		if err := k.moveNFToken(ctx, contractID, from, to, nft); err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) TransferFTFrom(ctx sdk.Context, contractID string, proxy sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, amount ...types.Coin) sdk.Error {
	if !k.IsApproved(ctx, contractID, proxy, from) {
		return types.ErrCollectionNotApproved(types.DefaultCodespace, proxy.String(), from.String(), contractID)
	}

	if err := k.transferFT(ctx, contractID, from, to, amount); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransferFTFrom,
			sdk.NewAttribute(types.AttributeKeyContractID, contractID),
			sdk.NewAttribute(types.AttributeKeyProxy, proxy.String()),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTo, to.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, types.NewCoins(amount...).String()),
		),
	})

	return nil
}

//nolint:dupl
func (k Keeper) TransferNFTFrom(ctx sdk.Context, contractID string, proxy sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, tokenIDs ...string) sdk.Error {
	if !k.IsApproved(ctx, contractID, proxy, from) {
		return types.ErrCollectionNotApproved(types.DefaultCodespace, proxy.String(), from.String(), contractID)
	}

	for _, tokenID := range tokenIDs {
		if err := k.transferNFT(ctx, contractID, from, to, tokenID); err != nil {
			return err
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransferNFTFrom,
			sdk.NewAttribute(types.AttributeKeyContractID, contractID),
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

func (k Keeper) moveNFToken(ctx sdk.Context, contractID string, from sdk.AccAddress, to sdk.AccAddress, token types.NFT) sdk.Error {
	if from.Equals(to) {
		return nil
	}
	children, err := k.ChildrenOf(ctx, token.GetContractID(), token.GetTokenID())
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
		err := k.moveNFToken(ctx, contractID, from, to, child.(types.NFT))
		if err != nil {
			return err
		}
	}

	amount := types.NewCoins(types.NewInt64Coin(token.GetTokenID(), 1))
	if err := k.SendCoins(ctx, contractID, from, to, amount); err != nil {
		return err
	}

	return k.UpdateToken(ctx, contractID, token.SetOwner(to))
}
