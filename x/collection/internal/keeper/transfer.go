package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/types"
)

type TransferKeeper interface {
	TransferFT(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string, amount sdk.Int) sdk.Error
	TransferNFT(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string) sdk.Error
	TransferFTFrom(ctx sdk.Context, proxy sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string, amount sdk.Int) sdk.Error
	TransferNFTFrom(ctx sdk.Context, proxy sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string) sdk.Error
}

func (k Keeper) TransferFT(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string, amount sdk.Int) sdk.Error {
	if err := k.transferFT(ctx, from, to, symbol, tokenID, amount); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransferFT,
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTo, to.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, symbol),
			sdk.NewAttribute(types.AttributeKeyTokenID, tokenID),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})

	return nil
}

func (k Keeper) transferFT(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string, amount sdk.Int) sdk.Error {
	coin := sdk.NewCoins(sdk.NewCoin(symbol+tokenID, amount))
	_, err := k.bankKeeper.SubtractCoins(ctx, from, coin)
	if err != nil {
		return err
	}

	_, err = k.bankKeeper.AddCoins(ctx, to, coin)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) TransferNFT(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string) sdk.Error {
	if err := k.transferNFT(ctx, from, to, symbol, tokenID); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransferNFT,
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTo, to.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, symbol),
			sdk.NewAttribute(types.AttributeKeyTokenID, tokenID),
		),
	})

	return nil
}

func (k Keeper) transferNFT(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string) sdk.Error {
	store := ctx.KVStore(k.storeKey)

	token, err := k.GetToken(ctx, symbol, tokenID)
	if err != nil {
		return err
	}

	nft, ok := token.(types.NFT)
	if !ok {
		return types.ErrTokenNotNFT(types.DefaultCodespace, token.GetDenom())
	}
	childToParentKey := types.TokenChildToParentKey(nft)
	if store.Has(childToParentKey) {
		return types.ErrTokenCannotTransferChildToken(types.DefaultCodespace, token.GetDenom())
	}
	if !from.Equals(nft.GetOwner()) {
		return types.ErrTokenNotOwnedBy(types.DefaultCodespace, token.GetDenom(), from)
	}
	if !from.Equals(to) {
		if err := k.moveNFToken(ctx, from, to, nft); err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) TransferFTFrom(ctx sdk.Context, proxy sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string, amount sdk.Int) sdk.Error {
	if !k.IsApproved(ctx, proxy, from, symbol) {
		return types.ErrCollectionNotApproved(types.DefaultCodespace, proxy.String(), from.String(), symbol)
	}

	if err := k.transferFT(ctx, from, to, symbol, tokenID, amount); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransferFTFrom,
			sdk.NewAttribute(types.AttributeKeyProxy, proxy.String()),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTo, to.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, symbol),
			sdk.NewAttribute(types.AttributeKeyTokenID, tokenID),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})

	return nil
}

//nolint:dupl
func (k Keeper) TransferNFTFrom(ctx sdk.Context, proxy sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string) sdk.Error {
	if !k.IsApproved(ctx, proxy, from, symbol) {
		return types.ErrCollectionNotApproved(types.DefaultCodespace, proxy.String(), from.String(), symbol)
	}

	if err := k.transferNFT(ctx, from, to, symbol, tokenID); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransferNFTFrom,
			sdk.NewAttribute(types.AttributeKeyProxy, proxy.String()),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTo, to.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, symbol),
			sdk.NewAttribute(types.AttributeKeyTokenID, tokenID),
		),
	})

	return nil
}

func (k Keeper) moveToken(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amount sdk.Coins) sdk.Error {
	_, err := k.bankKeeper.SubtractCoins(ctx, from, amount)
	if err != nil {
		return err
	}

	_, err = k.bankKeeper.AddCoins(ctx, to, amount)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) moveNFToken(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, token types.NFT) sdk.Error {
	if from.Equals(to) {
		return nil
	}
	children, err := k.ChildrenOf(ctx, token.GetSymbol(), token.GetTokenID())
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
		if err := k.moveNFToken(ctx, from, to, child.(types.NFT)); err != nil {
			return err
		}
	}

	amount := sdk.NewCoins(sdk.NewInt64Coin(token.GetDenom(), 1))
	if err := k.moveToken(ctx, from, to, amount); err != nil {
		return err
	}

	token.SetOwner(to)
	return k.UpdateToken(ctx, token)
}
