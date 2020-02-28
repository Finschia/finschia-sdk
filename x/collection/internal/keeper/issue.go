package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/types"
)

type IssueKeeper interface {
	IssueFT(ctx sdk.Context, contractID string, owner sdk.AccAddress, token types.FT, amount sdk.Int) sdk.Error
	IssueNFT(ctx sdk.Context, contractID string, owner sdk.AccAddress, tokenType string) sdk.Error
}

func (k Keeper) IssueFT(ctx sdk.Context, contractID string, owner sdk.AccAddress, token types.FT, amount sdk.Int) sdk.Error {
	err := k.SetToken(ctx, contractID, token)
	if err != nil {
		return err
	}

	err = k.MintSupply(ctx, contractID, owner, types.NewCoins(types.NewCoin(token.GetTokenID(), amount)))
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
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyMintable, strconv.FormatBool(token.GetMintable())),
			sdk.NewAttribute(types.AttributeKeyDecimals, token.GetDecimals().String()),
		),
	})

	return nil
}

func (k Keeper) IssueNFT(ctx sdk.Context, contractID string, tokenType types.TokenType, owner sdk.AccAddress) sdk.Error {
	err := k.SetTokenType(ctx, contractID, tokenType)
	if err != nil {
		return err
	}

	mintPerm := types.NewMintPermission(contractID)
	k.AddPermission(ctx, owner, mintPerm)
	burnPerm := types.NewBurnPermission(contractID)
	k.AddPermission(ctx, owner, burnPerm)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeIssueNFT,
			sdk.NewAttribute(types.AttributeKeyContractID, contractID),
			sdk.NewAttribute(types.AttributeKeyTokenType, tokenType.GetTokenType()),
		),
	})

	return nil
}
