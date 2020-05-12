package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link/x/collection/internal/types"
)

type MintKeeper interface {
	MintFT(ctx sdk.Context, contractID string, from, to sdk.AccAddress, amount types.Coins) error
	MintNFT(ctx sdk.Context, from sdk.AccAddress, token types.NFT) error
}

func (k Keeper) MintFT(ctx sdk.Context, contractID string, from, to sdk.AccAddress, amount types.Coins) error {
	for _, coin := range amount {
		token, err := k.GetToken(ctx, contractID, coin.Denom)
		if err != nil {
			return err
		}
		if err := k.isMintable(ctx, token, from); err != nil {
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

func (k Keeper) MintNFT(ctx sdk.Context, from sdk.AccAddress, token types.NFT) error {
	if !k.HasTokenType(ctx, token.GetContractID(), token.GetTokenType()) {
		return sdkerrors.Wrapf(types.ErrTokenTypeNotExist, "ContractID: %s, TokenType: %s", token.GetContractID(), token.GetTokenType())
	}

	perm := types.NewMintPermission()
	if !k.HasPermission(ctx, token.GetContractID(), from, perm) {
		return sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", from.String(), perm.String())
	}

	err := k.SetToken(ctx, token)
	if err != nil {
		return err
	}

	err = k.MintSupply(ctx, token.GetContractID(), token.GetOwner(), types.OneCoins(token.GetTokenID()))
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

func (k Keeper) isMintable(ctx sdk.Context, token types.Token, from sdk.AccAddress) error {
	ft, ok := token.(types.FT)
	if !ok {
		return sdkerrors.Wrapf(types.ErrTokenNotMintable, "ContractID: %s, TokenID: %s", token.GetContractID(), token.GetTokenID())
	}

	if !ft.GetMintable() {
		return sdkerrors.Wrapf(types.ErrTokenNotMintable, "ContractID: %s, TokenID: %s", token.GetContractID(), token.GetTokenID())
	}
	perm := types.NewMintPermission()
	if !k.HasPermission(ctx, token.GetContractID(), from, perm) {
		return sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", from.String(), perm.String())
	}
	return nil
}
