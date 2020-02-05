package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
)

func (k Keeper) TransferFT(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, symbol string, amount sdk.Int) sdk.Error {
	coin := sdk.NewCoins(sdk.NewCoin(symbol, amount))

	_, err := k.bankKeeper.SubtractCoins(ctx, from, coin)
	if err != nil {
		return err
	}

	_, err = k.bankKeeper.AddCoins(ctx, to, coin)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransfer,
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTo, to.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, symbol),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})

	return nil
}

func (k Keeper) TransferCFT(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string, amount sdk.Int) sdk.Error {
	coin := sdk.NewCoins(sdk.NewCoin(symbol+tokenID, amount))
	_, err := k.bankKeeper.SubtractCoins(ctx, from, coin)
	if err != nil {
		return err
	}

	_, err = k.bankKeeper.AddCoins(ctx, to, coin)
	if err != nil {
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

func (k Keeper) TransferNFT(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, symbol string) sdk.Error {
	store := ctx.KVStore(k.storeKey)

	token, err := k.GetToken(ctx, symbol, "")
	if err != nil {
		return err
	}

	nft, ok := token.(*types.BaseNFT)
	if !ok {
		return types.ErrTokenNotNFT(types.DefaultCodespace, token.GetDenom())
	}
	if !from.Equals(nft.Owner) {
		return types.ErrTokenNotOwnedBy(types.DefaultCodespace, token.GetDenom(), from)
	}
	if !from.Equals(to) {
		if err := k.moveNFToken(ctx, store, from, to, nft); err != nil {
			return err
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransferNFT,
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTo, to.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, symbol),
		),
	})

	return nil
}

func (k Keeper) TransferCNFT(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string) sdk.Error {
	store := ctx.KVStore(k.storeKey)

	token, err := k.GetToken(ctx, symbol, tokenID)
	if err != nil {
		return err
	}

	cnft, ok := token.(types.CollectiveNFT)
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
		if err := k.moveCNFToken(ctx, store, from, to, cnft); err != nil {
			return err
		}
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

func (k Keeper) moveTokenInternal(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amount sdk.Coins) sdk.Error {
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

func (k Keeper) moveNFToken(ctx sdk.Context, store sdk.KVStore, from sdk.AccAddress, to sdk.AccAddress, token *types.BaseNFT) sdk.Error {
	amount := sdk.NewCoins(sdk.NewInt64Coin(token.GetDenom(), 1))
	if err := k.moveTokenInternal(ctx, from, to, amount); err != nil {
		return err
	}

	token.Owner = to
	return k.ModifyToken(ctx, token)
}

func (k Keeper) moveCNFToken(ctx sdk.Context, store sdk.KVStore, from sdk.AccAddress, to sdk.AccAddress, token types.CollectiveNFT) sdk.Error {
	children, err := k.ChildrenOf(ctx, token.GetSymbol(), token.GetTokenID())
	if err != nil {
		return err
	}

	for _, child := range children {
		err := k.moveCNFToken(ctx, store, from, to, child.(types.CollectiveNFT))
		if err != nil {
			return err
		}
	}

	amount := sdk.NewCoins(sdk.NewInt64Coin(token.GetDenom(), 1))
	if err := k.moveTokenInternal(ctx, from, to, amount); err != nil {
		return err
	}

	token.SetOwner(to)
	return k.ModifyToken(ctx, token)
}
