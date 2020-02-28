package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/types"
)

type MintKeeper interface {
	MintFT(ctx sdk.Context, contractID string, from, to sdk.AccAddress, amount types.Coins) sdk.Error
	MintNFT(ctx sdk.Context, contractID string, from sdk.AccAddress, token types.NFT) sdk.Error
}

func (k Keeper) MintFT(ctx sdk.Context, contractID string, from, to sdk.AccAddress, amount types.Coins) sdk.Error {
	for _, coin := range amount {
		token, err := k.GetToken(ctx, contractID, coin.Denom)
		if err != nil {
			return err
		}
		if err := k.isMintable(ctx, contractID, token, from); err != nil {
			return err
		}
	}
	err := k.MintSupply(ctx, contractID, to, amount)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMintFT,
			sdk.NewAttribute(types.AttributeKeyContractID, contractID),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTo, to.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})
	return nil
}

func (k Keeper) MintNFT(ctx sdk.Context, contractID string, from sdk.AccAddress, token types.NFT) sdk.Error {
	if !k.HasTokenType(ctx, contractID, token.GetTokenType()) {
		return types.ErrTokenTypeNotExist(types.DefaultCodespace, token.GetContractID(), token.GetTokenType())
	}

	perm := types.NewMintPermission(token.GetContractID())
	if !k.HasPermission(ctx, from, perm) {
		return types.ErrTokenNoPermission(types.DefaultCodespace, from, perm)
	}

	err := k.SetToken(ctx, contractID, token)
	if err != nil {
		return err
	}

	err = k.MintSupply(ctx, contractID, token.GetOwner(), types.OneCoins(token.GetTokenID()))
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMintNFT,
			sdk.NewAttribute(types.AttributeKeyContractID, token.GetContractID()),
			sdk.NewAttribute(types.AttributeKeyName, token.GetName()),
			sdk.NewAttribute(types.AttributeKeyTokenID, token.GetTokenID()),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTo, token.GetOwner().String()),
		),
	})

	return nil
}
func (k Keeper) isMintable(ctx sdk.Context, contractID string, token types.Token, from sdk.AccAddress) sdk.Error {
	ft, ok := token.(types.FT)
	if !ok {
		return types.ErrTokenNotMintable(types.DefaultCodespace, contractID, token.GetTokenID())
	}

	if !ft.GetMintable() {
		return types.ErrTokenNotMintable(types.DefaultCodespace, contractID, token.GetTokenID())
	}
	perm := types.NewMintPermission(contractID)
	if !k.HasPermission(ctx, from, perm) {
		return types.ErrTokenNoPermission(types.DefaultCodespace, from, perm)
	}
	return nil
}
