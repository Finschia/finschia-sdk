package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/types"
)

type IssueKeeper interface {
	IssueFT(ctx sdk.Context, owner sdk.AccAddress, token types.FT, amount sdk.Int) sdk.Error
	IssueNFT(ctx sdk.Context, owner sdk.AccAddress, tokenType string) sdk.Error
}

func (k Keeper) IssueFT(ctx sdk.Context, owner sdk.AccAddress, to sdk.AccAddress, token types.FT, amount sdk.Int) sdk.Error {
	if !k.ExistCollection(ctx, token.GetContractID()) {
		return types.ErrCollectionNotExist(types.DefaultCodespace, token.GetContractID())
	}

	err := k.SetToken(ctx, token)
	if err != nil {
		return err
	}

	err = k.MintSupply(ctx, token.GetContractID(), to, types.NewCoins(types.NewCoin(token.GetTokenID(), amount)))
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeIssueFT,
			sdk.NewAttribute(types.AttributeKeyContractID, token.GetContractID()),
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

func (k Keeper) IssueNFT(ctx sdk.Context, tokenType types.TokenType, owner sdk.AccAddress) sdk.Error {
	if !k.ExistCollection(ctx, tokenType.GetContractID()) {
		return types.ErrCollectionNotExist(types.DefaultCodespace, tokenType.GetContractID())
	}

	err := k.SetTokenType(ctx, tokenType)
	if err != nil {
		return err
	}

	mintPerm := types.NewMintPermission(tokenType.GetContractID())
	k.AddPermission(ctx, owner, mintPerm)
	burnPerm := types.NewBurnPermission(tokenType.GetContractID())
	k.AddPermission(ctx, owner, burnPerm)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeIssueNFT,
			sdk.NewAttribute(types.AttributeKeyContractID, tokenType.GetContractID()),
			sdk.NewAttribute(types.AttributeKeyTokenType, tokenType.GetTokenType()),
		),
	})

	return nil
}
