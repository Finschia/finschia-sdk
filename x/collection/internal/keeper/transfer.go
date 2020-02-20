package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/types"
)

type TransferKeeper interface {
	TransferCFT(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string, amount sdk.Int) sdk.Error
	TransferCNFT(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string) sdk.Error
	TransferCFTFrom(ctx sdk.Context, proxy sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string, amount sdk.Int) sdk.Error
	TransferCNFTFrom(ctx sdk.Context, proxy sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string) sdk.Error
}

func (k Keeper) TransferCFT(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string, amount sdk.Int) sdk.Error {
	if err := k.transferCFT(ctx, from, to, symbol, tokenID, amount); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransferCFT,
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTo, to.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, symbol),
			sdk.NewAttribute(types.AttributeKeyTokenID, tokenID),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})

	return nil
}

func (k Keeper) transferCFT(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string, amount sdk.Int) sdk.Error {
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

func (k Keeper) TransferCNFT(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string) sdk.Error {
	if err := k.transferCNFT(ctx, from, to, symbol, tokenID); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransferCNFT,
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTo, to.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, symbol),
			sdk.NewAttribute(types.AttributeKeyTokenID, tokenID),
		),
	})

	return nil
}

func (k Keeper) transferCNFT(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string) sdk.Error {
	store := ctx.KVStore(k.storeKey)

	token, err := k.GetToken(ctx, symbol, tokenID)
	if err != nil {
		return err
	}

	cnft, ok := token.(types.NFT)
	if !ok {
		return types.ErrTokenNotCNFT(types.DefaultCodespace, token.GetDenom())
	}
	childToParentKey := types.TokenChildToParentKey(cnft)
	if store.Has(childToParentKey) {
		return types.ErrTokenCannotTransferChildToken(types.DefaultCodespace, token.GetDenom())
	}
	if !from.Equals(cnft.GetOwner()) {
		return types.ErrTokenNotOwnedBy(types.DefaultCodespace, token.GetDenom(), from)
	}
	if !from.Equals(to) {
		if err := k.moveCNFToken(ctx, from, to, cnft); err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) TransferCFTFrom(ctx sdk.Context, proxy sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string, amount sdk.Int) sdk.Error {
	if !k.IsApproved(ctx, proxy, from, symbol) {
		return types.ErrCollectionNotApproved(types.DefaultCodespace, proxy.String(), from.String(), symbol)
	}

	if err := k.transferCFT(ctx, from, to, symbol, tokenID, amount); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransferCFTFrom,
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
func (k Keeper) TransferCNFTFrom(ctx sdk.Context, proxy sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string) sdk.Error {
	if !k.IsApproved(ctx, proxy, from, symbol) {
		return types.ErrCollectionNotApproved(types.DefaultCodespace, proxy.String(), from.String(), symbol)
	}

	if err := k.transferCNFT(ctx, from, to, symbol, tokenID); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransferCNFTFrom,
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

func (k Keeper) moveCNFToken(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, token types.NFT) sdk.Error {
	if from.Equals(to) {
		return nil
	}
	children, err := k.ChildrenOf(ctx, token.GetSymbol(), token.GetTokenID())
	if err != nil {
		return err
	}

	for _, child := range children {
		err := k.moveCNFToken(ctx, from, to, child.(types.NFT))
		if err != nil {
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
