package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link-modules/x/collection/internal/types"
)

type IssueKeeper interface {
	IssueFT(ctx sdk.Context, owner sdk.AccAddress, token types.FT, amount sdk.Int) error
	IssueNFT(ctx sdk.Context, owner sdk.AccAddress, tokenType string) error
}

func (k Keeper) IssueFT(ctx sdk.Context, owner, to sdk.AccAddress, token types.FT, amount sdk.Int) error {
	if !k.ExistCollection(ctx) {
		return sdkerrors.Wrapf(types.ErrCollectionNotExist, "ContractID: %s", k.getContractID(ctx))
	}
	err := k.SetToken(ctx, token)
	if err != nil {
		return err
	}

	err = k.MintSupply(ctx, to, types.NewCoins(types.NewCoin(token.GetTokenID(), amount)))
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeIssueFT,
			sdk.NewAttribute(types.AttributeKeyContractID, k.getContractID(ctx)),
			sdk.NewAttribute(types.AttributeKeyName, token.GetName()),
			sdk.NewAttribute(types.AttributeKeyTokenID, token.GetTokenID()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyTo, to.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyMintable, strconv.FormatBool(token.GetMintable())),
			sdk.NewAttribute(types.AttributeKeyDecimals, token.GetDecimals().String()),
		),
	})

	return nil
}

func (k Keeper) IssueNFT(ctx sdk.Context, tokenType types.TokenType, owner sdk.AccAddress) error {
	if !k.ExistCollection(ctx) {
		return sdkerrors.Wrapf(types.ErrCollectionNotExist, "ContractID: %s", k.getContractID(ctx))
	}

	err := k.SetTokenType(ctx, tokenType)
	if err != nil {
		return err
	}

	mintPerm := types.NewMintPermission()
	k.AddPermission(ctx, owner, mintPerm)
	burnPerm := types.NewBurnPermission()
	k.AddPermission(ctx, owner, burnPerm)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeIssueNFT,
			sdk.NewAttribute(types.AttributeKeyContractID, k.getContractID(ctx)),
			sdk.NewAttribute(types.AttributeKeyTokenType, tokenType.GetTokenType()),
		),
	})

	return nil
}
